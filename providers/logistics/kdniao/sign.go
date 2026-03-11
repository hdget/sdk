package logistics_kdniao

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// Sign 生成快递鸟签名
// 签名算法: DataSign = URLEncode(Base64(MD5(RequestData + ApiKey)))
func Sign(requestData, appKey string) string {
	// 1. MD5加密
	h := md5.New()
	h.Write([]byte(requestData + appKey))
	md5Str := hex.EncodeToString(h.Sum(nil))

	// 2. Base64编码
	base64Str := base64.StdEncoding.EncodeToString([]byte(md5Str))

	// 3. URL编码
	encodedStr := url.QueryEscape(base64Str)

	return encodedStr
}

// VerifySign 验证签名
func VerifySign(requestData, appKey, dataSign string) bool {
	expectedSign := Sign(requestData, appKey)
	return expectedSign == dataSign
}