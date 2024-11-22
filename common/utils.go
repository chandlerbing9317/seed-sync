package common

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"seed-sync/config"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func GetProxyFunc(proxy bool) func(req *http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if proxy {
			proxy := config.Conf.ProxyConfig
			if proxy.ProxyUsername != "" && proxy.ProxyPassword != "" {
				proxyURL, err := url.Parse(proxy.ProxyURL)
				if err != nil {
					return nil, err
				}
				proxyURL.User = url.UserPassword(proxy.ProxyUsername, proxy.ProxyPassword)
				return proxyURL, nil
			}
			return url.Parse(proxy.ProxyURL)
		}
		return nil, nil
	}
}

func GetRequest(method string, url string, header map[string]string, body any) (*http.Request, error) {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	// 设置header
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return req, nil
}

func GetNextExecuteTime(cronExpr string) (time.Time, error) {
	// 创建一个解析器，支持秒级别的cron表达式
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, err
	}

	// 计算下一次执行时间
	next := schedule.Next(time.Now())
	return next, nil
}

type lockedRandomSource struct {
	mut sync.Mutex
	src rand.Source
}

func NewLockedRandomSource(seed int64) rand.Source {
	return &lockedRandomSource{
		src: rand.NewSource(seed),
	}
}

func (l *lockedRandomSource) Int63() int64 {
	l.mut.Lock()
	defer l.mut.Unlock()
	return l.src.Int63()
}
func (l *lockedRandomSource) Seed(seed int64) {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.src.Seed(seed)
}
