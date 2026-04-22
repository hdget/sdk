package lib_ws

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
)

type ModuleInfo struct {
	ApiVersion int    // API版本
	Name       string // 规范化的模块名
	Dir        string // 模块所在目录
}

type Module interface {
	GetApp() string
	GetInfo() *ModuleInfo
}

type baseModule struct {
	app        string      // 应用名称
	moduleInfo *ModuleInfo // 模块的信息
}

var (
	errInvalidModule = errors.New("invalid module, it must be struct")
	moduleNameSuffix = "Module"
	regexFirstV      = regexp.MustCompile(`(?:^|\/)v(\d+)(?:\/([^\/]+.*))?`)
)

var (
	_modules = make([]Module, 0)
)

// newModule 从约定的结构名中解析模块名和版本, 结构名需要为v<number>_<module>
func newModule(app string, moduleObject any) (Module, error) {
	structName := reflectUtils.GetStructName(moduleObject)
	if structName == "" {
		return nil, errInvalidModule
	}

	if !reflectUtils.IsAssignableStruct(moduleObject) {
		return nil, fmt.Errorf("module object: %s must be a pointer to struct", structName)
	}

	// 模块结构体所在的包路径
	pkgPath := getPkgPath(moduleObject)

	// 通过包路径来解析模块信息
	moduleInfo, err := ParseModuleInfo(pkgPath, structName)
	if err != nil {
		return nil, err
	}

	return &baseModule{
		app:        app,
		moduleInfo: moduleInfo,
	}, nil
}

func (m *baseModule) GetApp() string {
	return m.app
}

// GetInfo 获取模块元数据信息
func (m *baseModule) GetInfo() *ModuleInfo {
	return m.moduleInfo
}

// ParseModuleInfo 合法的包路径可能为以下格式：
// * /path/to/v1
// * /path/to/v1/pc
// * /path/to/v2/wxmp
func ParseModuleInfo(pkgPath, moduleName string) (*ModuleInfo, error) {
	strVer, subDirs := getSubDirsAfterFirstV(pkgPath)
	if strVer == "" {
		return nil, errors.New("invalid module path, e,g: /path/to/v1")
	}

	apiVersion, err := strconv.Atoi(strVer)
	if err != nil {
		return nil, errors.New("invalid apiVersion")
	}

	var dir string
	switch len(subDirs) {
	case 0:
		dir = "" // 允许dir为空， 默认为内部调用
	case 1:
		dir = subDirs[0]
	default:
		return nil, errors.New("invalid module path, only supports one sub level")
	}

	return &ModuleInfo{
		ApiVersion: apiVersion,
		Dir:        dir,
		Name:       trimSuffixIgnoreCase(moduleName, moduleNameSuffix),
	}, nil
}

func GetModules() []GinModule {
	results := make([]GinModule, 0)
	for _, m := range _modules {
		if got, ok := m.(GinModule); ok {
			results = append(results, got)
		}
	}
	return results
}

func registerModule(module any) {
	switch m := module.(type) {
	case GinModule:
		_modules = append(_modules, m)
	}
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

func trimSuffixIgnoreCase(s, suffix string) string {
	if len(suffix) > len(s) {
		return s
	}
	if strings.EqualFold(s[len(s)-len(suffix):], suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
