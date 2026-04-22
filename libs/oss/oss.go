package oss

import "context"

// ObjectACL 对象访问权限
type ObjectACL string

const (
	ACLPrivate      ObjectACL = "private"           // 私有读写
	ACLPublicRead   ObjectACL = "public-read"       // 公共读，私有写
	ACLPublicRW     ObjectACL = "public-read-write" // 公共读写
	ACLDefault      ObjectACL = "default"           // 继承Bucket权限
)

// API object storage service api
type API interface {
	Upload(ctx context.Context, dir, filename string, data []byte) (string, error)                                          // 上传文件
	GetPresignedURL(ctx context.Context, dir, filename, contentType string) (string, map[string]string, error)             // 生成预签名URL, 返回URL,headers
	GetPostSignature(ctx context.Context, dir, filename string) (map[string]string, error)                                 // 生成POST签名
}

// Option 配置选项函数
type Option func(API)
