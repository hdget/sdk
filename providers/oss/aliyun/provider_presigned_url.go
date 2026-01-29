package oss_aliyun

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
)

// GetPresignedURL 生成PutObject的预签名URL
func (p *aliyunOssProvider) GetPresignedURL(dir, filename, contentType string) (string, map[string]string, error) {
	if !pie.Contains(p.allowContentTypes, contentType) {
		return "", nil, errors.New("content type not allowed")
	}

	if dir == "" || filename == "" {
		return "", nil, errors.New("dir or filename is empty")
	}

	objectKey := p.getObjectKey(dir, filename)

	result, err := p.newOSSClient().Presign(context.TODO(), &oss.PutObjectRequest{
		Bucket:       oss.Ptr(p.config.Bucket),
		Key:          oss.Ptr(objectKey),
		ContentType:  oss.Ptr(contentType),
		StorageClass: oss.StorageClassStandard,
	}, oss.PresignExpires(p.signExpiresIn))
	if err != nil {
		return "", nil, errors.Wrapf(err, "presign, dir: %s, filename: %s", dir, filename)
	}

	return result.URL, result.SignedHeaders, nil
}
