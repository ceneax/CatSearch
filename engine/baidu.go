package engine

import (
	"CatSearch/config"
	"CatSearch/model"
	"CatSearch/util"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type baiduRes struct {
	Feed feed `json:"feed"`
}

type feed struct {
	Entry []entry `json:"entry"`
}

type entry struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Abs   string `json:"abs"`
}

type baidu struct {
	name string
}

func (e *baidu) GetName() string {
	return e.name
}

func (e *baidu) Request(p model.Params) Request {
	headers := http.Header{}
	headers.Set("User-Agent", config.HEADER_USER_AGENT)
	headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	page := 1
	intPage, err := strconv.Atoi(p.Page)
	if err == nil {
		page = intPage
	}

	return Request{
		Url:     "https://www.baidu.com/s?ie=UTF-8&tn=json&wd=" + url.QueryEscape(p.Query) + "&pn=" + strconv.Itoa((page-1)*10),
		Headers: headers,
	}
}

func (e *baidu) Response(res *http.Response) Response {
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

	var body baiduRes
	err = json.Unmarshal(responseBody, &body)
	if err != nil {
		log.Println(e.name + " Response Body 转换json异常！" + err.Error())
		return Response{}
	}

	result := model.Result{}

	for _, item := range body.Feed.Entry {
		if item.Url == "" || item.Title == "" {
			continue
		}
		result.Results = append(result.Results, model.Results{
			Url:     util.ParseUrlHostPath(item.Url),
			Title:   template.HTML(item.Title),
			Content: template.HTML(item.Abs),
			Source:  e.name,
		})
	}

	return Response{
		Result: result,
	}
}
