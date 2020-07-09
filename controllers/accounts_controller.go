package controllers

import (
	"encoding/json"
	"github.com/nishant01/procard-go-api/models/accounts"
	"github.com/nishant01/procard-go-api/services"
	"github.com/nishant01/procard-go-api/utils/http_utils"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
	"io/ioutil"
	"net/http"
)

var (
	AccountController AccountsControllerInterface = &accountsController{}
)

type AccountsControllerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type accountsController struct {

}

func (a *accountsController) Register(w http.ResponseWriter, r *http.Request) {
	account := &accounts.Account{}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid json body")
		http_utils.HttpUtils.RespondError(w, respErr)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(requestBody, account); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid accounts params")
		http_utils.HttpUtils.RespondError(w, respErr)
		return
	}

	//parseErr := json.NewDecoder(r.Body).Decode(account)
	//
	//if parseErr != nil {
	//	respErr := rest_errors.NewBadRequestError("invalid request")
	//	http_utils.HttpUtils.RespondError(w, respErr)
	//	return
	//}

	result, createErr := services.AccountService.CreateAccount(*account)

	if createErr != nil {
		http_utils.HttpUtils.RespondError(w, createErr)
		return
	}

	http_utils.HttpUtils.SetTokenCookie(w, result.Token)

	http_utils.HttpUtils.RespondJson(w, http.StatusCreated, result)
}

func (a *accountsController) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := &accounts.LoginRequest{}

	parseErr := json.NewDecoder(r.Body).Decode(loginRequest)

	if parseErr != nil {
		respErr := rest_errors.NewBadRequestError("invalid request")
		http_utils.HttpUtils.RespondError(w, respErr)
		return
	}

	result, loginErr := services.AccountService.LoginUser(*loginRequest)

	if loginErr != nil {
		http_utils.HttpUtils.RespondError(w, loginErr)
		return
	}

	http_utils.HttpUtils.SetTokenCookie(w, result.Token)

	http_utils.HttpUtils.RespondJson(w, http.StatusOK, result)
}
