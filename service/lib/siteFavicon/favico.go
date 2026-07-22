package siteFavicon

import (
	"context"
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
	"sun-panel/global"
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

// getOneFaviconTimeout 是 GetOneFaviconURL 的整体超时时间。
// 自托管/内网环境下公共 favicon 服务可能不可达，整体超时可避免请求长时间挂起。
const getOneFaviconTimeout = 3 * time.Second

// faviconServiceDesc 描述一个公共 favicon 兜底服务
type faviconServiceDesc struct {
	name  string
	build func(host string) string
}

// defaultFaviconServices 默认兜底服务列表，顺序即优先级。
// 默认首选 DuckDuckGo（国内通常可达），Google 作为备选保留，其余为补充源。
var defaultFaviconServices = []faviconServiceDesc{
	{
		name: "duckduckgo",
		build: func(host string) string {
			return fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", host)
		},
	},
	{
		name: "google",
		build: func(host string) string {
			return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=128", host)
		},
	},
	{
		name: "iconhorse",
		build: func(host string) string {
			return fmt.Sprintf("https://icon.horse/icon/%s", host)
		},
	},
	{
		name: "clearbit",
		build: func(host string) string {
			return fmt.Sprintf("https://logo.clearbit.com/%s", host)
		},
	},
}

// faviconServiceByName 通过名称快速查找兜底服务描述
var faviconServiceByName = func() map[string]faviconServiceDesc {
	m := make(map[string]faviconServiceDesc, len(defaultFaviconServices))
	for _, s := range defaultFaviconServices {
		m[s.name] = s
	}
	return m
}()

// resolveFaviconServices 解析实际使用的兜底服务列表。
// 若配置文件（conf.ini）[favicon] fallback_services 指定了以逗号分隔的服务名，则按该顺序使用；
// 否则使用默认列表（DuckDuckGo 优先）。非法/未知名称会被忽略，全部无效时回退到默认列表。
func resolveFaviconServices() []faviconServiceDesc {
	cfg := ""
	if global.Config != nil {
		cfg = strings.TrimSpace(global.Config.GetValueString("favicon", "fallback_services"))
	}
	if cfg == "" {
		return defaultFaviconServices
	}

	parts := strings.Split(cfg, ",")
	resolved := make([]faviconServiceDesc, 0, len(parts))
	for _, p := range parts {
		name := strings.ToLower(strings.TrimSpace(p))
		if name == "" {
			continue
		}
		if svc, ok := faviconServiceByName[name]; ok {
			resolved = append(resolved, svc)
		}
	}
	if len(resolved) == 0 {
		return defaultFaviconServices
	}
	return resolved
}

func IsHTTPURL(url string) bool {
	httpPattern := `^(http://|https://|//)`
	match, err := regexp.MatchString(httpPattern, url)
	if err != nil {
		return false
	}
	return match
}

// GetOneFaviconURL 从目标站点获取一个可用的 favicon URL（使用默认 3s 整体超时）。
// 策略：先解析页面 link 标签；失败或无结果时尝试站点默认路径；最后使用多个公共 favicon 服务兜底。
// 约定：当所有源都失败时返回空字符串且 err 为 nil，交由上层（前端）走占位逻辑，而不返回错误。
func GetOneFaviconURL(urlStr string) (string, error) {
	return GetOneFaviconURLWithTimeout(urlStr, getOneFaviconTimeout)
}

// GetOneFaviconURLWithTimeout 与 GetOneFaviconURL 相同，但允许自定义整体超时。
func GetOneFaviconURLWithTimeout(urlStr string, timeout time.Duration) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Host == "" {
		return "", errors.New("invalid URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	iconURL, _ := getOneFaviconURLInner(ctx, parsedURL)
	if iconURL == "" {
		// 所有源均失败（含整体超时）：返回空字符串，不返回错误，交由前端走占位逻辑
		return "", nil
	}
	return iconURL, nil
}

func getOneFaviconURLInner(ctx context.Context, parsedURL *url.URL) (string, error) {
	urlStr := parsedURL.String()

	// 1. 解析页面中的 link[rel*=icon]
	iconURLs, err := getFaviconURLContext(ctx, urlStr)
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
		if _, herr := GetRemoteFileSizeContext(ctx, candidate); herr == nil {
			return candidate, nil
		}
	}

	// 3. 使用可配置的兜底服务列表（默认 DuckDuckGo 优先）
	services := resolveFaviconServices()
	host := parsedURL.Host
	for _, svc := range services {
		candidate := svc.build(host)
		if _, ferr := GetRemoteFileSizeContext(ctx, candidate); ferr == nil {
			return candidate, nil
		}
	}

	// 全部失败：返回空字符串（上层据此走占位逻辑）
	return "", nil
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

// setProbeHeaders 为可用性探测请求设置通用请求头，尽量模拟浏览器以避免被拒绝。
func setProbeHeaders(req *http.Request) {
	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.google.com/")
}

// probeHeadContext 发起 HEAD 请求探测远程文件可访问性。
// 返回是否可达（2xx）以及响应中的内容长度。HEAD 通常不带响应体，速度最快，优先使用。
func probeHeadContext(ctx context.Context, url string) (ok bool, size int64) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, 0
	}
	setProbeHeaders(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()
	// 丢弃可能携带的少量响应体，保持连接整洁
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4*1024))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, resp.ContentLength
	}
	// 非 2xx（如 405 Method Not Allowed / 403 / 501 等）视作需要降级用 GET 复核
	return false, 0
}

// probeGetContext 在 HEAD 不受支持或返回非 2xx 时，降级用 GET（Range: bytes=0-0）校验可访问性。
// 通过 Range 头只请求首个字节，并用 io.LimitReader 限制最多读取 4KB，避免无限下载整张图片。
// 返回是否可达以及文件大小（优先使用 Content-Range 中的总大小，其次回退到 Content-Length）。
func probeGetContext(ctx context.Context, url string) (ok bool, size int64) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, 0
	}
	setProbeHeaders(req)
	// 仅请求第一个字节，大多数服务器支持 Range 并返回 206，不支持时回退到 200 全量（仍被下方限制读取）
	req.Header.Set("Range", "bytes=0-0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()
	// 仅读取少量字节，避免下载整个文件
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4*1024))

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
		total := parseContentRangeTotal(resp.Header.Get("Content-Range"))
		if total <= 0 {
			total = resp.ContentLength
		}
		if total <= 0 {
			// 无法获知确切大小，但资源确实可用
			total = 1
		}
		return true, total
	}
	return false, 0
}

// parseContentRangeTotal 从 Content-Range 头（如 "bytes 0-0/1234"）中解析文件总大小。
func parseContentRangeTotal(header string) int64 {
	if header == "" {
		return 0
	}
	// 格式：bytes <start>-<end>/<total>
	slash := strings.LastIndex(header, "/")
	if slash < 0 {
		return 0
	}
	total, err := strconv.ParseInt(strings.TrimSpace(header[slash+1:]), 10, 64)
	if err != nil {
		return 0
	}
	return total
}

// GetRemoteFileSizeContext 在给定 context 下校验远程图标可访问性并返回文件大小。
// 先尝试 HEAD；若 HEAD 不受支持（如 405）、返回非 2xx 或出错，则降级用 GET（Range 限制读取）复核，
// 以避免因部分网站拒绝 HEAD 而误判图标不可用（批量获取失败的主因之一）。
func GetRemoteFileSizeContext(ctx context.Context, url string) (int64, error) {
	if ok, size := probeHeadContext(ctx, url); ok {
		return size, nil
	}
	// HEAD 不可用，降级到 GET 复核
	if ok, size := probeGetContext(ctx, url); ok {
		return size, nil
	}
	return 0, fmt.Errorf("remote file not reachable: %s", url)
}

// GetRemoteFileSize 获取远程文件大小并校验可访问性（无 context 版本，向后兼容）。
func GetRemoteFileSize(url string) (int64, error) {
	return GetRemoteFileSizeContext(context.Background(), url)
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

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
	fileExt := extractImageExt(url)
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

// extractImageExt 从 URL 中提取图片扩展名，去掉查询串/锚点；无法识别时默认 .ico。
// 这样保存在本地的 favicon 始终带有正确扩展名（.ico/.svg/.png 等），便于前端正确识别图标类型。
func extractImageExt(rawURL string) string {
	clean := rawURL
	if idx := strings.IndexAny(clean, "?#"); idx >= 0 {
		clean = clean[:idx]
	}
	ext := strings.ToLower(path.Ext(clean))
	if ext == "" || !strings.Contains(".png.jpg.jpeg.gif.webp.ico.svg.bmp", ext) {
		return ".ico"
	}
	return ext
}

func GetOneFaviconURLAndUpload(urlStr string) (string, bool) {
	iconUrl, err := GetOneFaviconURL(urlStr)
	if err != nil || iconUrl == "" {
		return "", false
	}
	return iconUrl, true
}

func getFaviconURLContext(ctx context.Context, url string) ([]string, error) {
	var icons []string
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return icons, err
	}

	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := httpClient.Do(req)
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
