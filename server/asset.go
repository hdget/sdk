package server

import (
	"embed"
	"encoding/json"
	"github.com/hdget/utils/convert"
	"os"
	"path"
	"path/filepath"
)

type AssetManager interface {
	Load(file string) ([]byte, error)
	Store(file string, data any) error
}

type assetManagerImpl struct {
	embedAbsPath string
	embedRelPath string
	fs           embed.FS
}

func newAssetManager(embedAbsPath, embedRelPath string, fs embed.FS) AssetManager {
	return &assetManagerImpl{embedAbsPath: embedAbsPath, embedRelPath: embedRelPath, fs: fs}
}

func (m *assetManagerImpl) Load(file string) ([]byte, error) {
	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
	return m.fs.ReadFile(path.Join(m.embedRelPath, file))
}

func (m *assetManagerImpl) Store(file string, data any) error {
	var err error
	var content []byte
	switch v := data.(type) {
	case string:
		content = convert.StringToBytes(v)
	case []byte:
		content = v
	default:
		content, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	fullPath := filepath.Join(m.embedAbsPath, file)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
