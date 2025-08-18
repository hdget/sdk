package encoding

import (
	"crypto/rand"
	"github.com/sqids/sqids-go"
	"io"
)

type encodingConfig struct {
	sqidsOption *sqids.Options
	salt        string
}

type Option func(o *encodingConfig)

const (
	defaultMinLength = 6
	defaultAlphabet  = "0123456789ABCDEFGHJKMNPQRSTVWXYZ" // 去除I（避免与数字1混淆）,L（避免与数字1混淆）,O（避免与数字0混淆）, U（避免与V混淆）
)

func WithAlphabet(alphabet string) Option {
	return func(o *encodingConfig) {
		o.sqidsOption.Alphabet = alphabet
	}
}

func WithMinLength(minLength uint8) Option {
	return func(o *encodingConfig) {
		o.sqidsOption.MinLength = minLength
	}
}

func WithSalt(salt string) Option {
	return func(o *encodingConfig) {
		o.salt = salt
	}
}

func WithRandomSalt(saltLength int) Option {
	return func(o *encodingConfig) {
		salt := make([]byte, saltLength) // 生成8字节的盐
		if _, err := io.ReadFull(rand.Reader, salt); err == nil {
			o.salt = string(salt)
		}
	}
}
