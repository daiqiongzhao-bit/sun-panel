package assets

import (
	"embed"
	"strings"
)

//go:embed *
var embeddedAssets embed.FS

// Asset 返回嵌入的静态资源。
// 调用方传入的 name 形如 "assets/version" 或 "assets/lang/zh-cn.ini"，
// 这里统一去掉 "assets/" 前缀后从 embed.FS 读取。
func Asset(name string) ([]byte, error) {
	name = strings.TrimPrefix(name, "assets/")
	return embeddedAssets.ReadFile(name)
}
