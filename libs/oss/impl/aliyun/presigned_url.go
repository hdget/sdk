package aliyun

import (
	"context"

	alisdk "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
)

// GetPresignedURL 生成PutObject的预签名URL
func (impl *aliyunOssImpl) GetPresignedURL(ctx context.Context, dir, filename, contentType string) (string, map[string]string, error) {
	if !pie.Contains(impl.allowContentTypes, contentType) {
		return "", nil, errors.New("content type not allowed")
	}

	if dir == "" || filename == "" {
		return "", nil, errors.New("dir or filename is empty")
	}

	objectKey := impl.getObjectKey(dir, filename)

	result, err := impl.client.Presign(ctx, &alisdk.PutObjectRequest{
		Bucket:       alisdk.Ptr(impl.config.Bucket),
		Key:          alisdk.Ptr(objectKey),
		ContentType:  alisdk.Ptr(contentType),
		StorageClass: alisdk.StorageClassStandard,
	}, alisdk.PresignExpires(impl.signExpiresIn))
	if err != nil {
		return "", nil, errors.Wrapf(err, "presign, dir: %s, filename: %s", dir, filename)
	}

	return result.URL, result.SignedHeaders, nil
}
