package config

import "github.com/go-ego/gse"

var (
	Debug = false

	ServerIp      = "127.0.0.1"
	ServerPort    = "4000"
	ServerApiPort = "4010"

	EngineTimeout = 10000
	EngineList    = map[string]map[string]string{
		"google": {
			"weight": "0.5",
		},
		"baidu": {
			"weight": "0.5",
		},
		//"stackoverflow": {
		//	"weight": "0.5",
		//},
	}

	Seg gse.Segmenter
)

const (
	HEADER_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.56"
	GFW_LIST          = ``
)
