package sourcecode

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type metaDataManager struct {
	srcDir string
}

//type EmbedInfo struct {
//	VarName string
//	AbsPath string // 调用方嵌入文件系统的绝对路径
//	RelPath string // 调用方嵌入文件系统的相对路径
//}

type MetaData struct {
	ModulePaths     map[string]string // 模块的路径
	ServerEntryFile string            // appServer.Run的入口文件即appServer开始运行所在的go文件
	// EmbedInfos      []*EmbedInfo      // 获取所有嵌入类型的定义
}

const (
	fileMeta = ".meta" // 源代码信息
)

func newMetaDataManager(srcDir string) *metaDataManager {
	return &metaDataManager{
		srcDir: srcDir,
	}
}

func (m *metaDataManager) Store(data any) error {
	absPath := filepath.Join(m.srcDir, fileMeta)
	return saveFile(absPath, data)
}

func (m *metaDataManager) Load() (*MetaData, error) {
	// 读取文件内容
	// 使用 os.ReadFile 读取文件内容
	content, err := os.ReadFile(filepath.Join(m.srcDir, fileMeta))
	if err != nil {
		return nil, errors.Wrapf(err, "failed read meta file, file: %s", fileMeta)
	}

	var meta MetaData
	err = json.Unmarshal(content, &meta)
	if err != nil {
		return nil, errors.Wrap(err, "invalid meta file")
	}

	return &meta, nil
}

func (m *metaDataManager) Print() error {
	meta, err := m.Load()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "PATH"})
	table.SetRowLine(true)
	table.Append([]string{
		"SERVER_ENTRY", meta.ServerEntryFile,
	})
	for k, v := range meta.ModulePaths {
		table.Append([]string{k, v})
	}
	table.Render() // Send output
	return nil
}
