package oss_aliyun

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/pkg/errors"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (p *aliyunOssProvider) Upload(dir, filename string, data []byte) (string, error) {
	objectKey := p.getObjectKey(dir, filename)

	putRequest := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(p.config.Bucket), // 存储空间名称
		Key:          oss.Ptr(objectKey),       // 存储对象路径
		Body:         bytes.NewReader(data),
		StorageClass: oss.StorageClassStandard, // 指定对象的存储类型为标准存储
		Acl:          oss.ObjectACLPublicRead,  // 指定对象的访问权限
	}

	// 执行上传对象的请求
	_, err := p.newOSSClient().PutObject(context.TODO(), putRequest)
	if err != nil {
		return "", errors.Wrapf(err, "oss put object, dir: %s, filename: %s", dir, filename)
	}

	return objectKey, nil
}

func (p *aliyunOssProvider) getObjectKey(dir, filename string) string {
	strDate := time.Now().Format("20060102")
	year, month, day := strDate[:4], strDate[4:6], strDate[6:8]
	return path.Join(dir, year, month, day, p.generateSafeFileName(filename))
}

func (p *aliyunOssProvider) generateSafeFileName(filename string) string {
	safeFileName := filepath.Base(filename)                   // 移除路径分隔符
	safeFileName = strings.ReplaceAll(safeFileName, " ", "_") // 替换空格等特殊字符

	ext := filepath.Ext(safeFileName)
	name := safeFileName[0 : len(safeFileName)-len(ext)]

	return fmt.Sprintf("%s_%s%s", name, randStr(6), ext) // 防止相同文件名被覆盖
}
