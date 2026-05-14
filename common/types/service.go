package types

import (
	"context"

	"github.com/hdget/sdk/common/protobuf"
)

type ServiceCreate[TBizObject any] interface {
	Create(ctx context.Context, object TBizObject) (int64, error) // 创建业务对象
}

type ServiceRetrieve[TObjectKey ObjectKeyType, TBizObject any] interface {
	ServiceGet[TObjectKey, TBizObject]
	ServiceQuery[TBizObject]
}

type ServiceGet[TObjectKey ObjectKeyType, TBizObject any] interface {
	Get(ctx context.Context, key TObjectKey) (TBizObject, error) // 获取业务对象
}

type ServiceQuery[TBizObject any] interface {
	Query(ctx context.Context, filters map[string]string, list ...*protobuf.ListParam) (int64, []TBizObject, error) // 查询业务对象
}

type ServiceEdit[TBizObject any] interface {
	Edit(ctx context.Context, bizObject TBizObject) error // 编辑业务对象
}

type ServiceUpdate[TModelObject any] interface {
	Update(ctx context.Context, modelObj TModelObject) error // 更新业务对象
}

type ServiceDelete[TObjectKey ObjectKeyType] interface {
	Delete(ctx context.Context, key TObjectKey) error // 删除业务对象
}

// ============================================================
// 批量操作接口
// ============================================================

// ServiceBulkGet 批量读取（返回 map）
// 调用者需要 slice 时，使用 pie.Values(map) 或 slices.Collect(maps.Values(map))
type ServiceBulkGet[TObjectKey ObjectKeyType, TModelObject any] interface {
	BulkGet(ctx context.Context, ids []TObjectKey) (map[TObjectKey]TModelObject, error)
}

// ServiceBulkCreate 批量创建
type ServiceBulkCreate[TModelObject any] interface {
	BulkCreate(ctx context.Context, items []TModelObject) ([]TModelObject, error)
}

// ============================================================
// 组合接口（常用 Service 动作组合）
// ============================================================

type ServiceOperation[TObjectKey ObjectKeyType, TBizObject any] interface {
	ServiceCreate[TBizObject]
	ServiceRetrieve[TObjectKey, TBizObject]
	ServiceEdit[TBizObject]
	ServiceDelete[TObjectKey]
}
