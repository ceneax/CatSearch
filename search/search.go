package search

import (
	"CatSearch/config"
	"CatSearch/engine"
	"CatSearch/model"
	"CatSearch/net"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

func Search(context *gin.Context, p model.Params) {
	results := map[string]*model.Result{}

	wg := &sync.WaitGroup{}
	responseChannel := make(chan map[string]model.Result, len(config.EngineList))
	wgResponse := &sync.WaitGroup{}

	go func() {
		wgResponse.Add(1)
		for response := range responseChannel {
			for k, v := range response {
				results[k] = &v
			}
		}
		wgResponse.Done()
	}()

	for engineName := range config.EngineList {
		wg.Add(1)
		go doSearch(engineName, p, responseChannel, wg)
	}

	wg.Wait()
	close(responseChannel)
	wgResponse.Wait()

	DoResponse(context, p, processResult(results, p.Query))
}

func doSearch(engineName string, p model.Params, result chan map[string]model.Result, wg *sync.WaitGroup) {
	defer wg.Done()

	searchEngine, err := engine.GetEngine(engineName)
	if err != nil {
		log.Println("Engine初始化失败！" + err.Error())
		return
	}

	proxy := ""
	if config.Debug {
		if engineName == "google" {
			proxy = "http://127.0.0.1:10809"
		}
	}
	req := searchEngine.Request(p)
	res, err := net.Get(req.Url, req.Headers, proxy, time.Duration(config.EngineTimeout)*time.Millisecond)
	if err != nil {
		log.Println("内部错误！" + err.Error())
		return
	}

	result <- map[string]model.Result{
		engineName: searchEngine.Response(res).Result,
	}
}

func processResult(results map[string]*model.Result, q string) model.Result {
	// 最终返回的结果
	result := model.Result{}
	// 按照config配置项里搜索引擎排序后的结果数组
	var sortedResult []model.Result
	// 搜索结果里结果最多的数量
	maxResultsNum := 0

	// 按照config配置项里搜索引擎的排序来进行结果的排序
	for k := range config.EngineList {
		if results[k] != nil {
			sortedResult = append(sortedResult, *results[k])
			if len(results[k].Results) > maxResultsNum {
				maxResultsNum = len(results[k].Results)
			}

			// 顺便在这个循环里处理Card和Answer
			if k == "google" {
				result.Answer = results[k].Answer
				result.Card = results[k].Card
			}
			// 顺便在这个循环里处理CodeAnswer
			if k == "stackoverflow" {
				result.CodeAnswer = results[k].CodeAnswer
			}
		}
	}

	// 根据搜索引擎的排序来对所有搜索结果交替合并
	for i := 0; i < maxResultsNum; i++ {
		for j, item := range sortedResult {
			if len(item.Results) > 0 {
				// 过滤无效结果
				if item.Results[0].Url.Raw != "" {
					// 顺便在这个循环里处理内容着色
					//item.Results[0].Content = util.ContentHighlight(string(item.Results[0].Content), q)
					item.Results[0].Title = template.HTML(strings.ReplaceAll(string(item.Results[0].Title), "<", "&lt;"))
					item.Results[0].Title = template.HTML(strings.ReplaceAll(string(item.Results[0].Title), ">", "&gt;"))
					item.Results[0].Content = template.HTML(strings.ReplaceAll(string(item.Results[0].Content), "<", "&lt;"))
					item.Results[0].Content = template.HTML(strings.ReplaceAll(string(item.Results[0].Content), ">", "&gt;"))

					// 顺便在这个循环里标记被Block的网站
					uri, err := url.Parse(item.Results[0].Url.Raw)
					if err == nil && uri.Host != "" && strings.Contains(config.GFW_LIST, uri.Host) {
						item.Results[0].Blocked = true
					}

					// 加入到最终结果变量中
					result.Results = append(result.Results, item.Results[0])
				}
				sortedResult[j].Results = append(sortedResult[j].Results[:0], sortedResult[j].Results[1:]...)
			}
		}
	}

	// 结果去重
	result.Results = removeDuplicationMap(result.Results)

	return result
}

func removeDuplicationMap(arr []model.Results) []model.Results {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		if v.Url.Raw[len(v.Url.Raw)-1:] == "/" {
			v.Url.Raw = v.Url.Raw[:len(v.Url.Raw)-1]
		}
		_, ok := set[v.Url.Raw]
		if ok {
			log.Println(v.Url.Raw)
			continue
		}
		set[v.Url.Raw] = struct{}{}
		arr[j] = v
		j++
	}
	return arr[:j]
}
