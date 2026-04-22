package aliyun

import (
	"time"

	"github.com/hdget/sdk/libs/oss"

	alisdk "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type aliyunOssImpl struct {
	config            oss.Config
	client            *alisdk.Client
	allowContentTypes []string
	maxFileSize       int64
	signExpiresIn     time.Duration
	objectACL         oss.ObjectACL
}

const (
	defaultSignatureExpiresIn = 180 * time.Second        // 上传签名默认失效时间, 3分钟
	defaultMaxFileSize        = int64(100 * 1024 * 1024) // 上传文件的最大尺寸, 100M
)

func New(cfg oss.Config, options ...oss.Option) (oss.API, error) {
	impl := &aliyunOssImpl{
		config:            cfg,
		allowContentTypes: oss.ImageContentTypes,     // 默认允许图片文件上传
		maxFileSize:       defaultMaxFileSize,        // 默认文件上传大小为100M
		signExpiresIn:     defaultSignatureExpiresIn, // 默认签名过期时间为3分钟
		objectACL:         oss.ACLDefault,            // 默认继承Bucket权限
	}

	for _, option := range options {
		option(impl)
	}

	impl.client = newClient(cfg)

	return impl, nil
}

func newClient(cfg oss.Config) *alisdk.Client {
	// 构建凭证提供者
	credProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AccessSecret)

	// 创建OSS配置
	c := alisdk.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithRegion(cfg.Region) // region: cn-shanghai, 不需要带oss

	return alisdk.NewClient(c)
}
