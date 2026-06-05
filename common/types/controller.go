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
	Keys []TKey `json:"keys"`
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
