package module

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
)

type ModuleKind int

const (
	ModuleKindUnknown    ModuleKind = iota
	ModuleKindInvocation            // service invocation module
	ModuleKindEvent                 // topic event module
	ModuleKindDelayEvent            // delay event module
	ModuleKindHealth                // health module
)

type Info struct {
	ApiVersion int    // API版本
	Name       string // 规范化的模块名
	Dir        string // 模块所在目录
}

type Module interface {
	GetApp() string
	GetInfo() *Info
	GetKind() ModuleKind
}

type baseModule struct {
	app        string // 应用名称
	moduleInfo *Info  // 模块的信息
}

var (
	errInvalidModule = errors.New("invalid module, it must be struct")
	moduleNameSuffix = "Module"
	regexFirstV      = regexp.MustCompile(`(?:^|\/)v(\d+)(?:\/([^\/]+.*))?`)
)

var (
	_modules = make(map[ModuleKind][]Module)
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

func (m *baseModule) GetKind() ModuleKind {
	return ModuleKindUnknown
}

func (m *baseModule) GetApp() string {
	return m.app
}

// GetInfo 获取模块元数据信息
func (m *baseModule) GetInfo() *Info {
	return m.moduleInfo
}

// ParseModuleInfo 合法的包路径可能为以下格式：
// * /path/to/v1
// * /path/to/v1/pc
// * /path/to/v2/wxmp
func ParseModuleInfo(pkgPath, moduleName string) (*Info, error) {
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

	return &Info{
		ApiVersion: apiVersion,
		Dir:        dir,
		Name:       trimSuffixIgnoreCase(moduleName, moduleNameSuffix),
	}, nil
}

func Get[T Module](kind ModuleKind) []T {
	results := make([]T, 0)
	for _, mod := range _modules[kind] {
		if got, ok := mod.(T); ok {
			results = append(results, got)
		}
	}
	return results
	//switch m := T.(type) {
	//case module.InvocationModule:
	//	_invocationModules = append(_invocationModules, m)
	//case module.EventModule:
	//	_eventModules = append(_eventModules, m)
	//case module.DelayEventModule:
	//	_delayEventModules = append(_delayEventModules, m)
	//case module.HealthModule:
	//	_healthModules = append(_healthModules, m)
	//}
}

func register(module any) {
	switch m := module.(type) {
	case InvocationModule:
		_modules[ModuleKindInvocation] = append(_modules[ModuleKindInvocation], m)
	case EventModule:
		_modules[ModuleKindEvent] = append(_modules[ModuleKindEvent], m)
	case DelayEventModule:
		_modules[ModuleKindDelayEvent] = append(_modules[ModuleKindDelayEvent], m)
	case HealthModule:
		_modules[ModuleKindHealth] = append(_modules[ModuleKindHealth], m)
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
