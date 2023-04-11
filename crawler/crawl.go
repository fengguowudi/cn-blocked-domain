package crawler

import (
    "errors"
    "net/http"
    "net/url"
    "runtime"
    "strings"
)

const getMethod = http.MethodGet

var userAgent = genUA()
var allowedLanguages = []string{"zh-CN", "zh", "en"}

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

// Crawl crawls webpage content and returns *gzip.Reader
func Crawl(target, referer string) (*http.Response, error) {
    if _, err := url.Parse(target); err != nil {
        return nil, err
    }

    acceptLanguage := strings.Join(allowedLanguages, ",")
    req, err := http.NewRequest(getMethod, target, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("User-Agent", userAgent)
    req.Header.Set("Referer", referer)
    req.Header.Set("Accept-Language", acceptLanguage)

    var client = &http.Client{
        Transport: &http.Transport{
            DisableCompression: true,
        },
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("bad status code: " + http.StatusText(resp.StatusCode))
    }

    return resp, nil
}
