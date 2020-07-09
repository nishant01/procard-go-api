package controllers

import (
	"github.com/nishant01/procard-go-api/utils/http_utils"
	"net/http"
)

const (
	pong = "pong"
)

var (
	PingController PingControllerInterface = &pingController{}
)

type PingControllerInterface interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type pingController struct {

}

func (c *pingController) Ping(w http.ResponseWriter, r *http.Request) {
	http_utils.HttpUtils.RespondJson(w, http.StatusOK, pong)
}
