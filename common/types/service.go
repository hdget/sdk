package types

import (
	"github.com/hdget/sdk/common/biz"
	"github.com/hdget/sdk/common/protobuf"
)

type ServiceOperation[TObjectId ObjectIdentifier, TBizObject any] interface {
	ServiceCreate[TBizObject]
	ServiceRetrieve[TObjectId, TBizObject]
	ServiceUpdate[TBizObject]
	ServiceDelete[TObjectId]
}

type ServiceCreate[TBizObject any] interface {
	Create(ctx biz.Context, object TBizObject) (int64, error) // 创建业务对象
}

type ServiceRetrieve[TObjectId ObjectIdentifier, TBizObject any] interface {
	Get(ctx biz.Context, id TObjectId) (TBizObject, error)                                                      // 获取业务对象
	Query(ctx biz.Context, filters map[string]string, list ...*protobuf.ListParam) (int64, []TBizObject, error) // 查询业务对象
}

type ServiceUpdate[TBizObject any] interface {
	Edit(ctx biz.Context, bizObject TBizObject) error // 编辑业务对象
}

type ServiceDelete[TObjectId ObjectIdentifier] interface {
	Delete(ctx biz.Context, id TObjectId) error // 删除业务对象
}
