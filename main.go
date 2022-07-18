package main

import (
	"CatSearch/config"
	"CatSearch/route"
	"flag"
	"log"
)

var debug = flag.Bool("d", false, "Debug 模式")

func main() {
	flag.Parse()
	config.Debug = *debug

	err := config.Seg.LoadDict("./s_1.txt")
	if err != nil {
		log.Fatal("分词词典加载失败")
		return
	}

	go route.StartSearcher()
	route.StartApi()
}
