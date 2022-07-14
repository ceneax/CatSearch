package engine

import (
	"CatSearch/model"
	"errors"
	"net/http"
)

type Engine interface {
	GetName() string
	Request(p model.Params) Request
	Response(res *http.Response) Response
}

type Request struct {
	Url     string
	Headers http.Header
}

type Response struct {
	Result model.Result
}

func GetEngine(name string) (Engine, error) {
	switch name {
	case "google":
		return &google{
			name: name,
		}, nil
	case "baidu":
		return &baidu{
			name: name,
		}, nil
	case "stackoverflow":
		return &stackoverflow{
			name: name,
		}, nil
	}
	return nil, errors.New("未找到对应的Engine: " + name)
}
