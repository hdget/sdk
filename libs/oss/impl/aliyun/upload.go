package aliyun

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	alisdk "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/hdget/sdk/libs/oss"
	"github.com/pkg/errors"
)

func (impl *aliyunOssImpl) Upload(ctx context.Context, dir, filename string, data []byte) (string, error) {
	// 路径遍历防护：检查dir和filename参数
	if err := validatePath(dir, filename); err != nil {
		return "", err
	}

	objectKey := impl.getObjectKey(dir, filename)

	putRequest := &alisdk.PutObjectRequest{
		Bucket:       alisdk.Ptr(impl.config.Bucket), // 存储空间名称
		Key:          alisdk.Ptr(objectKey),          // 存储对象路径
		Body:         bytes.NewReader(data),
		StorageClass: alisdk.StorageClassStandard, // 指定对象的存储类型为标准存储
		Acl:          impl.getObjectACL(),         // 指定对象的访问权限
	}

	// 执行上传对象的请求，使用传入的context
	_, err := impl.client.PutObject(ctx, putRequest)
	if err != nil {
		return "", errors.Wrapf(err, "oss put object, dir: %s, filename: %s", dir, filename)
	}

	return objectKey, nil
}

// validatePath 验证路径参数，防止路径遍历攻击
func validatePath(dir, filename string) error {
	// 检查目录路径
	if strings.Contains(dir, "..") {
		return errors.New("invalid directory path: path traversal detected")
	}

	// 检查文件名
	if strings.Contains(filename, "..") {
		return errors.New("invalid filename: path traversal detected")
	}

	// 使用filepath.Base清理文件名，确保不包含路径分隔符
	cleanedFilename := filepath.Base(filename)
	if cleanedFilename != filename && strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		// 文件名包含路径分隔符，可能尝试路径遍历
		return errors.New("invalid filename: path separators not allowed")
	}

	return nil
}

func (impl *aliyunOssImpl) getObjectACL() alisdk.ObjectACLType {
	switch impl.objectACL {
	case oss.ACLPrivate:
		return alisdk.ObjectACLPrivate
	case oss.ACLPublicRead:
		return alisdk.ObjectACLPublicRead
	case oss.ACLPublicRW:
		return alisdk.ObjectACLPublicReadWrite
	default:
		return alisdk.ObjectACLDefault // 继承Bucket权限
	}
}

func (impl *aliyunOssImpl) getObjectKey(dir, filename string) string {
	strDate := time.Now().Format("20060102")
	year, month, day := strDate[:4], strDate[4:6], strDate[6:8]
	return path.Join(dir, year, month, day, impl.generateSafeFileName(filename))
}

func (impl *aliyunOssImpl) generateSafeFileName(filename string) string {
	safeFileName := filepath.Base(filename)                   // 移除路径分隔符
	safeFileName = strings.ReplaceAll(safeFileName, " ", "_") // 替换空格等特殊字符

	ext := filepath.Ext(safeFileName)
	name := safeFileName[0 : len(safeFileName)-len(ext)]

	return fmt.Sprintf("%s_%s%s", name, randStr(6), ext) // 防止相同文件名被覆盖
}
