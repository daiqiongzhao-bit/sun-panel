package siteFavicon

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sun-panel/lib/cmn"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var httpClient = &http.Client{
	Timeout: 8 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("too many redirects")
		}
		return nil
	},
}

const chromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

func IsHTTPURL(url string) bool {
	httpPattern := `^(http://|https://|//)`
	match, err := regexp.MatchString(httpPattern, url)
	if err != nil {
		return false
	}
	return match
}

// GetOneFaviconURL 从目标站点获取一个可用的 favicon URL。
// 策略：先解析页面 link 标签；失败或无结果时尝试站点默认路径；最后使用多个公共 favicon 服务兜底。
func GetOneFaviconURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Host == "" {
		return "", errors.New("invalid URL")
	}

	iconURLs, err := getFaviconURL(urlStr)
	if err == nil && len(iconURLs) > 0 {
		return resolveIconURL(iconURLs[0], parsedURL), nil
	}

	// 2. 尝试常见默认路径
	defaultPaths := []string{
		"/favicon.ico",
		"/favicon.png",
		"/apple-touch-icon.png",
		"/apple-touch-icon-precomposed.png",
	}
	base := parsedURL.Scheme + "://" + parsedURL.Host
	for _, p := range defaultPaths {
		candidate := base + p
		if _, herr := GetRemoteFileSize(candidate); herr == nil {
			return candidate, nil
		}
	}

	// 3. 使用公共 favicon 服务兜底
	host := parsedURL.Host
	fallbacks := []string{
		fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", host),
		fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=128", host),
		fmt.Sprintf("https://icon.horse/icon/%s", host),
		fmt.Sprintf("https://logo.clearbit.com/%s", host),
	}
	for _, candidate := range fallbacks {
		if _, ferr := GetRemoteFileSize(candidate); ferr == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("not found ico: %v", err)
}

func resolveIconURL(iconURL string, baseURL *url.URL) string {
	if IsHTTPURL(iconURL) {
		return iconURL
	}
	if strings.HasPrefix(iconURL, "//") {
		return baseURL.Scheme + ":" + iconURL
	}
	u, err := url.Parse(iconURL)
	if err != nil {
		return baseURL.Scheme + "://" + baseURL.Host + "/" + strings.TrimPrefix(iconURL, "/")
	}
	return baseURL.ResolveReference(u).String()
}

// GetRemoteFileSize 获取远程文件大小并校验可访问性
func GetRemoteFileSize(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.google.com/")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed, status code: %d", resp.StatusCode)
	}

	return resp.ContentLength, nil
}

// DownloadImage 下载图片到本地
func DownloadImage(url, savePath string, maxSize int64) (*os.File, error) {
	fileSize, err := GetRemoteFileSize(url)
	if err != nil {
		return nil, err
	}

	if fileSize > maxSize {
		return nil, fmt.Errorf("文件太大，不下载。大小：%d字节", fileSize)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Referer", "https://www.google.com/")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed, status code: %d", response.StatusCode)
	}

	urlFileName := path.Base(url)
	fileExt := path.Ext(url)
	if fileExt == "" {
		fileExt = ".ico"
	}
	fileName := cmn.Md5(fmt.Sprintf("%s%s", urlFileName, time.Now().String())) + fileExt

	destination := savePath + "/" + fileName

	file, err := os.Create(destination)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetOneFaviconURLAndUpload(urlStr string) (string, bool) {
	iconUrl, err := GetOneFaviconURL(urlStr)
	if err != nil || iconUrl == "" {
		return "", false
	}
	return iconUrl, true
}

func getFaviconURL(url string) ([]string, error) {
	var icons []string
	client := httpClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return icons, err
	}

	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return icons, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return icons, errors.New("HTTP request failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return icons, err
	}

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		href, _ := s.Attr("href")

		if strings.Contains(rel, "icon") && href != "" {
			icons = append(icons, href)
		}
	})

	if len(icons) == 0 {
		return icons, errors.New("favicon not found on the page")
	}

	return icons, nil
}
