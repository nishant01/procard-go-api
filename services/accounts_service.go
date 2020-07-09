package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nishant01/procard-go-api/models/accounts"
	"time"

	// "github.com/nishant01/procard-go-api/utils/crypto_utils"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
)

const (
	tokenKey = "thisIsTheJwtSecretPassword" //os.Getenv("token_key")
)

var (
	AccountService AccountServiceInterface = &accountService{}
)

type AccountServiceInterface interface {
	CreateAccount(accounts.Account) (*accounts.Account, rest_errors.RestErr)
	LoginUser(accounts.LoginRequest) (*accounts.Account, rest_errors.RestErr)
}

type accountService struct {}

func (a *accountService) CreateAccount(account accounts.Account) (*accounts.Account, rest_errors.RestErr) {
	if err := account.Validate(); err != nil {
		return nil, err
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	hashedPassword, err := account.HashPassword(account.Password)

	if err != nil {
		fmt.Println(err)
		encrErr := rest_errors.NewInternalServerError(fmt.Sprintf("Password Encryption failed"), errors.New("error processing json"))
		//http_utils.HttpUtils.RespondError(w, passEncrErr)
		return nil, encrErr
	}

	account.Password = string(hashedPassword)
	//account.Password = crypto_utils.GetMd5(account.Password)
	account.IsActive = accounts.IsActive
	account.IsConfigured = accounts.IsConfigured

	if err := account.Save(); err != nil {
		return nil, err
	}

	//Create new JWT token for the newly registered account
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	tk := &accounts.Token{
		UserID: account.ID,
		Username: account.Username,
		Email: account.Email,
		IsActive: account.IsActive,
		IsConfigured: account.IsConfigured,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(tokenKey))

	if err != nil {
		tokenErr := rest_errors.NewInternalServerError("Token generation error", errors.New("token error"))
		return nil, tokenErr
	}

	account.Token = tokenString
	account.Password = "" //delete password
	//var resp = map[string]interface{}{"status": false, "message": "logged in"}
	//resp["token"] = tokenString //Store the token in the response
	//resp["account"] = account

	return &account, nil
}

func (a *accountService) LoginUser(loginRequest accounts.LoginRequest) (*accounts.Account, rest_errors.RestErr) {
	account := &accounts.Account{
		Email: loginRequest.Email,
		Password: loginRequest.Password,
	}

	if err := account.FindByEmailAndPassword(loginRequest.Email, loginRequest.Password); err != nil {
		return nil, err
	}

	//Create JWT token
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	tk := &accounts.Token{
		UserID: account.ID,
		Username: account.Username,
		Email: account.Email,
		IsActive: account.IsActive,
		IsConfigured: account.IsConfigured,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(tokenKey))

	if err != nil {
		tokenErr := rest_errors.NewInternalServerError("token generation error", errors.New("token error"))
		return nil, tokenErr
	}

	account.Token = tokenString //Store the token in the response

	//var resp = map[string]interface{}{"status": false, "message": "logged in"}
	//resp["token"] = tokenString //Store the token in the response
	//resp["account"] = account

	return account, nil
}
