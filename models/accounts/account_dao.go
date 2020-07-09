package accounts

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/nishant01/procard-go-api/config/database"
	"github.com/nishant01/procard-go-api/utils/logger"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
)

var (
	db = database.GetDB()
)

func (a *Account) dbValidation() rest_errors.RestErr {
	//Email and Username must be unique
	temp := &Account{}

	err := db.Table("accounts").Where("email = ? OR username = ?", a.Email, a.Username).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("Error when tying to connect database", err)
		return rest_errors.NewInternalServerError("Connection error. Please retry", errors.New("database error"))
	}

	if temp.Email != "" && temp.Email == a.Email {
		logger.Error("Email already in use by another user", err)
		return rest_errors.NewInternalServerError("Email already in use by another user.", errors.New("database error"))
	}

	if temp.Username != "" && temp.Username == a.Username{
		logger.Error("Username already in use by another user", err)
		return rest_errors.NewInternalServerError("Username already in use by another user.", errors.New("database error"))
	}

	return nil
}

func (a *Account) Save() rest_errors.RestErr {
	//check for errors and duplicate emails and username
	dbValidationErr := a.dbValidation()

	if dbValidationErr != nil {
		return dbValidationErr
	}

	createdAccount := db.Create(&a)

	// close database when not in use
	// defer db.Close()

	if createdAccount.Error != nil || a.ID <= 0 {
		logger.Error("Error when trying to save account", createdAccount.Error)
		return rest_errors.NewInternalServerError("Unable to create account this time, please retry", errors.New("database error"))
	}

	return nil
}

func (a *Account) FindByEmailAndPassword(email string, password string) rest_errors.RestErr {
	dbErr := db.Table("accounts").Where("email = ?", email).First(&a).Error

	// defer db.Close()

	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			logger.Error("Email address not found", dbErr)
			return rest_errors.NewInternalServerError("Email address not found", errors.New("database error"))
		}
		logger.Error("Connection error", dbErr)
		return rest_errors.NewInternalServerError("Connection error. Please retry", errors.New("database error"))
	}

	err := a.CheckPasswordHash(a.Password, password)

	if err != nil {
		logger.Error("Invalid password.", err)
		return err
	}

	a.Password = ""

	return nil
}
