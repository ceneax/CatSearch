package net

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

func Get(reqUrl string, headers http.Header, proxy string, timeout time.Duration) (*http.Response, error) {
	var client http.Client

	if len(proxy) > 0 {
		uri, _ := url.Parse(proxy)
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(uri),
			},
			Timeout: timeout,
		}
	} else {
		client = http.Client{
			Timeout: timeout,
		}
	}

	req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)
	for key := range headers {
		req.Header.Set(key, headers.Get(key))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Post(reqUrl string, headers http.Header, body map[string]interface{}, proxy string, timeout time.Duration) (*http.Response, error) {
	var client http.Client

	if len(proxy) > 0 {
		uri, _ := url.Parse(proxy)
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(uri),
			},
			Timeout: timeout,
		}
	} else {
		client = http.Client{
			Timeout: timeout,
		}
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(b))
	for key := range headers {
		req.Header.Set(key, headers.Get(key))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
