package encoding

import (
	"encoding/base64"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
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
	salt  []byte
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

	if len(impl.salt) > 0 {
		value, err = addSalt(value, impl.salt)
		if err != nil {
			return ""
		}
	}

	return value
}

func (impl encodingImpl) DecodeInt64(value string) int64 {
	if strings.TrimSpace(value) == "" {
		return 0
	}

	if len(impl.salt) > 0 {
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

func (impl encodingImpl) DecodeInt64Slice(code string) []int64 {
	if strings.TrimSpace(code) == "" {
		return nil
	}

	uint64s := impl.sqids.Decode(code)
	return pie.Map(uint64s, func(v uint64) int64 { return int64(v) })
}

// 加密函数
func addSalt(plaintext string, salt []byte) (string, error) {
	// 将盐和明文拼接在一起
	data := append(salt, []byte(plaintext)...)
	// 使用固定密钥进行XOR加密
	for i := range data {
		data[i] ^= salt[i%len(salt)]
	}
	// 返回Base64编码后的结果
	return base64.URLEncoding.EncodeToString(data), nil
}

// 解密函数
func removeSalt(ciphertextBase64 string, salt []byte) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// 使用固定密钥进行XOR解密
	for i := range data {
		data[i] ^= salt[i%len(salt)]
	}

	// 检查并移除盐
	if len(data) < len(salt) {
		return "", errors.New("ciphertext too short")
	}

	return string(data[len(salt):]), nil
}
