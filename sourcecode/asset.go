package sourcecode

import (
	"embed"
	"path"
)

type assetManager struct {
	absPath string // 调用方嵌入文件系统的绝对路径
	relPath string // 调用方嵌入文件系统的相对路径
	fs      embed.FS
}

func newAssetManager(fs embed.FS, embedAbsPath, embedRelPath string) *assetManager {
	return &assetManager{absPath: embedAbsPath, relPath: embedRelPath, fs: fs}
}

func (m *assetManager) Load(file string) ([]byte, error) {
	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
	return m.fs.ReadFile(path.Join(m.relPath, file))
}

func (m *assetManager) Store(file string, data any) error {
	storePath := path.Join(m.absPath, file)
	return saveFile(storePath, data)
}
