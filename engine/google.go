package engine

import (
	"CatSearch/config"
	"CatSearch/model"
	"CatSearch/util"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type google struct {
	name string
}

func (e *google) GetName() string {
	return e.name
}

func (e *google) Request(p model.Params) Request {
	headers := http.Header{}
	headers.Set("User-Agent", config.HEADER_USER_AGENT)
	headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	page := 1
	intPage, err := strconv.Atoi(p.Page)
	if err == nil {
		page = intPage
	}

	return Request{
		Url:     "https://www.google.com/search?q=" + url.QueryEscape(p.Query) + "&ie=utf8&oe=utf8&start=" + strconv.Itoa((page-1)*10),
		Headers: headers,
	}
}

func (e *google) Response(res *http.Response) Response {
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
		return Response{}
	}

	result := model.Result{}

	answer := doc.Find(".g.wF4fFd.JnwWd.g-blk")
	if answer.Length() > 0 {
		gClass := answer.Find(".g")
		if gClass.Length() > 0 {
			a := answer.Find(".yuRUbf a").First()
			result.Answer.Url = a.AttrOr("href", "")
			result.Answer.Title = a.Find("h3").Text()
			gClass.Remove()
		}

		c := answer.Find(".V3FYCf")
		if c.Find(".LGOjhe").Length() > 0 {
			result.Answer.Ans = append(result.Answer.Ans, c.Text())
		} else if c.Find(".di3YZe").Length() > 0 {
			title := c.Find(".co8aDb").Text()
			result.Answer.Ans = append(result.Answer.Ans, title)
			c.Find(".RqBzHd .TrT0Xe").Each(func(i int, selection *goquery.Selection) {
				result.Answer.Ans = append(result.Answer.Ans, selection.Text())
			})
		}

		answer.Remove()
	}

	doc.Find("#search .g").Each(func(i int, selection *goquery.Selection) {
		if selection.Find(".g").Length() <= 0 && selection.Find(".yuRUbf").Length() > 0 {
			a := selection.Find(".yuRUbf a").First()
			result.Results = append(result.Results, model.Results{
				Url:     util.ParseUrlHostPath(a.AttrOr("href", "")),
				Title:   template.HTML(a.Find("h3").Text()),
				Content: template.HTML(selection.Find(".VwiC3b").First().Text()),
				Source:  e.name,
			})
		}
	})

	card := doc.Find(".TQc1id")
	if card.Length() > 0 {
		one := card.Find(".g").First()
		//selection.Find(".umyQi img").Each(func(j int, selection2 *goquery.Selection) {
		//	log.Println(selection2.AttrOr("src", ""))
		//})
		head := one.Find(".SPZz6b")
		result.Card.Title = head.Find(".qrShPb").Text()
		result.Card.Category = head.Find(".wwUB2c").Text()
		body := one.Find("#kp-wp-tab-overview").Find(".UDZeY.OTFaAf").First()
		body.Find(".wDYxhc").Each(func(i int, selection *goquery.Selection) {
			if selection.Find(".kno-rdesc").Length() > 0 {
				result.Card.Content = selection.Find("span").Text()
			} else if selection.Find(".rVusze").Length() > 0 {
				content := selection.Find(".rVusze")
				if result.Card.Info == nil {
					result.Card.Info = map[string]string{}
				}
				result.Card.Info[content.Find("span.w8qArf>a").Text()] = content.Find("span.LrzXr").Text()
			}
		})
	}

	return Response{
		Result: result,
	}
}
