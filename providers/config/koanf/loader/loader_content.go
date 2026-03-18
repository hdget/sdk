package loader

import (
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

type contentConfigLoader struct {
	content []byte
	reader  *koanf.Koanf
}

func NewContentConfigLoader(reader *koanf.Koanf, content []byte) Loader {
	return &contentConfigLoader{
		content: content,
		reader:  reader,
	}
}

func (l *contentConfigLoader) Load() error {
	return l.reader.Load(rawbytes.Provider(l.content), toml.Parser())
}
