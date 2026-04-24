package biz

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

// hdCtxValueKey 定义 context key 类型，避免与其他包冲突
type hdCtxValueKey struct{}

// ctxValue 存储在 context 中的内部结构体
// 包含 metadata、Transactor 和缓存字段
type ctxValue struct {
	metadata   MetaData
	transactor Transactor
}

// NewContext 创建一个kvs的 context.Context，包含空的 metadata 和 Transactor
func NewContext(kvs ...any) context.Context {
	return context.WithValue(context.Background(), hdCtxValueKey{}, &ctxValue{
		metadata:   newMetaData(kvs...),
		transactor: newTransactor(),
	})
}

// NewFromIncomingGrpcContext 从 gRPC context 中提取 metadata 并存储到 context.Context
func NewFromIncomingGrpcContext(ctx context.Context) context.Context {
	cv := &ctxValue{
		metadata:   newMetaData(),
		transactor: newTransactor(),
	}

	// 尝试从 gRPC context 中获取 metadata
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return context.WithValue(ctx, hdCtxValueKey{}, cv)
	}

	for key, values := range md {
		switch key {
		case MetaKeyTid, MetaKeyUid: // int64
			val, err := strconv.ParseInt(values[0], 10, 64)
			if err == nil {
				cv.metadata.Set(key, val)
			}
		case MetaKeyRoleIds:
			var val []int64
			if values[0] != "" {
				strIds := strings.Split(values[0], ",")
				val = make([]int64, 0, len(strIds))
				for _, s := range strIds {
					id, err := strconv.ParseInt(s, 10, 64)
					if err == nil {
						val = append(val, id)
					}
				}
			}
			cv.metadata.Set(key, val)
		default:
			cv.metadata.Set(key, values[0])
		}
	}

	return context.WithValue(ctx, hdCtxValueKey{}, cv)
}

// NewOutgoingGrpcContext 将 context 中的 metadata 转换为 gRPC outgoing context
func NewOutgoingGrpcContext(ctx context.Context) context.Context {
	cv := getCtxValue(ctx)
	if cv == nil {
		return metadata.NewOutgoingContext(ctx, newMetaData().AsGRPCMetaData())
	}
	return metadata.NewOutgoingContext(ctx, cv.metadata.AsGRPCMetaData())
}

// GetMetaData 从 context 中获取 MetaData
func GetMetaData(ctx context.Context) MetaData {
	cv := getCtxValue(ctx)
	if cv == nil {
		return newMetaData()
	}
	return cv.metadata
}

// GetTransactor 从 context 中获取 Transactor
func GetTransactor(ctx context.Context) Transactor {
	cv := getCtxValue(ctx)
	if cv == nil {
		return newTransactor()
	}
	return cv.transactor
}

// GetTid 从 context 中获取租户 ID（带缓存优化）
func GetTid(ctx context.Context) int64 {
	cv := getCtxValue(ctx)
	if cv == nil {
		return 0
	}
	// 从 metadata 获取，metadata 内部已有缓存机制
	return cv.metadata.GetInt64(MetaKeyTid)
}

// GetUid 从 context 中获取用户 ID（带缓存优化）
func GetUid(ctx context.Context) int64 {
	cv := getCtxValue(ctx)
	if cv == nil {
		return 0
	}
	return cv.metadata.GetInt64(MetaKeyUid)
}

// GetAppId 从 context 中获取应用 app_id（带缓存优化）
func GetAppId(ctx context.Context) string {
	cv := getCtxValue(ctx)
	if cv == nil {
		return ""
	}
	return cv.metadata.GetString(MetaKeyAppKey)
}

// GetSource 从 context 中获取请求来源（带缓存优化）
func GetSource(ctx context.Context) string {
	cv := getCtxValue(ctx)
	if cv == nil {
		return ""
	}
	return cv.metadata.GetString(MetaKeySource)
}

// GetAppCode 从 context 中获取请求应用类型标识（带缓存优化）
func GetAppCode(ctx context.Context) string {
	cv := getCtxValue(ctx)
	if cv == nil {
		return ""
	}
	return cv.metadata.GetString(MetaKeyAppCode)
}

// RoleIds 从 context 中获取角色 ID 列表（带缓存优化）
func RoleIds(ctx context.Context) []int64 {
	cv := getCtxValue(ctx)
	if cv == nil {
		return nil
	}
	return cv.metadata.GetInt64Slice(MetaKeyRoleIds)
}

// WithTid 将租户 ID 存入 context
func WithTid(ctx context.Context, tid int64) context.Context {
	cv := mustGetCtxValue(ctx)
	cv.metadata.Set(MetaKeyTid, tid)
	return context.WithValue(ctx, hdCtxValueKey{}, cv)
}

// WithUid 将用户 ID 存入 context
func WithUid(ctx context.Context, uid int64) context.Context {
	cv := mustGetCtxValue(ctx)
	cv.metadata.Set(MetaKeyUid, uid)
	return context.WithValue(ctx, hdCtxValueKey{}, cv)
}

// WithRoleIds 将角色 ID 列表存入 context
func WithRoleIds(ctx context.Context, roleIds []int64) context.Context {
	cv := mustGetCtxValue(ctx)
	cv.metadata.Set(MetaKeyRoleIds, roleIds)
	return context.WithValue(ctx, hdCtxValueKey{}, cv)
}

// WithMetaData 将 metadata 存入 context
func WithMetaData(ctx context.Context, md MetaData) context.Context {
	cv := getCtxValue(ctx)
	if cv == nil {
		cv = &ctxValue{
			metadata:   md,
			transactor: newTransactor(),
		}
	} else {
		cv.metadata = md
	}
	return context.WithValue(ctx, hdCtxValueKey{}, cv)
}

// getCtxValue 从 context 中获取 ctxValue，如果不存在返回 nil
func getCtxValue(ctx context.Context) *ctxValue {
	if cv, ok := ctx.Value(hdCtxValueKey{}).(*ctxValue); ok {
		return cv
	}
	return nil
}

// mustGetCtxValue 确保有个合法的*ctxValue
func mustGetCtxValue(ctx context.Context) *ctxValue {
	cv := getCtxValue(ctx)
	if cv == nil {
		cv = &ctxValue{
			metadata:   newMetaData(),
			transactor: newTransactor(),
		}
	}
	return cv
}
