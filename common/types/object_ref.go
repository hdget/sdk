package types

import (
	"github.com/hdget/sdk/common/protobuf"
)

/* request */

type CreateRefObjectRequest[TObjectId ObjectIdentifier, TBizObject any] struct {
	Id   TObjectId  `json:"id"`
	Item TBizObject `json:"item"`
}

type EditRefObjectRequest[TObjectId ObjectIdentifier, TBizObject any] struct {
	Id   TObjectId  `json:"id"`
	Item TBizObject `json:"item"`
}

type DeleteRefObjectRequest[TObjectId ObjectIdentifier] struct {
	Id     TObjectId `json:"id"`
	ItemId TObjectId `json:"itemId"`
}

type GetRefObjectRequest[TObjectId ObjectIdentifier] struct {
	Id     TObjectId `json:"id"`
	ItemId TObjectId `json:"itemId"`
}

type QueryRefObjectRequest[TObjectId ObjectIdentifier] struct {
	Id      TObjectId           `json:"id"`
	Filters map[string]string   `json:"filters,omitempty"`
	List    *protobuf.ListParam `json:"list,omitempty"`
}

/* response */

type QueryRefObjectResponse[TBizObject any] struct {
	Total int64        `json:"total"`
	Items []TBizObject `json:"items"`
}
