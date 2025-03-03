package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/gidyon/micro/v2/utils/errs"
	"github.com/gidyon/micro/v2/utils/mdutil"
	"github.com/gidyon/micro/v2/utils/templateutil"
	"github.com/gidyon/services/pkg/api/account"
	"github.com/gidyon/services/pkg/api/messaging"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func (accountAPI *accountAPIServer) validateAdminUpdateAccountRequest(
	ctx context.Context, changeReq *account.AdminUpdateAccountRequest,
) (*Account, error) {
	// Request should not be nil
	if changeReq == nil {
		return nil, errs.NilObject("AdminUpdateAccountRequest")
	}

	// Authorize the admin
	_, err := accountAPI.AuthAPI.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}

	accountID := changeReq.GetAccountId()
	adminID := changeReq.GetAdminId()

	// Validation
	var ID, ID2 int
	switch {
	case accountID == "":
		return nil, errs.MissingField("account id")
	case adminID == "":
		return nil, errs.MissingField("admin id")
	case changeReq.UpdateOperation == account.UpdateOperation_UPDATE_OPERATION_INSPECIFIED:
		return nil, errs.WrapMessage(codes.InvalidArgument, "update operation is uknown")
	default:
		ID, err = strconv.Atoi(adminID)
		if err != nil {
			return nil, errs.WrapMessage(codes.InvalidArgument, "incorrect admin id")
		}
		ID2, err = strconv.Atoi(accountID)
		if err != nil {
			return nil, errs.WrapMessage(codes.InvalidArgument, "incorrect account id")
		}
	}

	// Get admin
	admin := &Account{}
	err = accountAPI.SQLDBWrites.Unscoped().Select("account_state,primary_group").
		First(admin, "account_id=?", ID).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.WrapMessagef(codes.NotFound, "admin account with id: %s doesn't exist", accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	// Admin account must be active
	if admin.AccountState != account.AccountState_ACTIVE.String() {
		return nil, errs.WrapMessage(codes.PermissionDenied, "admin account not active")
	}

	// Admin must be admin
	if !accountAPI.AuthAPI.IsAdmin(admin.PrimaryGroup) {
		return nil, errs.WrapMessage(codes.PermissionDenied, "only admins allowed")
	}

	// Get user
	accountDB := &Account{}
	err = accountAPI.SQLDBWrites.Unscoped().First(accountDB, "account_id=?", ID2).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.WrapMessagef(codes.NotFound, "account with id: %s doesn't exist", accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	return accountDB, nil
}

func (accountAPI *accountAPIServer) AdminUpdateAccount(
	ctx context.Context, req *account.AdminUpdateAccountRequest,
) (*empty.Empty, error) {
	// Validate the request, super admin credentials and account owner
	accountDB, err := accountAPI.validateAdminUpdateAccountRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var (
		fullName    = accountDB.Names
		messageType messaging.MessageType
		title       string
		data        string
		link        string
	)

	// Start a transaction
	tx := accountAPI.SQLDBWrites.Begin()
	defer func() {
		if err := recover(); err != nil {
			accountAPI.Logger.Errorln(err)
		}
	}()

	if tx.Error != nil {
		return nil, errs.FailedToBeginTx(err)
	}

	// Update the model
	tx = tx.Unscoped().Model(&Account{}).Where("account_id=?", req.AccountId)

	switch req.UpdateOperation {
	case account.UpdateOperation_UNDELETE:
		err = tx.Update("deleted_at", nil).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to delete account")
		}
		messageType = messaging.MessageType_INFO
		title = "Your Account Has Been Restored"
		data = fmt.Sprintf(
			"Hello %s, we are glad to inform you that your account has been restored", fullName,
		)

	case account.UpdateOperation_DELETE:
		err = tx.Update("deleted_at", time.Now()).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to undelete account")
		}
		messageType = messaging.MessageType_ALERT
		title = "Your Account Has Been Deleted"
		data = fmt.Sprintf(
			"Hello %s, we are sad to inform you that your account has been scheduled for deletion", fullName,
		)

	case account.UpdateOperation_UNBLOCK:
		// The state must be blocked in order to unblock it
		if accountDB.AccountState != account.AccountState_BLOCKED.String() {
			tx.Rollback()
			return nil, errs.WrapMessage(codes.FailedPrecondition, "account is not blocked")
		}
		err = tx.Update("account_state", account.AccountState_ACTIVE.String()).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to unblock account")
		}
		messageType = messaging.MessageType_INFO
		title = "Your Account Has Been Unblock"
		data = fmt.Sprintf(
			"Hello %s, we are glad to inform you that your account has been unblocked", fullName,
		)

	case account.UpdateOperation_BLOCK:
		// The state must be active in order to block it
		if accountDB.AccountState != account.AccountState_ACTIVE.String() {
			tx.Rollback()
			return nil, errs.WrapMessage(codes.FailedPrecondition, "account is not active")
		}
		err = tx.Update("account_state", account.AccountState_BLOCKED.String()).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to block account")
		}
		messageType = messaging.MessageType_ALERT
		title = "Your Account Has Been Blocked"
		data = fmt.Sprintf(
			"Hello %s, we are sad to inform you that your account has been blocked", fullName,
		)

	case account.UpdateOperation_CHANGE_GROUP:
		// The state must be active in order to change group it
		if accountDB.AccountState != account.AccountState_ACTIVE.String() {
			tx.Rollback()
			return nil, errs.WrapMessage(codes.FailedPrecondition, "account to change group is not active")
		}
		// Update the model
		bs, err := json.Marshal(req.Payload)
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to json unmarshal")
		}
		err = tx.Update("secondary_groups", bs).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to update secondary groups")
		}
		messageType = messaging.MessageType_INFO
		title = "Your Account Has Group Has Been Changed"
		data = fmt.Sprintf(
			"Hello %s, your account has been added to the following groups %s", fullName, req.Payload,
		)

	case account.UpdateOperation_CHANGE_PRIMARY_GROUP:
		// The state must be active in order to change group it
		if accountDB.AccountState != account.AccountState_ACTIVE.String() {
			tx.Rollback()
			return nil, errs.WrapMessage(codes.FailedPrecondition, "account to change group is not active")
		}
		if len(req.Payload) == 0 {
			return nil, errs.WrapMessage(codes.InvalidArgument, "missing group")
		}
		group := req.Payload[0]
		// Update the model
		err = tx.Update("primary_group", group).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to update primary group")
		}
		messageType = messaging.MessageType_INFO
		title = fmt.Sprintf("Your Group Has Been Changed To %s", group)
		data = fmt.Sprintf(
			"Hello %s, your account group has been updated to %s", fullName, group,
		)

	case account.UpdateOperation_ADMIN_ACTIVATE:
		err = tx.Update("account_state", account.AccountState_ACTIVE.String()).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to activate account")
		}
		messageType = messaging.MessageType_INFO
		title = "Your Account Has Been Activated"
		data = fmt.Sprintf(
			"Hello %s. Your account has been activated by the administrator", fullName,
		)

	case account.UpdateOperation_PASSWORD_RESET:
		newPass, err := genHash(strings.Join(req.Payload, ""))
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to generate hash password")
		}
		err = tx.Update("password", newPass).Error
		if err != nil {
			tx.Rollback()
			return nil, errs.WrapErrorWithMsg(err, "failed to activate account")
		}
		messageType = messaging.MessageType_INFO
		title = "Your Account Pasword Has Been Updated"
		data = fmt.Sprintf(
			"Hello %s. Your account password has been updated by the administrator. <br>New password is: %s", fullName, strings.Join(req.Payload, ""),
		)
	}

	if req.Notify {
		// Email template
		emailContent := templateutil.EmailData{
			Names:        accountDB.Names,
			AccountID:    req.AccountId,
			AppName:      firstVal(req.GetSender().GetAppName(), accountAPI.AppName),
			Reason:       req.Reason,
			TemplateName: templateName,
		}

		var emailData string
		if accountAPI.tpl != nil {
			content := bytes.NewBuffer(make([]byte, 0, 64))
			err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
			if err != nil {
				tx.Rollback()
				return nil, errs.WrapErrorWithMsg(err, "failed to exucute template")
			}
			emailData = content.String()
		} else {
			emailData = data
		}

		ctx, cancel := context.WithTimeout(mdutil.AddFromCtx(ctx), 5*time.Second)
		defer cancel()

		// Send message to inform necessary audience
		_, err = accountAPI.MessagingClient.SendMessage(ctx, &messaging.SendMessageRequest{
			Message: &messaging.Message{
				UserId:      req.AccountId,
				Title:       title,
				Data:        data,
				EmailData:   emailData,
				Link:        link,
				Save:        true,
				Type:        messageType,
				SendMethods: []messaging.SendMethod{messaging.SendMethod_SMSV2, messaging.SendMethod_EMAIL},
			},
			Sender:          req.GetSender(),
			SmsAuth:         req.GetSmsAuth(),
			SmsCredentialId: req.SmsCredentialId,
			FetchSmsAuth:    req.FetchSmsAuth,
		}, grpc.WaitForReady(true))
		if err != nil {
			accountAPI.Logger.Errorf("error while sending account changed message: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, errs.WrapErrorWithMsg(err, "failed to commit transation")
	}

	return &empty.Empty{}, nil
}
