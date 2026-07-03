package kdniao

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
)

// sign 生成快递鸟签名
// 签名算法: DataSign = Base64(MD5(RequestData + ApiKey))
//
// 安全说明: MD5 仅用于满足快递鸟 ApiCommon 的签名要求，不应用于其他安全敏感场景。
// 参考: 快递鸟开放平台文档
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
// 使用 constant time 比较防止时序攻击
func verifySign(requestData, appKey, dataSign string) bool {
	expectedSign := sign(requestData, appKey)
	// 使用 constant time 比较防止时序攻击
	return subtle.ConstantTimeCompare([]byte(expectedSign), []byte(dataSign)) == 1
}
