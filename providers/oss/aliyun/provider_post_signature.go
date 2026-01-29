package oss_aliyun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"path"
	"time"
)

// GetPostSignature 生成oss直传post签名
func (p *aliyunOssProvider) GetPostSignature(dir, filename string) (map[string]string, error) {
	now := time.Now()
	ossDate := now.UTC().Format("20060102T150405Z")
	ossCredential := path.Join(p.config.AccessKey, time.Now().UTC().Format("20060102"), p.config.Region, "oss", "aliyun_v4_request")

	policyBase64, policySigned, err := p.generatePolicy(dir, now, ossDate, ossCredential)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"key":                     p.getObjectKey(dir, filename), // 返回自定义的Object名字给前端
		"policy":                  policyBase64,
		"x-oss-signature":         policySigned,
		"x-oss-credential":        ossCredential,
		"x-oss-access-key-id":     p.config.AccessKey,
		"x-oss-signature-version": "OSS4-HMAC-SHA256",
		"x-oss-date":              ossDate,
	}, nil
}

// generatePolicy 生成访问策略
func (p *aliyunOssProvider) generatePolicy(dir string, now time.Time, ossDate, ossCredential string) (string, string, error) {
	// 定义策略
	policy := map[string]any{
		// 多少秒后签名过期
		"expiration": now.Add(p.signExpiresIn).Format("2006-01-02T15:04:05Z"),
		"conditions": []any{
			map[string]string{"bucket": p.config.Bucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"}, // 签名的版本和算法
			map[string]string{"x-oss-credential": ossCredential},             // 派生密钥的参数集
			map[string]string{"x-oss-date": ossDate},                         // 请求的时间，其格式遵循ISO 8601日期和时间标准，例如20231203T121212Z
			[]any{"starts-with", "$key", dir},                                // 限制上传目录， 上传的文件名必须以dir开头
			[]any{"content-length-range", 0, p.maxFileSize},                  // 文件大小限制
			[]any{"in", "$content-type", p.allowContentTypes},                // 文件内容限制
		},
	}

	// 编码Policy为Base64
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return "", "", err
	}

	policyBase64 := base64.StdEncoding.EncodeToString(policyJSON)

	// 计算SigningKey (HMAC阶梯)
	// 第一级：日期密钥
	dateKey := hmacSHA256([]byte("aliyun_v4"+p.config.AccessSecret), now.UTC().Format("20060102"))
	// 第二级：区域密钥
	dateRegionKey := hmacSHA256(dateKey, p.config.Region)
	// 第三级：服务密钥
	dateRegionServiceKey := hmacSHA256(dateRegionKey, "oss")
	// 最终签名密钥
	signingKey := hmacSHA256(dateRegionServiceKey, "aliyun_v4_request")

	// 最终计算Signature
	policySigned := hex.EncodeToString(hmacSHA256(signingKey, policyBase64))

	return policyBase64, policySigned, nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
