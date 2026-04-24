package types

import (
	"context"

	"github.com/hdget/sdk/common/protobuf"
)

type ServiceOperation[TObjectKey ObjectKeyType, TBizObject any] interface {
	ServiceCreate[TBizObject]
	ServiceRetrieve[TObjectKey, TBizObject]
	ServiceUpdate[TBizObject]
	ServiceDelete[TObjectKey]
}

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

type ServiceUpdate[TBizObject any] interface {
	Edit(ctx context.Context, bizObject TBizObject) error // 编辑业务对象
}

type ServiceDelete[TObjectKey ObjectKeyType] interface {
	Delete(ctx context.Context, key TObjectKey) error // 删除业务对象
}
