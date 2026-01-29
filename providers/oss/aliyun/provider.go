package oss_aliyun

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/hdget/sdk/common/types"
	"time"
)

type aliyunOssProvider struct {
	config            *aliyunOssConfig
	allowContentTypes []string
	maxFileSize       int64
	signExpiresIn     time.Duration
}

const (
	defaultSignatureExpiresIn = 180 * time.Second        // 上传签名默认失效时间, 3分钟
	defaultMaxFileSize        = int64(100 * 1024 * 1024) // 上传文件的最大尺寸, 100M
)

var (
	// ImageContentTypes 图像类
	ImageContentTypes = []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/bmp",
		"image/webp",
		"image/svg+xml",
		"image/webp",
		"image/tiff",
		"image/vnd.microsoft.icon",
	}

	// VideoContentTypes 视频类
	VideoContentTypes = []string{
		"video/mp4",
		"video/mpeg",
		"video/ogg",
		"video/webm",
		"video/quicktime",
		"video/x-msvideo",
		"video/x-ms-wmv",
	}

	// DocumentContentTypes 文档类
	DocumentContentTypes = []string{
		"text/plain",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/html",
		"application/json",
	}

	// ZipContentTypes 压缩类
	ZipContentTypes = []string{
		"application/zip",
		"application/gzip",
		"application/x-tar",
		"application/x-rar-compressed",
	}
)

func New(configProvider types.ConfigProvider) (types.OssProvider, error) {
	config, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	return &aliyunOssProvider{
		config:            config,
		allowContentTypes: ImageContentTypes,         // 默认允许图片文件上传
		maxFileSize:       defaultMaxFileSize,        // 默认文件上传大小为100M
		signExpiresIn:     defaultSignatureExpiresIn, // 默认签名过期时间为3分钟
	}, nil
}

func (p *aliyunOssProvider) GetCapability() types.Capability {
	return Capability
}

// WithContentTypes 设置允许的文件类型
func (p *aliyunOssProvider) WithContentTypes(contentTypes []string) types.OssProvider {
	if len(contentTypes) > 0 {
		p.allowContentTypes = contentTypes
	}
	return p
}

// WithMaxFileSize 设置允许最大的文件大小
func (p *aliyunOssProvider) WithMaxFileSize(size int64) types.OssProvider {
	if size > 0 {
		p.maxFileSize = size
	}
	return p
}

// WithSignExpiresIn 设置签名过期时间，单位为秒
func (p *aliyunOssProvider) WithSignExpiresIn(expiresIn time.Duration) types.OssProvider {
	if expiresIn > 0 {
		p.signExpiresIn = expiresIn
	}
	return p
}

func (p *aliyunOssProvider) newOSSClient() *oss.Client {
	// 构建凭证提供者
	credProvider := credentials.NewStaticCredentialsProvider(p.config.AccessKey, p.config.AccessSecret)

	// 创建OSS配置
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithRegion(p.config.Region) // region: cn-shanghai, 不需要带oss

	return oss.NewClient(ossCfg)
}
