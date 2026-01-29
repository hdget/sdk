package types

import "time"

type OssProvider interface {
	Provider
	WithContentTypes(contentTypes []string) OssProvider                                   // 设置允许的文件类型
	WithMaxFileSize(fileSize int64) OssProvider                                           // 设置允许最大的文件大小
	WithSignExpiresIn(duration time.Duration) OssProvider                                 // 设置签名过期时间
	Upload(dir, filename string, data []byte) (string, error)                             // 上传文件
	GetPresignedURL(dir, filename, contentType string) (string, map[string]string, error) // 生成预签名URL, 返回URL,headers
	GetPostSignature(dir, filename string) (map[string]string, error)                     // 生成POST签名
}
