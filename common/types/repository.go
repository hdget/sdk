package types

import (
	"context"

	"github.com/hdget/sdk/common/protobuf"
)

// ============================================================
// 基础 CRUD 接口
// ============================================================

// RepoCreate 创建操作:C
type RepoCreate[TBizObject any, TModelObject any] interface {
	Create(ctx context.Context, bizObj TBizObject) (TModelObject, error)
}

// RepoRetrieve 读操作:R
type RepoRetrieve[TObjectKey ObjectKeyType, TModelObject any] interface {
	Get(ctx context.Context, key TObjectKey) (TModelObject, error)
	Count(ctx context.Context, filters map[string]string) (int64, error)
	List(ctx context.Context, filters map[string]string, list ...*protobuf.ListParam) ([]TModelObject, error)
}

// RepoUpdate 更新:U
type RepoUpdate[TModelObject any] interface {
	Update(ctx context.Context, modelObj TModelObject) error
}

// RepoEdit 编辑操作（基于业务对象，返回模型对象）
type RepoEdit[TBizObject any, TModelObject any] interface {
	Edit(ctx context.Context, bizObj TBizObject) (TModelObject, error)
}

// RepoDelete 删除:D
type RepoDelete[TObjectKey ObjectKeyType] interface {
	Delete(ctx context.Context, key TObjectKey) (int64, error)
}

// RepoUpsert Upsert 操作（创建或更新）
type RepoUpsert[TBizObject any, TModelObject any] interface {
	Upsert(ctx context.Context, bizObj TBizObject) (TModelObject, error)
}

// ============================================================
// 批量操作接口
// ============================================================

// RepoBulkGet 批量读取（返回 map）
// 调用者需要 slice 时，使用 pie.Values(map) 或 slices.Collect(maps.Values(map))
type RepoBulkGet[TObjectKey ObjectKeyType, TModelObject any] interface {
	BulkGet(ctx context.Context, ids []TObjectKey) (map[TObjectKey]TModelObject, error)
}

// RepoBulkCreate 批量创建
type RepoBulkCreate[TModelObject any] interface {
	BulkCreate(ctx context.Context, items []TModelObject) ([]TModelObject, error)
}

// ============================================================
// 组合接口（常用 Repo 动作组合）
// ============================================================

// RepoCRUD 标准 CRUD 操作接口
// TObjectKey: 主键类型（int64, int32, int, string）
// TBizObject: 业务对象类型（pb.XXX）
// TModelObject: 数据模型类型（db.XXX）
type RepoCRUD[TObjectKey ObjectKeyType, TBizObject any, TModelObject any] interface {
	RepoCreate[TBizObject, TModelObject]
	RepoRetrieve[TObjectKey, TModelObject]
	RepoEdit[TBizObject, TModelObject]
	RepoDelete[TObjectKey]
}

// RepoCRUDWithBulk 标准 CRUD + 批量读取
type RepoCRUDWithBulk[TObjectKey ObjectKeyType, TBizObject any, TModelObject any] interface {
	RepoCRUD[TObjectKey, TBizObject, TModelObject]
	RepoBulkGet[TObjectKey, TModelObject]
}

// RepoCRUDFull 完整 CRUD（包含所有常用操作）
type RepoCRUDFull[TObjectKey ObjectKeyType, TBizObject any, TModelObject any] interface {
	RepoCreate[TBizObject, TModelObject]
	RepoRetrieve[TObjectKey, TModelObject]
	RepoBulkGet[TObjectKey, TModelObject]
	RepoEdit[TBizObject, TModelObject]
	RepoDelete[TObjectKey]
}
