package types

import "github.com/hdget/sdk/common/protobuf"

type KeyType interface {
	~int64 | ~int32 | ~int | ~string
}

type GetRequest[TKey KeyType] struct {
	Key TKey `json:"key"`
}

type CreateRequest[TCreateModel any] struct {
	Item TCreateModel `json:"item"`
}

type DeleteRequest[TKey KeyType] struct {
	Key TKey `json:"key"`
}

type UpdateRequest[TUpdateModel any] struct {
	Item TUpdateModel `json:"item"`
}

type QueryRequest[TFilter any] struct {
	Filters TFilter             `json:"filters,omitempty"`
	List    *protobuf.ListParam `json:"list,omitempty"`
}

type QueryResponse[TResult any] struct {
	Total int64     `json:"total"`
	Items []TResult `json:"items"`
}

// ============================================================
// 范围操作接口
// ============================================================

type ScopedCreateRequest[TScopedKey KeyType, TCreate any] struct {
	ScopedKey TScopedKey `json:"scoped_key"`
	Item      TCreate    `json:"item"`
}

type ScopedUpdateRequest[TScopedKey KeyType, TUpdate any] struct {
	ScopedKey TScopedKey `json:"scoped_key"`
	Item      TUpdate    `json:"item"`
}

type ScopedDeleteRequest[TScopedKey, TKey KeyType] struct {
	ScopedKey TScopedKey `json:"scoped_key"`
	Key       TKey       `json:"key"`
}

type ScopedGetRequest[TScopedKey, TKey KeyType] struct {
	ScopedKey TScopedKey `json:"scoped_key"`
	Key       TKey       `json:"key"`
}

type ScopedQueryRequest[TScopedKey KeyType, TFilter any] struct {
	ScopedKey TScopedKey          `json:"scoped_key"`
	Filters   TFilter             `json:"filters,omitempty"`
	List      *protobuf.ListParam `json:"list,omitempty"`
}
