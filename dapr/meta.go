package dapr

import (
	"context"
	"github.com/hdget/sdk/encoding"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

const (
	MetaKeyClient  = "hd-client"
	MetaKeyRelease = "hd-release"
	MetaKeyTid     = "hd-tid"
	MetaKeyEtid    = "hd-etid"
	MetaKeyEuid    = "hd-euid"  // encoded user id
	MetaKeyErids   = "hd-erids" // encoded role ids
	MetaKeyCaller  = "dapr-caller-app-id"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
		MetaKeyEtid,
		MetaKeyClient,
		MetaKeyRelease,
	}
)

type MetaManager interface {
	GetHttpHeaderKeys() []string
	GetClientId(ctx context.Context) string
	GetRelease(ctx context.Context) string
	GetCaller(ctx context.Context) string
	GetUserId(ctx context.Context) int64
	GetRoleIds(ctx context.Context) []int64
	GetTenantId(ctx context.Context) int64
	GetEtid(ctx context.Context) string
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

func (m metaManagerImpl) GetClientId(ctx context.Context) string {
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
	return encoding.New().DecodeInt64Slice(getGrpcMdFirstValue(ctx, MetaKeyErids))
}

func (m metaManagerImpl) GetUserId(ctx context.Context) int64 {
	return encoding.New().DecodeInt64(getGrpcMdFirstValue(ctx, MetaKeyEuid))
}

func (m metaManagerImpl) GetTenantId(ctx context.Context) int64 {
	if v := getGrpcMdFirstValue(ctx, MetaKeyTid); v != "" {
		return cast.ToInt64(v)
	}
	return encoding.New().DecodeInt64(getGrpcMdFirstValue(ctx, MetaKeyEtid))
}

func (m metaManagerImpl) GetEtid(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, MetaKeyEtid)
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
