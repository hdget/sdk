package encoding

import (
	"crypto/rand"
	"github.com/sqids/sqids-go"
	"io"
)

type encodingConfig struct {
	sqidsOption *sqids.Options
	salt        []byte
}

type Option func(o *encodingConfig)

const (
	defaultMinLength = 6
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
		o.salt = []byte(salt)
	}
}

func WithRandomSalt(saltLength int) Option {
	return func(o *encodingConfig) {
		salt := make([]byte, saltLength) // 生成8字节的盐
		if _, err := io.ReadFull(rand.Reader, salt); err == nil {
			o.salt = salt
		}
	}
}
