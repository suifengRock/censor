package service

import (
	"strings"
	"sync"
)

type urlFilter struct {
	urls map[string]uint8
	mu   sync.Mutex
}

func NewUrlFilter() *urlFilter {
	u := new(urlFilter)
	u.urls = make(map[string]uint8)
	return u
}

func (u *urlFilter) IsExist(url string) bool {
	u.mu.Lock()
	_, ok := u.urls[url]
	u.mu.Unlock()
	return ok
}

func (u *urlFilter) Store(url string) bool {
	u.mu.Lock()
	_, ok := u.urls[url]
	if ok {
		u.mu.Unlock()
		return false
	}
	u.urls[url] = 1
	u.mu.Unlock()
	return true
}

func (u *urlFilter) Reset() {
	u.mu.Lock()
	u.urls = make(map[string]uint8)
	u.mu.Unlock()
}

var requestUrlsFilter = NewUrlFilter()
var requestUrlFilterStack = make(chan string, 5120)
var allowDomain = []string{}
var urlFilterHandleRunner = 16

func PushRequestUrlFilterStack(url string) {
	requestUrlFilterStack <- url
}

func checkDomain(url string) bool {
	for _, domain := range allowDomain {
		index := strings.Index(url, domain)
		if index < 0 {
			continue
		}
		if index > 10 {
			return true
		}
		return false
	}
	return true
}

func urlFilterHandle() {
	for {
		select {
		case url := <-requestUrlFilterStack:
			// domaim filter
			if checkDomain(url) {
				continue
			}
			// uniq url
			if requestUrlsFilter.Store(url) {
				// beego.Info(url)
				PushRequestUrl(url)
			}
		}
	}
}

func init() {
	// urlFilterHandleRunner = beego.AppConfig.Int("url_filter_handle_runner", 16)
	for i := 0; i < urlFilterHandleRunner; i++ {
		go urlFilterHandle()
	}

}
