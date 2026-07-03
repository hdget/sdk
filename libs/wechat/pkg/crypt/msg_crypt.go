package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// WXBizMsgCrypt 微信消息加解密主类
type WXBizMsgCrypt struct {
	appIdOrCorpId string
	token         string
	key           []byte
}

type EncryptResult struct {
	Encrypt      string `xml:"Encrypt" json:"Encrypt"`
	MsgSignature string `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    string `xml:"Timestamp" json:"Timestamp"`
	Nonce        string `xml:"Nonce" json:"Nonce"`
}

// NewBizMsgCrypt 创建WXBizMsgCrypt实例
func NewBizMsgCrypt(appIdOrCorpId, token, encodingAESKey string) (*WXBizMsgCrypt, error) {
	// Base64解码AESKey
	key, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil || len(key) != 32 {
		return nil, errors.New("invalid EncodingAESKey")
	}

	return &WXBizMsgCrypt{
		appIdOrCorpId: appIdOrCorpId,
		token:         token,
		key:           key,
	}, nil
}

// Encrypt 加密消息
func (w *WXBizMsgCrypt) Encrypt(replyMsg, nonce string, timestamp ...string) (*EncryptResult, error) {
	var ts string
	if len(timestamp) > 0 {
		ts = timestamp[0]
	} else {
		ts = fmt.Sprintf("%d", time.Now().Unix())
	}

	pc, err := newCBCEncryptor(w.key)
	if err != nil {
		return nil, err
	}

	encrypt, err := pc.Encrypt(replyMsg, w.appIdOrCorpId)
	if err != nil {
		return nil, err
	}

	signature, err := CalculateSignature(w.token, ts, nonce, encrypt)
	if err != nil {
		return nil, err
	}

	return &EncryptResult{Encrypt: encrypt, MsgSignature: signature, Timestamp: ts, Nonce: nonce}, nil
}

// Decrypt 解密消息
func (w *WXBizMsgCrypt) Decrypt(msgSignature, timestamp, nonce, encryptedData string) ([]byte, error) {
	calculatedSignature, err := CalculateSignature(w.token, timestamp, nonce, encryptedData)
	if err != nil {
		return nil, err
	}

	if calculatedSignature != msgSignature {
		return nil, fmt.Errorf("signature not matched, recv: %s, calculated: %s", msgSignature, calculatedSignature)
	}

	pc, err := newCBCEncryptor(w.key)
	if err != nil {
		return nil, err
	}

	return pc.Decrypt(encryptedData, w.appIdOrCorpId)
}

// CalculateSignature 计算SHA1签名
func CalculateSignature(token, timestamp, nonce, data string) (string, error) {
	inputs := []string{token, timestamp, nonce, data}
	sort.Strings(inputs)
	var buffer bytes.Buffer
	for _, input := range inputs {
		buffer.WriteString(input)
	}
	h := sha1.New()
	_, err := h.Write(buffer.Bytes())
	if err != nil {
		return "", errors.New("compute Signature")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// PKCS7Encoder PKCS7填充
type PKCS7Encoder struct {
	BlockSize int
}

// Encode PKCS7填充
func (p *PKCS7Encoder) Encode(text []byte) []byte {
	textLength := len(text)
	amountToPad := p.BlockSize - (textLength % p.BlockSize)
	if amountToPad == 0 {
		amountToPad = p.BlockSize
	}

	pad := byte(amountToPad)
	padded := make([]byte, textLength+amountToPad)
	copy(padded, text)
	for i := textLength; i < len(padded); i++ {
		padded[i] = pad
	}

	return padded
}

// Decode PKCS7去除填充
func (p *PKCS7Encoder) Decode(decrypted []byte) []byte {
	pad := int(decrypted[len(decrypted)-1])
	if pad < 1 || pad > 32 {
		pad = 0
	}
	return decrypted[:len(decrypted)-pad]
}

// CBCEncryptor 加解密核心类
type CBCEncryptor struct {
	key  []byte
	mode cipher.BlockMode
}

// newCBCEncryptor 创建Prpcrypt实例
func newCBCEncryptor(key []byte) (*CBCEncryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用CBC模式
	mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])

	return &CBCEncryptor{
		key:  key,
		mode: mode,
	}, nil
}

// Encrypt 加密
func (p *CBCEncryptor) Encrypt(text, appid string) (string, error) {
	// 16位随机字符串
	randomStr, err := p.getRandomStr()
	if err != nil {
		return "", errors.New("replyAsXML random str")
	}

	// 添加长度信息和appid
	textLength := make([]byte, 4)
	binary.BigEndian.PutUint32(textLength, uint32(len(text)))
	textWithInfo := randomStr + string(textLength) + text + appid

	// PKCS7填充
	pkcs7 := PKCS7Encoder{BlockSize: 32}
	paddedText := pkcs7.Encode([]byte(textWithInfo))

	// 加密
	ciphertext := make([]byte, len(paddedText))
	p.mode.CryptBlocks(ciphertext, paddedText)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密
func (p *CBCEncryptor) Decrypt(text, appid string) ([]byte, error) {
	// Base64解码
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, errors.Wrap(err, "base64 decode text")
	}

	// 创建解密器
	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, errors.Wrap(err, "new cipher")
	}

	mode := cipher.NewCBCDecrypter(block, p.key[:aes.BlockSize])

	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除PKCS7填充
	pkcs7 := PKCS7Encoder{BlockSize: 32}
	content := pkcs7.Decode(plaintext)

	// 解析内容
	if len(content) < 20 {
		return nil, errors.New("illegal buffer")
	}

	// 16位随机字符串 + 4位长度 + 数据 + appIdOrCorpId
	xmlLen := binary.BigEndian.Uint32(content[16:20])
	if len(content) < int(20+xmlLen) {
		return nil, errors.New("illegal buffer")
	}

	xmlContent := content[20 : 20+xmlLen]
	fromAppid := string(content[20+xmlLen:])
	if fromAppid != appid {
		return nil, errors.New("invalid app id")
	}

	return xmlContent, nil
}

// getRandomStr 生成16位随机字符串
func (p *CBCEncryptor) getRandomStr() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}
