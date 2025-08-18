package encoding

import (
	"encoding/base64"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/sqids/sqids-go"
	"strings"
	"sync"
)

type Coder interface {
	Encode(ids ...int64) string
	DecodeInt64(code string) int64
	DecodeInt64Slice(code string) []int64
}

type encodingImpl struct {
	sqids *sqids.Sqids
	salt  string
}

var (
	_once     sync.Once
	_encoding *encodingImpl
)

func New(options ...Option) Coder {
	_once.Do(func() {
		c := &encodingConfig{
			sqidsOption: &sqids.Options{
				MinLength: defaultMinLength,
				Alphabet:  defaultAlphabet, // 同ulid
			},
		}

		for _, apply := range options {
			apply(c)
		}

		sqidsInstance, _ := sqids.New(*c.sqidsOption)
		_encoding = &encodingImpl{
			sqids: sqidsInstance,
			salt:  c.salt,
		}
	})
	return _encoding
}

func (impl encodingImpl) Encode(ids ...int64) string {
	if pie.Contains(ids, 0) {
		return ""
	}

	uint64s := pie.Map(ids, func(v int64) uint64 { return uint64(v) })

	value, err := impl.sqids.Encode(uint64s)
	if err != nil {
		return ""
	}

	if impl.salt != "" {
		return addSalt(value, impl.salt)
	}

	return value
}

func (impl encodingImpl) DecodeInt64(value string) int64 {
	if strings.TrimSpace(value) == "" {
		return 0
	}

	if impl.salt != "" {
		var err error
		value, err = removeSalt(value, impl.salt)
		if err != nil {
			return 0
		}
	}

	uint64s := impl.sqids.Decode(value)
	if len(uint64s) <= 0 {
		return 0
	}

	return int64(uint64s[0])
}

func (impl encodingImpl) DecodeInt64Slice(value string) []int64 {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	if impl.salt != "" {
		var err error
		value, err = removeSalt(value, impl.salt)
		if err != nil {
			return nil
		}
	}

	uint64s := impl.sqids.Decode(value)
	return pie.Map(uint64s, func(v uint64) int64 { return int64(v) })
}

// 加密函数
func addSalt(plaintext, salt string) string {
	// 将盐和明文拼接在一起
	salted := salt + plaintext
	encoded := base64.URLEncoding.EncodeToString([]byte(salted))
	return strings.TrimRight(encoded, "=")
}

// 解密函数
func removeSalt(encoded, salt string) (string, error) {
	// 补全Base64填充
	if len(encoded)%4 != 0 {
		encoded += strings.Repeat("=", 4-len(encoded)%4)
	}

	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	// 去除盐值
	if len(decoded) < len(salt) {
		return "", fmt.Errorf("invalid encoded string")
	}
	return string(decoded[len(salt):]), nil
}
