package meta

import (
	"context"

	"github.com/hdget/common/constant"
	"github.com/hdget/utils/convert"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

type Meta interface {
	GetAppId(ctx context.Context) string
	GetClient(ctx context.Context) string
	GetRelease(ctx context.Context) string
	GetUid(ctx context.Context) int64  // 获取用户的id
	GetTsn(ctx context.Context) string // 获取租户的sn
	GetTid(ctx context.Context) int64  // 获取租户的id
	GetRoleIds(ctx context.Context) []int64
	GetCaller(ctx context.Context) string
	// DEPRECATED
	OldGetRoles(ctx context.Context) []*Role
	OldGetRoleValues(ctx context.Context) []string
	OldGetRoleIds(ctx context.Context) []int64
	OldGetPermIds(ctx context.Context) []int64
}

type metaImpl struct {
}

func New() Meta {
	return &metaImpl{}
}

func (m metaImpl) GetAppId(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, constant.MetaKeyAppId)
}

func (m metaImpl) GetClient(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, constant.MetaKeyClient)
}

func (m metaImpl) GetRelease(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, constant.MetaKeyRelease)
}

func (m metaImpl) GetCaller(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, constant.MetaKeyCaller)
}

func (m metaImpl) GetRoleIds(ctx context.Context) []int64 {
	return convert.CsvToInt64s(getGrpcMdFirstValue(ctx, constant.MetaKeyRoleIds))
}

func (m metaImpl) GetUid(ctx context.Context) int64 {
	return cast.ToInt64(getGrpcMdFirstValue(ctx, constant.MetaKeyUid))
}

func (m metaImpl) GetTsn(ctx context.Context) string {
	return getGrpcMdFirstValue(ctx, constant.MetaKeyTsn)
}

func (m metaImpl) GetTid(ctx context.Context) int64 {
	return cast.ToInt64(getGrpcMdFirstValue(ctx, constant.MetaKeyTid))
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
