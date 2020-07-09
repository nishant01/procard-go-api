package http_utils

import (
	"encoding/json"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
	"net/http"
	"time"
)

var (
	HttpUtils HttpUtilsInterface = &httpUtils{}
)

type HttpUtilsInterface interface {
	RespondJson(http.ResponseWriter, int, interface{})
	RespondError(http.ResponseWriter, rest_errors.RestErr)
	SetTokenCookie(http.ResponseWriter, string)
}

type httpUtils struct {

}

func (h *httpUtils) RespondJson(w http.ResponseWriter, statusCode int, body interface{}) {
	setupResponse(&w)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func (h *httpUtils) RespondError(w http.ResponseWriter, err rest_errors.RestErr) {
	h.RespondJson(w, err.Status(), err)
}

func (h *httpUtils) SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "procard_token",
		Value:   token,
		Expires: time.Now().Add(time.Minute * 10),
	})
}

func setupResponse(w *http.ResponseWriter) {
	//(*w).Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	//(*w).Header().Set("X-Frame-Options", "DENY")
	//(*w).Header().Set("X-Content-Type-Options", "nosniff")
	//(*w).Header().Set("Pragma", "no-cache")
	//
	//(*w).Header().Set("X-XSS-Protection", "1; mode=block")
	//(*w).Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT, DELETE, PATCH")
	//(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Expose-Headers, Access-Control-Allow-Origin, API-Key,Content-Type,If-Modified-Since,Cache-Control, Accept-Encoding, X-CSRF-Token, Authorization, Content-Length, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}


