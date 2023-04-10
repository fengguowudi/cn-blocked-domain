package crawler

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	client *http.Client
	filter *bloom.BloomFilter
)

func init() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 5,
	}
	client = &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
	filter = bloom.NewWithEstimates(1000000, 0.01)
}

func genUA() string {
	switch runtime.GOOS {
	case "linux":
		return `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0`
	case "darwin":
		return `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36`
	case "windows":
		return `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0`
	default:
		return ""
	}
}

// Crawl crawls webpage content and returns *http.Response
func Crawl(target, referer string) (*http.Response, error) {
	if _, err := url.Parse(target); err != nil {
		return nil, err
	}

	if filter.TestString(target) {
		return nil, errors.New("duplicate page")
	}

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", genUA())
	req.Header.Set("Referer", referer)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad status code: " + resp.Status)
	}

	filter.AddString(target)

	return resp, nil
}
