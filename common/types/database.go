package types

import (
	"github.com/hdget/sdk/common/protobuf"
)

// Database

type DbOperation[TObjectId ObjectIdentifier, TBizObject any, TModelObject any, TCondition any] interface {
	DbCreate[TBizObject, TModelObject]
	DbRetrieve[TObjectId, TModelObject, TCondition]
	DbUpdate[TModelObject]
	DbDelete[TObjectId]
}

// DbCreate 创建操作:C
type DbCreate[BizObject any, ModelObject any] interface {
	Create(bizObj BizObject) (ModelObject, error) // 创建对象
}

// DbRetrieve 读取操作:R
type DbRetrieve[TObjectId ObjectIdentifier, TModelObject any, TCondition any] interface {
	Get(id TObjectId) (TModelObject, error)                                              // 获取对象
	Count(filters map[string]string) (int64, error)                                      // 统计对象
	List(filters map[string]string, list ...*protobuf.ListParam) ([]TModelObject, error) // 列出对象, list不传的时候获取所有对象
	GetQueryConditions(filters map[string]string) []TCondition                           // 获取查询条件
}

// DbUpdate 更新：U
type DbUpdate[TModelObject any] interface {
	Update(modelObj TModelObject) error // 更新某个字段
}

type DbEdit[TBizObject any] interface {
	Edit(bizObj TBizObject) error // 编辑对象
}

// DbDelete 删除
type DbDelete[TObjectId ObjectIdentifier] interface {
	Delete(id TObjectId) error // 删除对象
}

// DbBulkRetrieve 批量读取
type DbBulkRetrieve[TObjectId ObjectIdentifier, ModelObject any] interface {
	BulkGet(ids []TObjectId) (map[TObjectId]ModelObject, error) // 批量获取对象
}

/* 关联数据表 */

type RefDbOperation[TObjectId ObjectIdentifier, TRefBizObject any, TRefModelObject any, TCondition any] interface {
	RefDbCreate[TObjectId, TRefBizObject, TRefModelObject]
	RefDbRetrieve[TObjectId, TRefModelObject, TCondition]
	RefDbUpdate[TRefModelObject]
	RefDbDelete[TObjectId]
}

// RefDbCreate 创建关联对象操作:C
type RefDbCreate[TObjectId ObjectIdentifier, TRefBizObject any, TRefModelObject any] interface {
	Create(id TObjectId, refBizObj TRefBizObject) (TRefModelObject, error) // 创建关联对象DAO
}

// RefDbRetrieve 读取关联对象操作:R
type RefDbRetrieve[TObjectId ObjectIdentifier, TRefModelObject any, Condition any] interface {
	Get(id, refId TObjectId) (TRefModelObject, error)                                                           // 获取关联对象DAO
	Count(id TObjectId, refObjFilters map[string]string) (int64, error)                                         // 统计关联对象DAO
	List(id TObjectId, refObjFilters map[string]string, list ...*protobuf.ListParam) ([]TRefModelObject, error) // 列出关联对象DAO
	GetQueryConditions(id TObjectId, refObjFilters map[string]string) []Condition                               // 获取关联对象DAO
}

// RefDbUpdate 更新关联对象：U
type RefDbUpdate[TRefObjectModel any] interface {
	Update(refObjModel TRefObjectModel) error // 更新数据库关联对象
}

type RefDbEdit[TObjectId ObjectIdentifier, TRefBizObject any] interface {
	Edit(id TObjectId, refBizObj TRefBizObject) error // 编辑数据库关联对象DAO
}

// RefDbDelete 删除关联对象
type RefDbDelete[TObjectId ObjectIdentifier] interface {
	Delete(id, refId TObjectId) error // 删除关联对象DAO
}

// RefDbBulkRetrieve 批量读取关联对象
type RefDbBulkRetrieve[TObjectId ObjectIdentifier, TRefModelObject any] interface {
	BulkGet(id TObjectId, refIds []TObjectId) (map[TObjectId]TRefModelObject, error) // 批量获取对象
}
