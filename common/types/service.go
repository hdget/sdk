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

type ServiceOperation[TKey KeyType, TCreate, TUpdate, TFilter, TResult any] interface {
	ServiceCreate[TCreate, TKey]
	ServiceGet[TKey, TResult]
	ServiceQuery[TFilter, TResult]
	ServiceUpdate[TUpdate]
	ServiceDelete[TKey]
}

// ============================================================
// 范围操作接口
// ============================================================

type ScopedServiceCreate[TScopedKey, TKey KeyType, TCreate any] interface {
	Create(ctx context.Context, scopedKey TScopedKey, model TCreate) (TKey, error) // 创建业务对象
}

type ScopedServiceGet[TScopedKey, TKey KeyType, TResult any] interface {
	Get(ctx context.Context, scopedKey TScopedKey, itemKey TKey) (TResult, error) // 获取业务对象
}

type ScopedServiceQuery[TScopedKey KeyType, TFilter, TResult any] interface {
	Query(ctx context.Context, scopedKey TScopedKey, filters TFilter, list ...*protobuf.ListParam) (int64, []TResult, error) // 查询业务对象
}

type ScopedServiceUpdate[TScopedKey KeyType, TUpdate any] interface {
	Update(ctx context.Context, scopedKey TScopedKey, model TUpdate) error // 更新业务对象
}

type ScopedServiceDelete[TScopedKey, TKey KeyType] interface {
	Delete(ctx context.Context, scopedKey TScopedKey, itemKey TKey) error
}

type ScopedServiceOperation[TScopedKey, TKey KeyType, TCreate, TUpdate, TFilter, TResult any] interface {
	ScopedServiceCreate[TScopedKey, TKey, TCreate]
	ScopedServiceGet[TScopedKey, TKey, TResult]
	ScopedServiceQuery[TScopedKey, TFilter, TResult]
	ScopedServiceUpdate[TScopedKey, TUpdate]
	ScopedServiceDelete[TScopedKey, TKey]
}
