package util

import (
	"CatSearch/config"
	"CatSearch/model"
	"CatSearch/net"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ParseUrlHostPath(rawUrl string) model.Url {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return model.Url{}
	}

	scheme, host, path := "", uri.Host, uri.Path

	if uri.Scheme == "" {
		scheme = "http"
	} else {
		scheme = uri.Scheme
	}

	if path != "" {
		if path[len(path)-1:] == "/" {
			path = path[:len(path)-1]
		}
	}

	return model.Url{
		Raw: rawUrl,
		Parse: model.Parse{
			Host: scheme + "://" + host,
			Path: strings.ReplaceAll(path, "/", " › "),
		},
	}
}

func ContentHighlight(content string, q string) template.HTML {
	for _, s := range config.Seg.Cut(q, true) {
		if strings.ReplaceAll(s, " ", "") == "" || IsSpecialString(s) {
			continue
		}
		//re := regexp.MustCompile("(?i)" + s)
		//content = re.ReplaceAllString(content, "<span class=\"highlight-style\">"+s+"</span>")
		//content = strings.ReplaceAll(content, s, "<span class=\"highlight-style\">"+s+"</span>")
		//if strings.Title(s) != s {
		//	content = strings.ReplaceAll(content, strings.Title(s), "<span class=\"highlight-style\">"+strings.Title(s)+"</span>")
		//}
	}
	return template.HTML(content)
}

func IsSpecialString(str string) bool {
	for _, s := range str {
		if unicode.IsOneOf([]*unicode.RangeTable{
			unicode.Han,
			unicode.Digit,
			unicode.Letter,
			unicode.Number,
		}, s) {
			return false
		}
	}
	return true
}

func Translate(content string, transType string) string {
	type transRes struct {
		Isdict     int     `json:"isdict"`
		Confidence float64 `json:"confidence"`
		Target     string  `json:"target"`
		Rc         int     `json:"rc"`
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("x-authorization", "token 9sdftiq37bnv410eon2l")

	res, err := net.Post("https://api.interpreter.caiyunai.com/v1/translator", headers, map[string]interface{}{
		"source":     content,
		"trans_type": transType, // auto2en
		"request_id": strconv.FormatInt(time.Now().Unix(), 10),
		"detect":     true,
	}, "", time.Duration(config.EngineTimeout)*time.Millisecond)
	if err != nil {
		return content
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("翻译 Response Body 关闭异常！" + err.Error())
			return
		}
	}(res.Body)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("翻译 Response Body 转换异常！" + err.Error())
		return content
	}

	var body transRes
	err = json.Unmarshal(responseBody, &body)
	if err != nil {
		log.Println("翻译 Response Body 转换json异常！" + err.Error())
		return content
	}

	return body.Target
}
