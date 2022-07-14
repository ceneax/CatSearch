package controller

import (
	"CatSearch/config"
	"CatSearch/engine"
	"CatSearch/model"
	"CatSearch/net"
	"CatSearch/util"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Translate(context *gin.Context) {
	type transReq struct {
		Source string `json:"source"`
		Type   string `json:"type"`
	}

	var req transReq
	if err := context.ShouldBindJSON(&req); err != nil {
		model.ApiFail(context, 100, "参数有误")
		return
	}

	model.ApiSuccess(context, util.Translate(req.Source, req.Type))

	//type transRes struct {
	//	Isdict     int     `json:"isdict"`
	//	Confidence float64 `json:"confidence"`
	//	Target     string  `json:"target"`
	//	Rc         int     `json:"rc"`
	//}
	//
	//var req transReq
	//if err := context.ShouldBindJSON(&req); err != nil {
	//	model.ApiFail(context, 100, "参数有误")
	//	return
	//}
	//
	//headers := http.Header{}
	//headers.Set("Content-Type", "application/json")
	//headers.Set("x-authorization", "token 9sdftiq37bnv410eon2l")
	//
	//res, err := net.Post("https://api.interpreter.caiyunai.com/v1/translator", headers, map[string]interface{}{
	//	"source":     req.Source,
	//	"trans_type": req.Type, // auto2en
	//	"request_id": strconv.FormatInt(time.Now().Unix(), 10),
	//	"detect":     true,
	//}, "", time.Duration(config.EngineTimeout)*time.Millisecond)
	//if err != nil {
	//	model.ApiFail(context, 101, "翻译失败")
	//	return
	//}
	//
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		return
	//	}
	//}(res.Body)
	//
	//responseBody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	model.ApiFail(context, 102, "翻译失败")
	//	return
	//}
	//
	//var body transRes
	//err = json.Unmarshal(responseBody, &body)
	//if err != nil {
	//	model.ApiFail(context, 103, "翻译失败")
	//	return
	//}
	//
	//model.ApiSuccess(context, body.Target)
}

func CodeAnswer(context *gin.Context) {
	kw := context.Query("q")
	if len(kw) <= 0 {
		model.ApiFail(context, 200, "参数有误")
		return
	}

	headers := http.Header{}
	headers.Set("User-Agent", config.HEADER_USER_AGENT)
	headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	proxy := ""
	if config.Debug {
		proxy = "http://127.0.0.1:10809"
	}
	res, err := net.Get(
		"https://www.google.com/search?q="+url.QueryEscape(util.Translate(kw, "auto2en")+" site:stackoverflow.com")+"&ie=utf8&oe=utf8&start=0",
		headers, proxy, time.Duration(config.EngineTimeout)*time.Millisecond)
	if err != nil {
		model.ApiFail(context, 201, "获取失败")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		model.ApiFail(context, 202, "获取失败")
		return
	}

	resUrl := ""

	doc.Find("#search .g").Each(func(i int, selection *goquery.Selection) {
		if resUrl != "" {
			return
		}
		if selection.Find(".g").Length() <= 0 && selection.Find(".yuRUbf").Length() > 0 {
			resUrl = selection.Find(".yuRUbf a").First().AttrOr("href", "")
			return
		}
	})

	if resUrl == "" {
		model.ApiFail(context, 203, "未获取到相关内容")
		return
	}

	parse, err := url.Parse(resUrl)
	if err != nil {
		model.ApiFail(context, 204, "获取失败")
		return
	}

	splitRes := strings.Split(parse.Path, "/")
	if len(splitRes) < 3 || splitRes[2] == "" {
		model.ApiFail(context, 205, "获取失败")
		return
	}

	stackoverflowUrl := "https://api.stackexchange.com/2.3/questions/" + splitRes[2] + "?key=U4DMV*8nvpm3EOpvf69Rxw((&site=stackoverflow&order=desc&sort=votes&filter=!6VvPDzQ)wlg1u"
	res2, err := net.Get(stackoverflowUrl, http.Header{}, proxy, time.Duration(config.EngineTimeout)*time.Millisecond)
	if err != nil {
		model.ApiFail(context, 206, "获取失败")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res2.Body)

	responseBody, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		model.ApiFail(context, 207, "获取失败")
		return
	}

	var body engine.StackoverflowRes
	err = json.Unmarshal(responseBody, &body)
	if err != nil {
		model.ApiFail(context, 208, "获取失败")
		return
	}

	if len(body.Items) <= 0 {
		model.ApiFail(context, 209, "未获取到相关内容")
		return
	}

	item := body.Items[0]

	if !item.Is_answered || len(item.Answers) <= 0 {
		model.ApiFail(context, 210, "未获取到相关内容")
		return
	}

	var answer engine.Answers
	for _, a := range item.Answers {
		if a.Is_accepted {
			answer = a
			break
		}
	}

	model.ApiSuccess(context, model.CodeAnswer{
		Url:   item.Link,
		Title: "stackoverflow - " + item.Title,
		Ans: template.HTML(strings.ReplaceAll(strings.ReplaceAll(answer.Body, "<pre><code>",
			"<pre class=\"code-answer-pre\"><code>"), "</code></pre>", "</code></pre>")),
		Tags: item.Tags,
	})
}
