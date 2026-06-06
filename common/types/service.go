package types

import (
	"context"

	"github.com/hdget/sdk/common/protobuf"
)

type ServiceCreate[TCreate any, TKey KeyType] interface {
	Create(ctx context.Context, model TCreate) (TKey, error) // 创建业务对象
}

type ServiceGet[TKey KeyType, TResult any] interface {
	Get(ctx context.Context, key TKey) (TResult, error) // 获取业务对象
}

type ServiceQuery[TFilter, TResult any] interface {
	Query(ctx context.Context, filters TFilter, list ...*protobuf.ListParam) (int64, []TResult, error) // 查询业务对象
}

type ServiceUpdate[TUpdate any] interface {
	Update(ctx context.Context, model TUpdate) error // 更新业务对象
}

type ServiceDelete[TKey KeyType] interface {
	Delete(ctx context.Context, key TKey) error
}

// ============================================================
// 批量操作接口
// ============================================================

// ServiceBulkGet 批量读取（返回 map）
// 调用者需要 slice 时，使用 pie.Values(map) 或 slices.Collect(maps.Values(map))
type ServiceBulkGet[TKey KeyType, TResult any] interface {
	BulkGet(ctx context.Context, keys ...TKey) (map[TKey]TResult, error)
}

// ============================================================
// 组合接口
// ============================================================

type ServiceOperation[TKey KeyType, TCreate, TUpdate, TFilter, TResult any] interface {
	ServiceCreate[TCreate, TKey]
	ServiceGet[TKey, TResult]
	ServiceQuery[TFilter, TResult]
	ServiceUpdate[TUpdate]
	ServiceDelete[TKey]
}
