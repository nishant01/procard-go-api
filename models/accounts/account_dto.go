package accounts

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	IsActive = true
	IsConfigured = false
)

//Token struct declaration
type Token struct {
	UserID	 			uint	`json:"user_id"`
	Username			string	`json:"username"`
	Email  				string	`json:"email"`
	IsActive			bool	`json:"is_active"`
	IsConfigured		bool	`json:"is_configured"`
	*jwt.StandardClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string	`json:"username"`
}

type Account struct {
	gorm.Model
	Username		string	`json:"username"`
	Email 			string	`json:"email"`
	Password 		string	`json:"password"`
	Token 			string 	`json:"token";sql:"-"`
	IsActive		bool	`json:"is_active"`
	IsConfigured	bool	`json:"is_configured"`
}

type Accounts []Account

func (a *Account) Validate() rest_errors.RestErr {
	// account.FirstName = strings.TrimSpace(account.FirstName)
	// account.LastName = strings.TrimSpace(account.LastName)

	a.Username = strings.TrimSpace(strings.ToLower(a.Username))
	if a.Username == "" {
		return rest_errors.NewBadRequestError("Username is required")
	}

	if strings.Contains(a.Username, "@") {
		return rest_errors.NewBadRequestError( "Invalid username")
	}

	a.Email = strings.TrimSpace(strings.ToLower(a.Email))
	if a.Email == "" {
		return rest_errors.NewBadRequestError("Email address is required")
	}

	if !strings.Contains(a.Email, "@") {
		return rest_errors.NewBadRequestError( "Invalid email address")
	}

	a.Password = strings.TrimSpace(a.Password)
	if a.Password == "" {
		return rest_errors.NewBadRequestError("Password is required")
	}

	if len(a.Password) < 6 {
		return rest_errors.NewBadRequestError("Password length should be minimum 6 character")
	}

	return nil
}

// CheckPasswordHash checks password hash and password from user input if they match
func (a *Account) CheckPasswordHash(hash, password string) rest_errors.RestErr {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			return rest_errors.NewInternalServerError("Incorrect password.", errors.New("database error"))
		}
		return rest_errors.NewBadRequestError("Incorrect password")
	}
	return nil
}

// HashPassword hashes password from user input
func (a *Account) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
