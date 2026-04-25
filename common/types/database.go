package types

import (
	"github.com/hdget/sdk/common/protobuf"
)

// Database

type DbOperation[TObjectKey ObjectKeyType, TBizObject any, TModelObject any, TCondition any] interface {
	DbCreate[TBizObject, TModelObject]
	DbRetrieve[TObjectKey, TModelObject, TCondition]
	DbUpdate[TModelObject]
	DbDelete[TObjectKey]
}

// DbCreate 创建操作:C
type DbCreate[BizObject any, ModelObject any] interface {
	Create(bizObj BizObject) (ModelObject, error) // 创建对象
}

// DbRetrieve 读操作
type DbRetrieve[TObjectKey ObjectKeyType, TModelObject any, TCondition any] interface {
	Get(key TObjectKey) (TModelObject, error)                                            // 获取对象
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
type DbDelete[TObjectKey ObjectKeyType] interface {
	Delete(key TObjectKey) error // 删除对象
}

// DbBulkRetrieve 批量读取
type DbBulkRetrieve[TObjectKey ObjectKeyType, ModelObject any] interface {
	BulkGet(ids []TObjectKey) (map[TObjectKey]ModelObject, error) // 批量获取对象
}

/* 关联数据表 */

type RefDbOperation[TObjectKey ObjectKeyType, TRefBizObject any, TRefModelObject any, TCondition any] interface {
	RefDbCreate[TObjectKey, TRefBizObject, TRefModelObject]
	RefDbRetrieve[TObjectKey, TRefModelObject, TCondition]
	RefDbUpdate[TRefModelObject]
	RefDbDelete[TObjectKey]
}

// RefDbCreate 创建关联对象操作:C
type RefDbCreate[TObjectKey ObjectKeyType, TRefBizObject any, TRefModelObject any] interface {
	Create(key TObjectKey, refBizObj TRefBizObject) (TRefModelObject, error) // 创建关联对象DAO
}

// RefDbRetrieve 读取关联对象操作:R
type RefDbRetrieve[TObjectKey ObjectKeyType, TRefModelObject any, Condition any] interface {
	Get(key, refKey TObjectKey) (TRefModelObject, error)                                                          // 获取关联对象DAO
	Count(key TObjectKey, refObjFilters map[string]string) (int64, error)                                         // 统计关联对象DAO
	List(key TObjectKey, refObjFilters map[string]string, list ...*protobuf.ListParam) ([]TRefModelObject, error) // 列出关联对象DAO
	GetQueryConditions(key TObjectKey, refObjFilters map[string]string) []Condition                               // 获取关联对象DAO
}

// RefDbUpdate 更新关联对象：U
type RefDbUpdate[TRefObjectModel any] interface {
	Update(refObjModel TRefObjectModel) error // 更新数据库关联对象
}

type RefDbEdit[TObjectKey ObjectKeyType, TRefBizObject any] interface {
	Edit(key TObjectKey, refBizObj TRefBizObject) error // 编辑数据库关联对象DAO
}

// RefDbDelete 删除关联对象
type RefDbDelete[TObjectKey ObjectKeyType] interface {
	Delete(key, refKey TObjectKey) error // 删除关联对象DAO
}

// RefDbBulkRetrieve 批量读取关联对象
type RefDbBulkRetrieve[TObjectKey ObjectKeyType, TRefModelObject any] interface {
	BulkGet(key TObjectKey, refKeys []TObjectKey) (map[TObjectKey]TRefModelObject, error) // 批量获取对象
}
