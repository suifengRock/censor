package service

import (
	"bufio"
	"bytes"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

type HtmlBody struct {
	Data []byte
	Url  string
}

var htmlBodyStack = make(chan *HtmlBody, 1024)
var htmlBodyHandleRunner = 64

func PushHtmlBody(data []byte, url string) {
	body := new(HtmlBody)
	body.Data = data
	body.Url = url
	htmlBodyStack <- body
}

func bodyDataHandle() {

	for {
		select {
		case body := <-htmlBodyStack:
			// 深度查找url
			DeepSearchUrl(body.Data, body.Url)
			// 检测内容
			CheckData(body.Data, body.Url)
		}
	}

}

func CheckData(data []byte, reqUrl string) error {
	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		_, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			beego.Error("read resp data err: ", err.Error())
			return err
		}
		// fmt.Println(line)
	}
}

func DeepSearchUrl(data []byte, reqUrl string) error {
	u, err := url.Parse(reqUrl)
	if err != nil {
		beego.Error("url parse err: ", err.Error())
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		beego.Error("goquery init doc err: ", err.Error())
		return err
	}

	doc.Find("body a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, _ := s.Attr("href")
		if href != "" {
			if IsInvalid(href) {
				return
			}
			if strings.Index(href, "http") != 0 {
				href = u.Scheme + "://" + u.Host + "/" + href
			}
			PushRequestUrlFilterStack(href)
		}
	})
	return nil
}

func init() {
	// urlFilterHandleRunner = beego.AppConfig.Int("url_filter_handle_runner", 16)
	for i := 0; i < htmlBodyHandleRunner; i++ {
		go bodyDataHandle()
	}
}
