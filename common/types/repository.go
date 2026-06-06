package types

import (
	"context"

	"github.com/hdget/sdk/common/protobuf"
)

// ============================================================
// 基础 CRUD 接口
// ============================================================

// RepoCreate 创建操作:C
type RepoCreate[TCreate any, TKey KeyType] interface {
	Create(ctx context.Context, model TCreate) (TKey, error)
}

type RepoGet[TKey KeyType, TModel any] interface {
	Get(ctx context.Context, key TKey) (TModel, error)
}

type RepoCount[TFilter any] interface {
	Count(ctx context.Context, filter TFilter) (int64, error)
}

type RepoList[TFilter any, TModel any] interface {
	List(ctx context.Context, filter TFilter, list ...*protobuf.ListParam) ([]TModel, error)
}

type RepoQuery[TFilter any, TModel any] interface {
	RepoCount[TFilter]
	RepoList[TFilter, TModel]
}

// RepoUpdate 更新:U
type RepoUpdate[TUpdate any] interface {
	Update(ctx context.Context, model TUpdate) error
}

// RepoDelete 删除:D
type RepoDelete[TKey KeyType] interface {
	Delete(ctx context.Context, keys ...TKey) (int64, error)
}

// ============================================================
// 批量操作接口
// ============================================================

// RepoBulkGet 批量读取（返回 map）
// 调用者需要 slice 时，使用 pie.Values(map) 或 slices.Collect(maps.Values(map))
type RepoBulkGet[TKey KeyType, TModel any] interface {
	BulkGet(ctx context.Context, keys ...TKey) (map[TKey]TModel, error)
}

// ============================================================
// 组合接口
// ============================================================

type RepoOperation[TKey KeyType, TCreate, TUpdate, TFilter, TModel any] interface {
	RepoCreate[TCreate, TKey]
	RepoGet[TKey, TModel]
	RepoCount[TFilter]
	RepoList[TFilter, TModel]
	RepoUpdate[TUpdate]
	RepoDelete[TKey]
}

// ============================================================
// 带 scope的关联资源 Repo
// ============================================================

type ScopedRepoCreate[TScopeKey, TKey KeyType, TCreate any] interface {
	Create(ctx context.Context, scopeKey TScopeKey, item TCreate) (TKey, error)
}

type ScopedRepoUpdate[TScopeKey KeyType, TUpdate any] interface {
	Update(ctx context.Context, scopeKey TScopeKey, item TUpdate) error
}

type ScopedRepoDelete[TScopeKey, TKey KeyType] interface {
	Delete(ctx context.Context, scopeKey TScopeKey, keys ...TKey) (int64, error)
}

type ScopedRepoCount[TScopeKey KeyType, TFilter any] interface {
	Count(ctx context.Context, scopeKey TScopeKey, filters TFilter) (int64, error)
}

type ScopedRepoList[TScopeKey KeyType, TFilter, TModel any] interface {
	List(ctx context.Context, scopeKey TScopeKey, filters TFilter, list ...*protobuf.ListParam) ([]TModel, error)
}

type ScopedRepoOperation[TScopeKey, TKey KeyType, TCreate, TUpdate, TFilter, TModel any] interface {
	RepoGet[TKey, TModel]
	ScopedRepoCreate[TScopeKey, TKey, TCreate]
	ScopedRepoUpdate[TScopeKey, TUpdate]
	ScopedRepoDelete[TScopeKey, TKey]
	ScopedRepoCount[TScopeKey, TFilter]
	ScopedRepoList[TScopeKey, TFilter, TModel]
}
