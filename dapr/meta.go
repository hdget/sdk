package dapr

import (
	"context"
	"github.com/hdget/sdk/encoding"
	"google.golang.org/grpc/metadata"
)

const (
	MetaKeyClient  = "hd-client"
	MetaKeyRelease = "hd-release"
	MetaKeyTsn     = "hd-tsn"  // tenant sn
	MetaKeyUsn     = "hd-usn"  // user sn
	MetaKeyRids    = "hd-rids" // encoded role ids
	MetaKeyCaller  = "dapr-caller-app-id"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
		MetaKeyTsn,
		MetaKeyUsn,
		MetaKeyRids,
		MetaKeyClient,
		MetaKeyRelease,
	}
)

type MetaManager interface {
	GetHttpHeaderKeys() []string
	GetClient(ctx context.Context) string
	GetRelease(ctx context.Context) string
	GetUsn(ctx context.Context) string // 获取用户的sn
	GetTsn(ctx context.Context) string // 获取租户的sn
	GetRoleIds(ctx context.Context) []int64
	GetCaller(ctx context.Context) string
	// DEPRECATED
	OldGetRoles(ctx context.Context) []*Role
	OldGetRoleValues(ctx context.Context) []string
	OldGetRoleIds(ctx context.Context) []int64
	OldGetPermIds(ctx context.Context) []int64
}

type metaManagerImpl struct {
}

func Meta() MetaManager {
	return &metaManagerImpl{}
}

func (m metaManagerImpl) GetClient(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyClient)
}

func (m metaManagerImpl) GetHttpHeaderKeys() []string {
	return _httpHeaderKeys
}

func (m metaManagerImpl) GetRelease(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyRelease)
}

func (m metaManagerImpl) GetCaller(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyCaller)
}

func (m metaManagerImpl) GetRoleIds(ctx context.Context) []int64 {
	return encoding.New().DecodeInt64Slice(getGrpcMdFirstValue(ctx, MetaKeyRids))
}

func (m metaManagerImpl) GetUsn(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyUsn)
}

func (m metaManagerImpl) GetTsn(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyTsn)
}

// getGrpcMdFirstValue get grpc metadata first value
func getGrpcMdFirstValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// getGrpcMdValues get grpc meta all values
func getGrpcMdValues(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}
