package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
)

func genUA() (userAgent string) {
	switch runtime.GOOS {
	case "linux":
		userAgent = `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0`
	case "darwin":
		userAgent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36`
	case "windows":
		userAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0`
	}
	return
}

func Crawl(target, referer string) (*http.Response, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"User-Agent":      {genUA()},
			"Referer":         {referer},
			"Accept-Encoding": {"gzip"},
			"Accept-Language": {"zh-CN,zh;q=0.9,en;q=0.8"},
			"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9"},
		},
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return resp, nil
}
