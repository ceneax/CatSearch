package engine

import (
	"CatSearch/config"
	"CatSearch/model"
	"CatSearch/net"
	"CatSearch/util"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type StackoverflowRes struct {
	Items           []Items `json:"items"`
	Has_more        bool    `json:"has_more"`
	Quota_max       int     `json:"quota_max"`
	Quota_remaining int     `json:"quota_remaining"`
}

type Items struct {
	Is_answered        bool      `json:"is_answered"`
	Accepted_answer_id int       `json:"accepted_answer_id"`
	Answer_count       int       `json:"answer_count"`
	Score              int       `json:"score"`
	Question_id        int       `json:"question_id"`
	Link               string    `json:"link"`
	Title              string    `json:"title"`
	Body               string    `json:"body"`
	Tags               []string  `json:"tags"`
	Answers            []Answers `json:"answers"`
}

type Answers struct {
	Is_accepted     bool   `json:"is_accepted"`
	Score           int    `json:"score"`
	Answer_id       int    `json:"answer_id"`
	Question_id     int    `json:"question_id"`
	Content_license string `json:"content_license"`
	Body            string `json:"body"`
}

type stackoverflow struct {
	name string
}

func (e *stackoverflow) GetName() string {
	return e.name
}

func (e *stackoverflow) Request(p model.Params) Request {
	page := 1
	intPage, err := strconv.Atoi(p.Page)
	if err == nil {
		page = intPage
	}
	if page > 1 {
		return Request{}
	}

	headers := http.Header{}
	headers.Set("User-Agent", config.HEADER_USER_AGENT)
	headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	proxy := ""
	if config.Debug {
		proxy = "http://127.0.0.1:10809"
	}
	res, err := net.Get(
		"https://www.google.com/search?q="+url.QueryEscape(util.Translate(p.Query, "auto2en")+" site:stackoverflow.com")+"&ie=utf8&oe=utf8&start=0",
		headers, proxy, time.Duration(config.EngineTimeout)*time.Millisecond)
	if err != nil {
		log.Println("内部错误！" + err.Error())
		return Request{}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(e.name + " Response Body 关闭异常！" + err.Error())
			return
		}
	}(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(e.name + " 解析出错了！" + err.Error())
		return Request{}
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
		return Request{}
	}

	parse, err := url.Parse(resUrl)
	if err != nil {
		return Request{}
	}

	splitRes := strings.Split(parse.Path, "/")
	if len(splitRes) < 3 || splitRes[2] == "" {
		return Request{}
	}

	return Request{
		Url: "https://api.stackexchange.com/2.3/questions/" + splitRes[2] + "?key=U4DMV*8nvpm3EOpvf69Rxw((&site=stackoverflow&order=desc&sort=votes&filter=!6VvPDzQ)wlg1u",
	}

	//page := 1
	//intPage, err := strconv.Atoi(p.Page)
	//if err == nil {
	//	page = intPage
	//}
	//if page > 1 {
	//	return Request{}
	//}
	//
	//q := util.Translate(p.Query, "auto2en")
	//
	//return Request{
	//	Url: "https://api.stackexchange.com/2.3/search/advanced?q=" + url.QueryEscape(q) + "&page=1&pagesize=1&site=stackoverflow&sort=votes" +
	//		"&order=desc&key=U4DMV*8nvpm3EOpvf69Rxw((&accepted=True&filter=!*MZqiH2o(M2z5E0D",
	//}
}

func (e *stackoverflow) Response(res *http.Response) Response {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(e.name + " Response Body 关闭异常！" + err.Error())
			return
		}
	}(res.Body)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(e.name + " Response Body 转换异常！" + err.Error())
		return Response{}
	}

	var body StackoverflowRes
	err = json.Unmarshal(responseBody, &body)
	if err != nil {
		log.Println(e.name + " Response Body 转换json异常！" + err.Error())
		return Response{}
	}

	if len(body.Items) <= 0 {
		return Response{}
	}

	item := body.Items[0]

	if !item.Is_answered || len(item.Answers) <= 0 {
		return Response{}
	}

	var answer Answers
	for _, a := range item.Answers {
		if a.Is_accepted {
			answer = a
			break
		}
	}

	return Response{
		Result: model.Result{
			CodeAnswer: model.CodeAnswer{
				Url:   item.Link,
				Title: "stackoverflow - " + item.Title,
				Ans: template.HTML(strings.ReplaceAll(strings.ReplaceAll(answer.Body, "<pre><code>",
					"<div class=\"code-answer-pre\"><code>"), "</code></pre>", "</code></div>")),
				Tags: item.Tags,
			},
		},
	}
}
