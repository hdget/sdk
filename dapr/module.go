package dapr

import (
	"fmt"
	reflectUtils "github.com/hdget/utils/reflect"
	"github.com/pkg/errors"
	"strconv"
)

type DaprModuleInfo struct {
	Version   int    // 版本
	Namespace string // can be used to distinguish different client
	Name      string // 处理后的模块名
}

type Module interface {
	GetApp() string
	GetModuleInfo() *DaprModuleInfo
}

type baseModule struct {
	app        string          // 应用名称
	moduleInfo *DaprModuleInfo // 模块的信息
}

//
//type ModuleInfo struct {
//	Version   int
//	Namespace string
//
//	PkgPath    string // 包所在的路径
//	ModuleName string // 结构体的全名, e,g: xxxModule
//}

var (
	errInvalidModule = errors.New("invalid module, it must be struct")
	moduleNameSuffix = "Module"
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
	moduleInfo, err := ParseDaprModuleInfo(pkgPath, structName)
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

// GetModuleInfo 获取模块元数据信息
func (m *baseModule) GetModuleInfo() *DaprModuleInfo {
	return m.moduleInfo
}

// ParseDaprModuleInfo 合法的包路径为下列
// /path/to/v1
// /path/to/v1/pc
// /path/to/v2/wxmp
func ParseDaprModuleInfo(pkgPath, moduleName string) (*DaprModuleInfo, error) {
	strVer, subDirs := getSubDirsAfterFirstV(pkgPath)
	if strVer == "" {
		return nil, errors.New("invalid module path, e,g: /path/to/v1")
	}

	version, err := strconv.Atoi(strVer)
	if err != nil {
		return nil, errors.New("invalid version")
	}

	var namespace string
	switch len(subDirs) {
	case 0:
		namespace = ""
	case 1:
		namespace = subDirs[0]
	default:
		return nil, errors.New("invalid module path, only supports one sub level")
	}

	return &DaprModuleInfo{
		Version:   version,
		Namespace: namespace,
		Name:      trimSuffixIgnoreCase(moduleName, moduleNameSuffix),
	}, nil
}
