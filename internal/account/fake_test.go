package account

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/services/pkg/api/account"
)

func fakePhone() string {
	phone := randomdata.PhoneNumber()
	if len(phone) > 10 {
		phone = phone[:10]
	}
	return fmt.Sprintf("%s%d", phone, randomdata.Number(1000, 9999))
}

func createAdmin(accountState account.AccountState) (string, error) {
	accountPB := fakeAccount()
	accountPB.Group = auth.DefaultAdminGroup()
	accountPB.State = accountState

	// Get admin model
	accountDB, err := GetAccountDB(accountPB)
	if err != nil {
		return "", err
	}

	// Save to database
	err = AccountAPIServer.SQLDBWrites.Create(accountDB).Error
	if err != nil {
		return "", err
	}

	// Return account ID
	return fmt.Sprint(accountDB.AccountID), nil
}

// creates a fake account
func fakeAccount() *account.Account {
	// randPayload := randomdata.RandStringRunes(10)
	return &account.Account{
		ProjectId:   "1",
		AccountId:   randomdata.RandStringRunes(32),
		Email:       randomdata.Email(),
		Phone:       fakePhone(),
		Names:       randomdata.SillyName(),
		BirthDate:   randomdata.FullDate(),
		Gender:      account.Account_Gender(randomdata.Number(1, len(account.Account_Gender_name))),
		Nationality: randomdata.Country(randomdata.FullCountry),
		ProfileUrl:  randomdata.MacAddress(),
		State:       account.AccountState_ACTIVE,
		Group:       auth.DefaultUserGroup(),
	}
}

// create a fake account private profile
func fakePrivateAccount() *account.PrivateAccount {
	return &account.PrivateAccount{
		Password:         "hakty11",
		ConfirmPassword:  "hakty11",
		SecurityQuestion: "What is your pets name",
		SecurityAnswer:   randomdata.SillyName(),
	}
}
