package data

import (
	"embed"
	"encoding/json"
	"github.com/hdget/utils/convert"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type DataManager interface {
	Load(file string) ([]byte, error)
	Store(file string, data any) error
	LoadTempFile(file string) ([]byte, error)
	StoreTempFile(file string, data any) error
}

type dataManagerImpl struct {
	fsDir   string
	fs      embed.FS
	tempDir string
}

const (
	tempDir = "temp"
)

func New(fs embed.FS, options ...DataOption) DataManager {
	_, callerPath, _, _ := runtime.Caller(0)
	m := &dataManagerImpl{fs: fs, tempDir: tempDir, fsDir: filepath.Base(filepath.Dir(callerPath))}
	for _, option := range options {
		option(m)
	}
	return m
}

func (m *dataManagerImpl) Load(file string) ([]byte, error) {
	content, err := m.fs.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (m *dataManagerImpl) LoadTempFile(file string) ([]byte, error) {
	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
	return m.Load(path.Join(tempDir, file))
}

func (m *dataManagerImpl) Store(file string, data any) error {
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

	// go generate运行时实在main那一级的目录
	absPath := filepath.Join(m.fsDir, file)

	dir := filepath.Dir(absPath)
	if dir != "." {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(absPath)
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

func (m *dataManagerImpl) StoreTempFile(file string, data any) error {
	return m.Store(path.Join(m.tempDir, file), data)
}
