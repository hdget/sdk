package dapr

import (
	"embed"
	"encoding/json"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/utils/convert"
	"github.com/hdget/utils/text"
	"path"
	"reflect"
	"regexp"
	"strings"
)

var (
	truncateSize = 200
	regexFirstV  = regexp.MustCompile(`(?:^|\/)v(\d+)(?:\/([^\/]+.*))?`)
)

const (
	fileExposedHandlers = ".exposed_handlers.json"
)

func truncate(data []byte) string {
	return text.Truncate(convert.BytesToString(data), truncateSize)
}

func trimSuffixIgnoreCase(s, suffix string) string {
	if len(suffix) > len(s) {
		return s
	}
	if strings.EqualFold(s[len(s)-len(suffix):], suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func getPkgPath(v interface{}) string {
	t := reflect.TypeOf(v)

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return ""
	}

	return t.PkgPath()
}

// getSubDirsAfterFirstV 在路径中找到第一个v<数字>出现的位置，获取其版本号，并获取后面的子目录
func getSubDirsAfterFirstV(path string) (version string, dirs []string) {
	match := regexFirstV.FindStringSubmatch(path)
	if len(match) < 2 {
		return "", nil // 无匹配
	}

	version = match[1] // 提取数字部分（如 "1"）
	if len(match) >= 3 && match[2] != "" {
		// 分割剩余路径，过滤空字符串和文件
		parts := strings.Split(match[2], "/")
		for _, part := range parts {
			if part != "" && !strings.Contains(part, ".") { // 忽略文件名
				dirs = append(dirs, part)
			}
		}
	}
	return
}

// LoadStoredExposedHandlers 从embed.FS中加载ast解析后保存的DaprHandlers
func LoadStoredExposedHandlers(fs embed.FS) ([]*protobuf.DaprHandler, error) {
	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
	data, err := fs.ReadFile(path.Join("json", fileExposedHandlers))
	if err != nil {
		return nil, err
	}

	var handlers []*protobuf.DaprHandler
	err = json.Unmarshal(data, &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}
