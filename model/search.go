package model

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type Params struct {
	Query string
	Page  string
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Result `json:"data"`
}

type Result struct {
	Answer     Answer     `json:"answer"`
	Card       Card       `json:"card"`
	Results    []Results  `json:"results"`
	CodeAnswer CodeAnswer `json:"codeAnswer"`
}

type CodeAnswer struct {
	Url   string        `json:"url"`
	Title string        `json:"title"`
	Ans   template.HTML `json:"ans"`
	Tags  []string      `json:"tags"`
}

type Answer struct {
	Url   string   `json:"url"`
	Title string   `json:"title"`
	Ans   []string `json:"ans"`
}

type Results struct {
	Url     Url           `json:"url"`
	Title   template.HTML `json:"title"`
	Content template.HTML `json:"content"`
	Source  string        `json:"source"`
	Blocked bool          `json:"blocked"`
}

type Url struct {
	Raw   string `json:"raw"`
	Parse Parse  `json:"parse"`
}

type Parse struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type Card struct {
	Title    string            `json:"title"`
	Category string            `json:"category"`
	Content  string            `json:"content"`
	Info     map[string]string `json:"info"`
}

func Success(context *gin.Context, data Result) {
	context.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

func Fail(context *gin.Context, code int, msg string) {
	context.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}
