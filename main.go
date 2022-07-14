package main

import (
	"CatSearch/config"
	"CatSearch/route"
)

func main() {
	config.Debug = false
	config.Seg.LoadDict("./s_1.txt")

	go route.StartSearcher()
	route.StartApi()
}
