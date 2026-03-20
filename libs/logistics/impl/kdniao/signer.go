package kdniao

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// sign 生成快递鸟签名
// 签名算法: DataSign = Base64(MD5(RequestData + ApiKey))
func sign(requestData, appKey string) string {
	// 1. MD5加密
	h := md5.New()
	h.Write([]byte(requestData + appKey))
	md5Str := hex.EncodeToString(h.Sum(nil))

	// 2. Base64编码
	base64Str := base64.StdEncoding.EncodeToString([]byte(md5Str))

	return base64Str
}

// verifySign 验证签名
func verifySign(requestData, appKey, dataSign string) bool {
	expectedSign := sign(requestData, appKey)
	return expectedSign == dataSign
}
