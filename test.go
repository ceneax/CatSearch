package main

import "github.com/go-ego/gse"

func main1() {
	var Seg gse.Segmenter
	Seg.LoadDict()
	println(Seg.CutStr(Seg.Cut("android怎么跳转activity并拿到返回结果", true), "|"))
}
