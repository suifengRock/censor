package service

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

var requestHandleRunner = 32
var requestUrlStack = make(chan string, 5120)
var invalidUrls = []string{
	"#",
	"javascript",
}

func IsInvalid(url string) bool {
	for _, val := range invalidUrls {
		if strings.Contains(url, val) {
			return true
		}
	}
	return false
}

func PushRequestUrl(url string) {
	requestUrlStack <- url
}

func requestHandle() {
	for {
		select {
		case url := <-requestUrlStack:
			data, err := requestUrl(url)
			if err != nil {
				beego.Error("request url err: ", err.Error())
				continue
			}
			PushHtmlBody(data, url)
		}
	}
}

func requestUrl(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		beego.Error("http get err: ", url, err.Error())
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return nil, errors.New("no found.")
	}
	if res.StatusCode != 200 {
		beego.Error("status code error: ", res.StatusCode, res.Status, url)
		return nil, errors.New("http return status != 200.")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Error("http readAll err: ", err.Error())
		return nil, err
	}
	return data, nil
}

func init() {
	// urlFilterHandleRunner = beego.AppConfig.Int("url_filter_handle_runner", 16)
	for i := 0; i < requestHandleRunner; i++ {
		go requestHandle()
	}
}
