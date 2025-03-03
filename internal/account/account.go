package account

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"gorm.io/gorm"
	"gorm.io/hints"

	"google.golang.org/grpc/grpclog"

	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/micro/v2/utils/dbutil"
	"github.com/gidyon/micro/v2/utils/errs"
	"github.com/gidyon/micro/v2/utils/mdutil"
	"github.com/gidyon/micro/v2/utils/templateutil"
	"github.com/gidyon/services/internal/pkg/fauth"
	"github.com/gidyon/services/pkg/api/account"
	"github.com/gidyon/services/pkg/api/messaging"
	"github.com/gidyon/services/pkg/utils/timeutil"

	redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/speps/go-hashids"
	"google.golang.org/grpc/codes"
)

const templateName = "base"

type accountAPIServer struct {
	account.UnsafeAccountAPIServer
	activationURL string
	tpl           *template.Template
	*Options
}

// Options contain parameters for NewAccountAPI
type Options struct {
	AppName            string
	EmailDisplayName   string
	DefaultEmailSender string
	TemplatesDir       string
	ActivationURL      string
	PaginationHasher   *hashids.HashID
	AuthAPI            auth.API
	SQLDBWrites        *gorm.DB
	SQLDBReads         *gorm.DB
	RedisDBWrites      *redis.Client
	RedisDBReads       *redis.Client
	Logger             grpclog.LoggerV2
	MessagingClient    messaging.MessagingClient
	FirebaseAuth       fauth.FirebaseAuthClient
}

// NewAccountAPI creates an account API singleton
func NewAccountAPI(ctx context.Context, opt *Options) (account.AccountAPIServer, error) {
	// Validation
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("context")
	case opt == nil:
		err = errs.NilObject("options")
	case opt.AppName == "":
		err = errs.MissingField("app name")
	case opt.EmailDisplayName == "":
		err = errs.MissingField("email display name")
	// case opt.TemplatesDir == "":
	// 	err = errs.MissingField("templates directory")
	// case opt.ActivationURL == "":
	// 	err = errs.MissingField("activation url")
	case opt.PaginationHasher == nil:
		err = errs.MissingField("pagination PaginationHasher")
	case opt.AuthAPI == nil:
		err = errs.MissingField("authentication API")
	case opt.SQLDBWrites == nil:
		err = errs.NilObject("sql writes db")
	case opt.SQLDBReads == nil:
		err = errs.NilObject("sql reads db")
	case opt.RedisDBWrites == nil:
		err = errs.NilObject("redis writes db")
	case opt.RedisDBReads == nil:
		err = errs.NilObject("redis reads db")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case opt.MessagingClient == nil:
		err = errs.NilObject("messaging client")
		// case opt.FirebaseAuth == nil:
		// 	err = errs.NilObject("firebase auth")
	}
	if err != nil {
		return nil, err
	}

	// Account API
	accountAPI := &accountAPIServer{
		activationURL: opt.ActivationURL,
		Options:       opt,
	}

	if opt.TemplatesDir != "" {
		// Read template files from directory
		tFiles, err := templateutil.ReadFiles(opt.TemplatesDir)
		if err != nil {
			return nil, errs.WrapErrorWithMsg(err, "failed to read template files in directory")
		}

		// Parse template
		accountAPI.tpl, err = templateutil.ParseTemplate(tFiles...)
		if err != nil {
			return nil, errs.WrapErrorWithMsg(err, "failed to parse template")
		}
	}

	// Perform auto migration
	if !accountAPI.SQLDBWrites.Migrator().HasTable(accountsTable) {
		err = accountAPI.SQLDBWrites.AutoMigrate(&Account{})
		if err != nil {
			return nil, errs.WrapErrorWithMsg(err, "failed to automigrate accounts table")
		}
	}

	if !accountAPI.SQLDBWrites.Migrator().HasIndex(&Account{}, dbutil.FullTextIndex) {
		// Create a full text search index
		err = dbutil.CreateFullTextIndex(accountAPI.SQLDBWrites, accountsTable, "names", "email", "phone", "linked_accounts")
		if err != nil {
			return nil, errs.WrapErrorWithMsg(err, "failed to create full text index")
		}
	}

	return accountAPI, nil
}

func (accountAPI *accountAPIServer) SignInExternal(
	ctx context.Context, signInReq *account.SignInExternalRequest,
) (*account.SignInResponse, error) {
	if accountAPI.FirebaseAuth == nil {
		return nil, errs.WrapMessage(codes.Unavailable, "firebase auth not available")
	}

	// Validation
	var err error
	switch {
	case signInReq == nil:
		err = errs.NilObject("sign in request")
	case signInReq.ProjectId == "":
		err = errs.MissingField("project id")
	case signInReq.Account == nil:
		err = errs.NilObject("account")
	case signInReq.Account.Names == "":
		err = errs.MissingField("names")
	case signInReq.Account.Email == "" && signInReq.Account.Phone == "":
		err = errs.MissingField("email and phone")
	case signInReq.AuthToken == "":
		err = errs.MissingField("auth token")
	}
	if err != nil {
		return nil, err
	}

	// Verify ID token
	_, err = accountAPI.FirebaseAuth.VerifyIDToken(ctx, signInReq.AuthToken)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to verify firebase ID token")
	}

	var (
		accountDB = &Account{}
		ID        uint
	)

	// Get user
	switch {
	case signInReq.Account.Email != "":
		err = accountAPI.SQLDBWrites.First(accountDB, "email=? AND project_id = ?", signInReq.Account.Email, signInReq.ProjectId).Error
	case signInReq.Account.Phone != "":
		err = accountAPI.SQLDBWrites.First(accountDB, "phone=? AND project_id = ?", signInReq.Account.Phone, signInReq.ProjectId).Error
	}
	switch {
	case err == nil:
		ID = accountDB.AccountID
	case errors.Is(err, gorm.ErrRecordNotFound):
		// Create user
		signInReq.Account.ProjectId = signInReq.ProjectId
		accountDB, err = GetAccountDB(signInReq.Account)
		if err != nil {
			return nil, err
		}
		accountDB.AccountState = account.AccountState_ACTIVE.String()
		err = accountAPI.SQLDBWrites.Create(accountDB).Error
		if err != nil {
			return nil, errs.FailedToSave("account", err)
		}
		return accountAPI.updateSession(ctx, accountDB, "")
	default:
		return nil, errs.FailedToSave("account", err)
	}

	// Omit fields
	omitFields := []string{"project_id", "id_number", "linked_accounts", "password", "primary_group", "account_state", "secondary_groups", "security_question", "security-answer", "account_id", "gender", "created_at"}

	// Update account
	err = accountAPI.SQLDBWrites.Table(accountsTable).Where("account_id= ?", ID).Omit(omitFields...).Updates(accountDB).Error
	if err != nil {
		return nil, errs.FailedToUpdate("account", err)
	}

	return accountAPI.updateSession(ctx, accountDB, "")
}

func (accountAPI *accountAPIServer) RefreshSession(
	ctx context.Context, req *account.RefreshSessionRequest,
) (*account.SignInResponse, error) {
	var (
		ID  int
		err error
	)
	// Validation
	switch {
	case req == nil:
		return nil, errs.NilObject("refresh token request")
	case req.RefreshToken == "":
		return nil, errs.NilObject("refresh token")
	case req.AccountId == "":
		return nil, errs.NilObject("account id")
	default:
		ID, err = strconv.Atoi(req.AccountId)
		if err != nil {
			return nil, errs.IncorrectVal("account id")
		}
	}

	// Ensure that refresh token already exists
	ok, err := accountAPI.RedisDBWrites.SIsMember(ctx, refreshTokenSet(), req.RefreshToken).Result()
	switch {
	case err == nil:
		if !ok {
			return nil, errs.WrapMessage(codes.Unauthenticated, "not signed in")
		}
	case errors.Is(err, redis.Nil) || !ok:
		return nil, errs.WrapMessage(codes.Unauthenticated, "not signed in")
	default:
		return nil, errs.RedisCmdFailed(err, "get")
	}

	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.First(accountDB, "account_id=?", ID).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", req.AccountId)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	switch {
	case accountDB.AccountState == account.AccountState_BLOCKED.String():
		return nil, errs.WrapMessage(codes.PermissionDenied, "account is blocked")
	case accountDB.AccountState == account.AccountState_DELETED.String():
		return nil, errs.WrapMessage(codes.PermissionDenied, "account is deleted")
	}

	return accountAPI.updateSession(ctx, accountDB, req.AccountGroup)
}

func (accountAPI *accountAPIServer) ActivateAccount(
	ctx context.Context, activateReq *account.ActivateAccountRequest,
) (*account.ActivateAccountResponse, error) {
	var (
		ID  int
		err error
	)
	// Validation 1
	switch {
	case activateReq == nil:
		return nil, errs.NilObject("ActivateAccountRequest")
	case activateReq.Token == "":
		return nil, errs.MissingField("token")
	case activateReq.AccountId == "":
		return nil, errs.MissingField("account id")
	default:
		ID, err = strconv.Atoi(activateReq.AccountId)
		if err != nil {
			return nil, errs.IncorrectVal("account id")
		}
	}

	// Retrieve token claims
	payload, err := accountAPI.AuthAPI.AuthorizeActorOrGroup(
		auth.AddTokenMD(ctx, activateReq.Token), activateReq.AccountId, accountAPI.AuthAPI.AdminGroups()...,
	)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Unauthenticated, err, "failed to authorize request")
	}

	// Validation 2
	dev := (strings.ToLower(os.Getenv("MODE")) == "development")
	switch {
	case payload.ID == "":
		if !dev {
			return nil, errs.MissingField("token id")
		}
	}

	// Compare if account account_id matches or if activated by admin
	isOwner := payload.ID == activateReq.AccountId
	isAdmin := accountAPI.AuthAPI.IsAdmin(payload.Group)
	if !isOwner && !isAdmin {
		if !dev {
			switch {
			case !isAdmin:
				return nil, errs.WrapMessage(codes.PermissionDenied, "not admin user")
			case !isOwner:
				return nil, errs.TokenCredentialNotMatching("account id")
			}
		}
	}

	// Check that account exists
	if errors.Is(accountAPI.SQLDBWrites.Select("account_state").
		First(&Account{}, "account_id=?", ID).Error, gorm.ErrRecordNotFound) {
		return nil, errs.DoesNotExist("account", activateReq.AccountId)
	}

	// Update the model of the user to activate their account
	err = accountAPI.SQLDBWrites.Table(accountsTable).Where("account_id=?", ID).
		Update("account_state", account.AccountState_ACTIVE.String()).Error
	if err != nil {
		return nil, errs.FailedToUpdate("account", err)
	}

	return &account.ActivateAccountResponse{}, nil
}

func (accountAPI *accountAPIServer) UpdateAccount(
	ctx context.Context, updateReq *account.UpdateAccountRequest,
) (*empty.Empty, error) {
	var err error
	// Validation
	switch {
	case updateReq == nil:
		return nil, errs.NilObject("update request")
	case updateReq.Account == nil:
		return nil, errs.NilObject("Account")
	case updateReq.Account.AccountId == "":
		return nil, errs.MissingField("AccountID")
	default:
		_, err = strconv.Atoi(updateReq.Account.AccountId)
		if err != nil {
			return nil, errs.WrapMessage(codes.InvalidArgument, "account id is incorrect")
		}
	}

	// Authorization
	payload, err := accountAPI.AuthAPI.AuthorizeActorOrGroup(ctx, updateReq.GetAccount().GetAccountId(), accountAPI.AuthAPI.AdminGroups()...)
	if err != nil {
		return nil, err
	}

	// GetAccount the account details from database
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.Select("account_state").
		First(accountDB, "account_id=?", updateReq.Account.AccountId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", updateReq.Account.AccountId)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	// Check that account is not blocked or deleted
	switch {
	case accountDB.AccountState == account.AccountState_BLOCKED.String(),
		accountDB.AccountState == account.AccountState_DELETED.String():
		return nil, errs.WrapMessage(codes.PermissionDenied, "account is blocked or deleted")
	}

	accountDBX, err := GetAccountDB(updateReq.Account)
	if err != nil {
		return nil, err
	}

	if accountDBX.AccountState == account.AccountState_ACCOUNT_STATE_UNSPECIFIED.String() {
		accountDBX.AccountState = ""
	}
	if accountDBX.Gender == account.Account_GENDER_UNSPECIFIED.String() {
		accountDBX.Gender = ""
	}

	if !accountAPI.AuthAPI.IsAdmin(payload.Group) {
		// Update the model; omit "id", "primary_group", "account_state" and "security profile"
		err = accountAPI.SQLDBWrites.Model(accountDBX).
			Omit("id", "primary_group", "account_state", "password", "security_answer", "security_question").
			Where("account_id=?", updateReq.Account.AccountId).
			Updates(accountDBX).Error
	} else {
		err = accountAPI.SQLDBWrites.Model(accountDBX).
			Where("account_id=?", updateReq.Account.AccountId).
			Updates(accountDBX).Error
	}
	if err != nil {
		return nil, errs.FailedToUpdate("account", err)
	}

	return &empty.Empty{}, nil
}

func updateToken(accountID string) string {
	return "updatetoken:" + accountID
}

func (accountAPI *accountAPIServer) RequestChangePrivateAccount(
	ctx context.Context, req *account.RequestChangePrivateAccountRequest,
) (*account.RequestChangePrivateAccountResponse, error) {
	var err error

	// Validation
	switch {
	case req == nil:
		return nil, errs.NilObject("request change private request")
	case req.Payload == "":
		return nil, errs.MissingField("payload")
	case req.Project == "":
		return nil, errs.MissingField("project")
	case req.FallbackUrl == "":
		return nil, errs.MissingField("fallback url")
	case req.SendMethod == messaging.SendMethod_SEND_METHOD_UNSPECIFIED:
		return nil, errs.WrapMessage(codes.InvalidArgument, "send method is unspecified")
	}

	// GetAccount the user from database
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.
		First(accountDB, "(email=? OR phone=?) AND project_id = ?", req.Payload, req.Payload, req.Project).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		emailOrPhone := func(err error) string {
			if strings.Contains(strings.ToLower(err.Error()), "email") {
				return "email " + req.Payload
			}
			if strings.Contains(strings.ToLower(err.Error()), "phone") {
				return "phone " + req.Payload
			}
			return "email or phone " + req.Payload
		}
		return nil, errs.WrapMessagef(codes.NotFound, "account with %s does not exist", emailOrPhone(err))
	default:
		return nil, errs.FailedToFind("account", err)
	}

	accountID := fmt.Sprint(accountDB.AccountID)

	uniqueNumber := rand.Intn(199999) + 500000

	// Set token with expiration of 5 minutes
	err = accountAPI.RedisDBWrites.Set(
		ctx, updateToken(accountID), uniqueNumber, time.Duration(5*time.Minute),
	).Err()
	if err != nil {
		return nil, errs.RedisCmdFailed(err, "SET")
	}

	// GetAccount jwt
	jwtToken, err := accountAPI.AuthAPI.GenToken(ctx, &auth.Payload{
		ID:           accountID,
		Names:        accountDB.Names,
		EmailAddress: accountDB.Email,
		PhoneNumber:  accountDB.Phone,
		ProjectID:    accountDB.ProjectID,
	}, time.Now().Add(10*time.Minute))
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to generate token")
	}

	link := fmt.Sprintf("%s?jwt=%s&username=%s&passphrase=%d", req.FallbackUrl, jwtToken, firstVal(accountDB.Email, accountDB.Phone), uniqueNumber)

	ctx, cancel := context.WithTimeout(mdutil.AddFromCtx(ctx), 5*time.Second)
	defer cancel()

	var data string
	if req.SendMethod == messaging.SendMethod_EMAIL {
		data = fmt.Sprintf(
			"You requested to change your account password credentials. Click on the following link in order to change your password. <br> <a href=\"%s?jwt=%s&passphrase=%d&username=%s\" target=\"blank\">Change password</a>",
			req.FallbackUrl, jwtToken, uniqueNumber, firstVal(accountDB.Email, accountDB.Phone),
		)
	} else if req.Project != "" {
		data = fmt.Sprintf("Password reset token for %s \n\nReset Token is %d \nExpires in 10 minutes", req.Project, uniqueNumber)
	} else {
		data = fmt.Sprintf("Password reset token is %d \n\nExpires in 10 minutes", uniqueNumber)
	}

	// Create an outgoing context
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(auth.Header(), fmt.Sprintf("Bearer %s", jwtToken)))

	// Send message
	_, err = accountAPI.MessagingClient.SendMessage(ctx, &messaging.SendMessageRequest{
		Message: &messaging.Message{
			UserId:      accountID,
			Title:       "Reset Account Password",
			Data:        data,
			EmailData:   data,
			Link:        link,
			Save:        true,
			Type:        messaging.MessageType_ALERT,
			SendMethods: []messaging.SendMethod{req.SendMethod},
		},
		Sender:          req.GetSender(),
		SmsAuth:         req.GetSmsAuth(),
		SmsCredentialId: req.SmsCredentialId,
		FetchSmsAuth:    req.FetchSmsAuth,
	})
	if err != nil {
		return nil, errs.WrapErrorWithMsg(err, "failed to send message")
	}

	return &account.RequestChangePrivateAccountResponse{
		ResponseMessage: "reset token was sent to " + req.Payload,
		Jwt:             jwtToken,
	}, nil
}

func (accountAPI *accountAPIServer) UpdatePrivateAccount(
	ctx context.Context, req *account.UpdatePrivateAccountRequest,
) (*empty.Empty, error) {
	var err error

	// Validation
	var ID int
	switch {
	case req == nil:
		return nil, errs.NilObject("UpdatePrivateRequest")
	case req.AccountId == "":
		return nil, errs.MissingField("account id")
	case req.PrivateAccount == nil:
		return nil, errs.NilObject("private account")
	case req.ChangeToken == "":
		return nil, errs.MissingField("change token")
	default:
		ID, err = strconv.Atoi(req.AccountId)
		if err != nil {
			return nil, errs.IncorrectVal("account id")
		}
	}

	// Authorization
	_, err = accountAPI.AuthAPI.AuthorizeActorOrGroup(ctx, req.AccountId, accountAPI.AuthAPI.AdminGroups()...)
	if err != nil {
		return nil, err
	}

	// GetAccount the account details from database
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.Select("account_state").First(accountDB, "account_id=?", ID).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", req.AccountId)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.WrapMessage(codes.PermissionDenied, "account not active")
	}

	// Lets get the update token
	token, err := accountAPI.RedisDBWrites.Get(ctx, updateToken(req.AccountId)).Result()
	switch {
	case err == nil:
	case err == redis.Nil:
		return nil, errs.WrapMessage(codes.NotFound, "update token not found")
	default:
		return nil, errs.RedisCmdFailed(err, "get token")
	}

	if token != req.ChangeToken {
		return nil, errs.WrapMessage(codes.InvalidArgument, "token is incorrect")
	}

	// Hash the password if not empty
	if req.PrivateAccount.Password != "" {
		// Passwords must be similar
		if req.PrivateAccount.ConfirmPassword != req.PrivateAccount.Password {
			return nil, errs.WrapMessage(codes.InvalidArgument, "passwords do not match")
		}

		req.PrivateAccount.Password, err = genHash(req.PrivateAccount.Password)
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to generate password hash")
		}
	}

	// Create database model of the new account
	privateDB := &Account{
		SecurityQuestion: req.PrivateAccount.SecurityQuestion,
		SecurityAnswer:   req.PrivateAccount.SecurityAnswer,
		Password:         req.PrivateAccount.Password,
	}

	// Update the model
	err = accountAPI.SQLDBWrites.Model(privateDB).Where("account_id=?", ID).Updates(privateDB).Error
	if err != nil {
		return nil, errs.FailedToUpdate("account", err)
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) UpdatePrivateAccountExternal(
	ctx context.Context, req *account.UpdatePrivateAccountExternalRequest,
) (*empty.Empty, error) {
	var err error

	// Validation
	switch {
	case req == nil:
		return nil, errs.NilObject("UpdatePrivateRequest")
	case req.Jwt == "":
		return nil, errs.MissingField("jwt")
	case req.Username == "":
		return nil, errs.MissingField("username")
	case req.ProjectId == "":
		return nil, errs.MissingField("project_id")
	case req.PrivateAccount == nil:
		return nil, errs.NilObject("private account")
	case req.ChangeToken == "":
		return nil, errs.MissingField("change token")
	default:
	}

	// Validate jwt token from request
	payload, err := accountAPI.AuthAPI.GetPayloadFromJwt(req.Jwt)
	if err != nil {
		return nil, err
	}

	// The username should match payload data
	if payload.EmailAddress != req.Username && payload.PhoneNumber != req.Username && req.Username != payload.ID {
		return nil, errs.WrapMessage(codes.PermissionDenied, "you are not allowed to perform this operation")
	}

	// GetAccount the account details from database
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.Select("account_id,account_state").
		First(accountDB, "(email=? OR phone=?) AND project_id = ?", req.Username, req.Username, req.ProjectId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", req.Username)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.WrapMessage(codes.PermissionDenied, "account not active")
	}

	// Lets get the update token
	token, err := accountAPI.RedisDBWrites.Get(ctx, updateToken(fmt.Sprint(accountDB.AccountID))).Result()
	switch {
	case err == nil:
	case err == redis.Nil:
		return nil, errs.WrapMessage(codes.NotFound, "reset token expired")
	default:
		return nil, errs.RedisCmdFailed(err, "get token")
	}

	if token != req.ChangeToken {
		return nil, errs.WrapMessage(codes.InvalidArgument, "reset token is incorrect")
	}

	// Hash the password if not empty
	if req.PrivateAccount.Password != "" {
		// Passwords must be similar
		if req.PrivateAccount.ConfirmPassword != req.PrivateAccount.Password {
			return nil, errs.WrapMessage(codes.InvalidArgument, "passwords do not match")
		}

		req.PrivateAccount.Password, err = genHash(req.PrivateAccount.Password)
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to generate password hash")
		}
	}

	// Create database model of the new account
	privateDB := &Account{
		SecurityQuestion: req.PrivateAccount.SecurityQuestion,
		SecurityAnswer:   req.PrivateAccount.SecurityAnswer,
		Password:         req.PrivateAccount.Password,
	}

	// Update the model
	err = accountAPI.SQLDBWrites.Model(privateDB).Where("account_id=?", accountDB.AccountID).Updates(privateDB).Error
	if err != nil {
		return nil, errs.FailedToUpdate("account", err)
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) DeleteAccount(
	ctx context.Context, delReq *account.DeleteAccountRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeleteAccountRequest")
	}

	// Authorization
	_, err := accountAPI.AuthAPI.AuthorizeActorOrGroup(ctx, delReq.AccountId, accountAPI.AuthAPI.AdminGroups()...)
	if err != nil {
		return nil, err
	}

	// Check account id is provided
	var ID int
	switch {
	case delReq.AccountId == "":
		return nil, errs.MissingField("AccountID")
	default:
		ID, err = strconv.Atoi(delReq.AccountId)
		if err != nil {
			return nil, errs.WrapMessage(codes.InvalidArgument, "incorrect account id")
		}
	}

	// Get the account details from database
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.Select("account_state").First(accountDB, "account_id=?", ID).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", delReq.AccountId)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.WrapMessage(codes.PermissionDenied, "account is blocked")
	}

	// Soft delete their account
	err = accountAPI.SQLDBWrites.Delete(accountDB, "account_id=?", ID).Error
	if err != nil {
		return nil, errs.FailedToDelete("account", err)
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) GetAccount(
	ctx context.Context, req *account.GetAccountRequest,
) (*account.Account, error) {
	// Request must not be nil
	if req == nil {
		return nil, errs.NilObject("GetAccountRequest")
	}

	// Authorization
	payload, err := accountAPI.AuthAPI.AuthorizeActorOrGroup(ctx, req.AccountId, accountAPI.AuthAPI.AdminGroups()...)
	if err != nil {
		return nil, err
	}

	// Validation
	var ID int
	switch {
	case req.AccountId == "":
		return nil, errs.MissingField("account id")
	case req.UseEmail || req.UsePhone:
	default:
		ID, err = strconv.Atoi(req.AccountId)
		if err != nil {
			return nil, errs.WrapMessage(codes.InvalidArgument, "incorrect accoint id")
		}
	}

	// GetAccount account from database
	accountDB := &Account{}

	if req.Priviledge {
		if accountAPI.AuthAPI.IsAdmin(payload.Group) {
			err = accountAPI.SQLDBWrites.Unscoped().First(accountDB, "account_id=?", ID).Error
		} else {
			err = accountAPI.SQLDBWrites.First(accountDB, "account_id=?", ID).Error
		}
	} else {
		if req.UsePhone {
			err = accountAPI.SQLDBWrites.First(accountDB, "phone=?", req.AccountId).Error
		} else if req.UseEmail {
			err = accountAPI.SQLDBWrites.First(accountDB, "email=?", req.AccountId).Error
		} else {
			err = accountAPI.SQLDBWrites.First(accountDB, "account_id=?", ID).Error
		}
	}
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("account", req.AccountId)
	default:
		return nil, errs.FailedToFind("account", err)
	}

	// Account should not be deleted
	if accountDB.DeletedAt.Valid && !req.Priviledge {
		return nil, errs.DoesExist("account", req.AccountId)
	}

	// Account should not be blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() && !req.Priviledge {
		return nil, errs.WrapMessage(codes.PermissionDenied, "account is blocked")
	}

	accountPB, err := GetAccountPB(accountDB)
	if err != nil {
		return nil, err
	}

	return GetAccountPBView(accountPB, req.GetView()), nil
}

func (accountAPI *accountAPIServer) BatchGetAccounts(
	ctx context.Context, batchReq *account.BatchGetAccountsRequest,
) (*account.BatchGetAccountsResponse, error) {
	return nil, nil
}

func (accountAPI *accountAPIServer) GetLinkedAccounts(
	ctx context.Context, req *account.GetLinkedAccountsRequest,
) (*account.GetLinkedAccountsResponse, error) {
	return nil, nil
}

func (accountAPI *accountAPIServer) ExistAccount(
	ctx context.Context, existReq *account.ExistAccountRequest,
) (*account.ExistAccountResponse, error) {
	// Request must not be nil
	if existReq == nil {
		return nil, errs.NilObject("ExistAccountRequest")
	}

	var (
		projectID = existReq.GetProjectId()
		email     = existReq.GetEmail()
		phone     = existReq.GetPhone()
		err       error
	)

	// Validation
	switch {
	case projectID == "":
		return nil, errs.MissingField("project id")
	case email == "" && phone == "":
		return nil, errs.MissingField("email, phone or external id")
	}

	// Fix phone
	phone = fixPhone(phone)

	accountDB := &Account{}

	// Query for account with email or phone
	err = accountAPI.SQLDBWrites.Select("account_id,email,phone").
		First(accountDB, "(phone=? OR email=?) AND project_id=?", phone, email, projectID).Error
	switch {
	case err == nil:
		existingFields := make([]string, 0)
		if accountDB.Email == email {
			existingFields = append(existingFields, "email")
		}
		if accountDB.Phone == phone {
			existingFields = append(existingFields, "phone")
		}
		// Account exist
		return &account.ExistAccountResponse{
			Exists:         true,
			AccountId:      fmt.Sprint(accountDB.AccountID),
			ExistingFields: existingFields,
		}, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		// Account doesn't exist
		return &account.ExistAccountResponse{
			Exists: false,
		}, nil
	default:
		return nil, errs.FailedToFind("account", err)
	}
}

const defaultPageSize = 100

func (accountAPI *accountAPIServer) ListAccounts(
	ctx context.Context, req *account.ListAccountsRequest,
) (*account.Accounts, error) {
	// Request must not be nil
	switch {
	case req == nil:
		return nil, errs.NilObject("ListRequest")
	case req.PageSize < 0:
		return nil, errs.IncorrectVal("page size")
	}

	// Authenticate the request
	payload, err := accountAPI.AuthAPI.AuthenticateRequestV2(ctx)
	if err != nil {
		return nil, err
	}

	// Parse page size and page token
	pageSize := req.GetPageSize()
	if pageSize > defaultPageSize && !accountAPI.AuthAPI.IsAdmin(payload.Group) {
		pageSize = defaultPageSize
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	var id uint

	// Get last id from page token
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := accountAPI.PaginationHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(
				codes.InvalidArgument, err, "failed to parse page token",
			)
		}
		id = uint(ids[0])
	}

	db := accountAPI.SQLDBWrites.Limit(int(pageSize) + 1).Order("account_id DESC").Clauses(hints.ForceIndex("PRIMARY").ForOrderBy()).Model(&Account{})

	// Apply filter criterias
	db = generateWhereCondition(db, req.GetListCriteria()).Debug()

	// For admins
	for _, group := range accountAPI.AuthAPI.AdminGroups() {
		if payload.Group == group {
			db = db.Unscoped()
			break
		}
	}

	// ID filter
	if id > 0 {
		db = db.Where("account_id<?", id)
	}

	// Apply project filter
	if payload.ProjectID != "" {
		db = db.Where("project_id=?", payload.ProjectID)
		if req.ListCriteria != nil {
			req.ListCriteria.ProjectIds = []string{}
		}
	} else {
		if !accountAPI.AuthAPI.IsAdmin(payload.Group) {
			return nil, errs.WrapMessage(codes.PermissionDenied, "permission denied to fetch all accounts")
		}
	}

	var collectionCount int64

	// Page token
	if pageToken == "" {
		err = db.Count(&collectionCount).Error
		if err != nil {
			return nil, errs.SQLQueryFailed(err, "count")
		}
	}

	accountsDB := make([]*Account, 0, pageSize+1)

	err = db.Find(&accountsDB).Error
	switch {
	case err == nil:
	default:
		return nil, errs.FailedToFind("accounts", err)
	}

	accountsPB := make([]*account.Account, 0, len(accountsDB))
	pageSize2 := int(pageSize)

	for i, accountDB := range accountsDB {
		accountPB, err := GetAccountPB(accountDB)
		if err != nil {
			return nil, err
		}

		if i == pageSize2 {
			break
		}

		accountsPB = append(accountsPB, GetAccountPBView(accountPB, req.GetView()))
		id = accountDB.AccountID
	}

	var token string
	if len(accountsDB) > pageSize2 {
		// Next page token
		token, err = accountAPI.PaginationHasher.EncodeInt64([]int64{int64(id)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate page token")
		}
	}

	return &account.Accounts{
		NextPageToken:   token,
		Accounts:        accountsPB,
		CollectionCount: collectionCount,
	}, nil
}

// Searches for accounts
func (accountAPI *accountAPIServer) SearchAccounts(
	ctx context.Context, req *account.SearchAccountsRequest,
) (*account.Accounts, error) {
	// Request must not be nil
	if req == nil {
		return nil, errs.NilObject("SearchRequest")
	}

	// Authenticate the request
	payload, err := accountAPI.AuthAPI.AuthenticateRequestV2(ctx)
	if err != nil {
		return nil, err
	}

	// For empty queries
	if req.Query == "" {
		return &account.Accounts{
			Accounts: []*account.Account{},
		}, nil
	}

	// Parse page size and page token
	pageSize := req.GetPageSize()
	if pageSize > defaultPageSize && !accountAPI.AuthAPI.IsAdmin(payload.Group) {
		pageSize = defaultPageSize
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	var ID uint

	// Get last id from page token
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := accountAPI.PaginationHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	accountsDB := make([]*Account, 0, pageSize)

	db := accountAPI.SQLDBReads.Limit(int(pageSize)).Order("account_id DESC").Model(&Account{})

	// Apply filter criterias
	db = generateWhereCondition(db, req.GetSearchCriteria())

	// For admins
	for _, group := range accountAPI.AuthAPI.AdminGroups() {
		if payload.Group == group {
			db = db.Unscoped()
			break
		}
	}

	// Apply project project
	if payload.ProjectID != "" {
		db = db.Where("project_id=?", payload.ProjectID)
		if req.SearchCriteria != nil {
			req.SearchCriteria.ProjectIds = []string{}
		}
	} else {
		if !accountAPI.AuthAPI.IsAdmin(payload.Group) {
			return nil, errs.WrapMessage(codes.PermissionDenied, "permission denied to search all accounts")
		}
	}

	var collectionCount int64

	// Page token
	if pageToken == "" {
		err = db.Count(&collectionCount).Error
		if err != nil {
			return nil, errs.SQLQueryFailed(err, "count")
		}
	}

	parsedQuery := dbutil.ParseQuery(req.Query)

	// "names", "email", "phone", "linked_accounts"
	err = db.Find(&accountsDB, "MATCH(names, email, phone, linked_accounts) AGAINST(? IN BOOLEAN MODE)", parsedQuery).
		Error
	switch {
	case err == nil:
	default:
		return nil, errs.FailedToFind("accounts", err)
	}

	accountsPB := make([]*account.Account, 0, len(accountsDB))
	pageSize2 := int(pageSize)

	for i, accountDB := range accountsDB {
		accountPB, err := GetAccountPB(accountDB)
		if err != nil {
			return nil, err
		}

		if pageSize2 == i {
			break
		}

		accountsPB = append(accountsPB, GetAccountPBView(accountPB, req.GetView()))
		ID = accountDB.AccountID
	}

	var token string
	if len(accountsDB) > pageSize2 {
		// Next page token
		token, err = accountAPI.PaginationHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate page token")
		}
	}

	return &account.Accounts{
		NextPageToken:   token,
		Accounts:        accountsPB,
		CollectionCount: collectionCount,
	}, nil
}

func generateWhereCondition(db *gorm.DB, criteria *account.Criteria) *gorm.DB {
	if criteria == nil {
		return db
	}

	// Filter by account state
	switch {
	case criteria.ShowActiveAccounts:
		db = db.Where("account_state = ?", account.AccountState_ACTIVE.String())
	case criteria.ShowInactiveAccounts:
		db = db.Where("account_state = ?", account.AccountState_INACTIVE.String())
	case criteria.ShowBlockedAccounts:
		db = db.Where("account_state = ?", account.AccountState_BLOCKED.String())
	}

	// Filter by gender
	switch {
	case criteria.ShowFemales:
		db = db.Where("gender = ?", account.Account_FEMALE.String())
	case criteria.ShowMales:
		db = db.Where("gender = ?", account.Account_MALE.String())
	}

	// Filter by date
	if criteria.FilterCreationDate {
		nowSecs := time.Now().Unix()
		switch {
		case criteria.CreatedFrom > 0 && criteria.CreatedUntil > 0 && criteria.CreatedFrom < criteria.CreatedUntil:
			db = db.Where(
				"UNIX_TIMESTAMP(created_at) BETWEEN ? AND ?",
				criteria.CreatedFrom, criteria.CreatedUntil,
			)
		case criteria.CreatedUntil > 0:
			db = db.Where(
				"UNIX_TIMESTAMP(created_at) < ?", criteria.CreatedUntil,
			)
		case criteria.CreatedFrom > 0:
			if criteria.CreatedFrom < nowSecs {
				db = db.Where(
					"UNIX_TIMESTAMP(created_at) > ?", criteria.CreatedFrom,
				)
			}
		}
	}

	// Filter by primary_groups
	if len(criteria.Groups) != 0 {
		db = db.Where("primary_group IN (?)", criteria.Groups)
	}

	// Filter by project id
	if len(criteria.ProjectIds) != 0 {
		db = db.Where("project_id IN (?)", criteria.ProjectIds)
	}

	// Filter by phones
	if len(criteria.Phones) != 0 {
		db = db.Where("phone IN (?)", criteria.Phones)
	}

	// Filter by email
	if len(criteria.Emails) != 0 {
		db = db.Where("email IN (?)", criteria.Emails)
	}

	return db
}

func (accountAPI *accountAPIServer) DailyRegisteredUsers(
	ctx context.Context, req *account.DailyRegisteredUsersRequest,
) (*account.CountStats, error) {
	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request")
	case len(req.Dates) == 0:
		return nil, errs.MissingField("dates")
	case req.DateIsRange && len(req.Dates) == 1:
		return nil, errs.WrapMessage(codes.InvalidArgument, "please provide start and end date for date ranges")
	}

	actor, err := accountAPI.AuthAPI.GetJwtPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Filter == nil {
		req.Filter = &account.DailyRegisteredUsersRequest_Filter{
			ProjectIds: []string{actor.ProjectID},
		}
	} else {
		req.Filter.ProjectIds = []string{actor.ProjectID}
	}

	var (
		wg    = &sync.WaitGroup{}
		mu    = &sync.Mutex{} // guards stats
		stats = make([]*account.CountStat, 0, len(req.Dates))
	)

	if req.DateIsRange {
		dates, err := timeutil.GetDateRanges(req.Dates[0], req.Dates[1])
		if err != nil {
			return nil, err
		}

		req.Dates = dates
	}

	for _, date := range req.Dates {
		wg.Add(1)

		go func(date string) {
			defer wg.Done()

			dateTime, err := timeutil.GetDateFromString(date)
			if err != nil {
				accountAPI.Logger.Errorf("FAILED DailyRegisteredUsers: %v", err)
				return
			}

			db := accountAPI.SQLDBReads.Model(&Account{}).Where("created_at BETWEEN ? AND ?", dateTime, dateTime.Add(24*time.Hour)).Debug()
			if req.Filter != nil {
				if len(req.Filter.ProjectIds) != 0 {
					db = db.Where("project_id IN (?)", req.Filter.ProjectIds)
				}
			}

			var count int64

			// Get count of new users
			err = db.Count(&count).Error
			if err != nil {
				accountAPI.Logger.Errorf("FAILED Count: %v", err)
				return
			}

			mu.Lock()
			stats = append(stats, &account.CountStat{
				Date:  date,
				Count: count,
			})
			mu.Unlock()
		}(date)
	}

	// Wait for results
	wg.Wait()

	return &account.CountStats{
		Stats: stats,
	}, nil
}
