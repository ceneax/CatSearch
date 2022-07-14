package search

import (
	"CatSearch/config"
	"CatSearch/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DoResponse(context *gin.Context, p model.Params, result model.Result) {
	//model.Success(context, *result["google"])

	context.HTML(http.StatusOK, "search.html", gin.H{
		"params":        p,
		"tags":          config.Seg.CutTrim(p.Query, true),
		"hasAnswer":     result.Answer.Ans != nil,
		"hasCard":       result.Card.Content != "",
		"hasCodeAnswer": result.CodeAnswer.Ans != "",
		"result":        result,
	})
}
