package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"seed-sync/config"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func GetProxyFunc(isProxy bool) func(req *http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if isProxy {
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
func HasSameElement[T comparable](list []T, elements []T) bool {
	for _, element := range elements {
		if slices.Contains(list, element) {
			return true
		}
	}
	return false
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

func FormatUrlTemplate(template string, params map[string]string) string {
	result := template
	for key, value := range params {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
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

// ValidateURL 校验URL是否合法
func ValidateURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL不能为空")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("URL格式不正确: %v", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL必须使用http或https协议")
	}

	if u.Host == "" {
		return fmt.Errorf("URL必须包含主机名")
	}

	return nil
}

// NormalizeURL 规范化URL
func NormalizeURL(urlStr string) (string, error) {
	urlStr = strings.TrimSpace(urlStr)

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("URL格式不正确: %v", err)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	u.Path = strings.TrimRight(u.Path, "/")

	return u.String(), nil
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
