package param

import "path/filepath"

type File struct {
	FileConfigType string   // 配置内容类型，e,g: toml, json
	RootDirs       []string // 配置文件所在的RootDirs
	File           string   // 指定的配置文件
	SearchDirs     []string // 未指定配置文件情况下，搜索的目录
	SearchFileName string   // 未指定配置文件情况下，搜索的文件名，不需要文件后缀
}

const (
	defaultConfigType = "toml"
)

var (
	// the default config file search pattern
	//
	//	./config/app/<app>/<app>.test.toml
	//	./common/config/app/<app>/<app>.test.toml
	//	../config/app/<app>/<app>.test.toml
	//  ../common/config/app/<app>/<app>.test.toml
	//  ...
	defaultConfigRootDirs = []string{
		".",                                      // current dir
		filepath.Join("config", "app"),           // default config root dir1
		filepath.Join("common", "config", "app"), // default config root dir2
	}
)

func NewFileDefaultParam() *File {
	return &File{
		FileConfigType: defaultConfigType,
		RootDirs:       defaultConfigRootDirs,
		File:           "",
	}
}
