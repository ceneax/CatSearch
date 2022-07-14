package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiSuccess(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, ApiResponse{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

func ApiFail(context *gin.Context, code int, msg string) {
	context.JSON(http.StatusOK, ApiResponse{
		Code: code,
		Msg:  msg,
	})
}
