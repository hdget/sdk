package types

import (
	"github.com/hdget/sdk/common/protobuf"
)

/* request */

type CreateRefObjectRequest[TObjectKey ObjectKeyType, TBizObject any] struct {
	Key  TObjectKey `json:"key"`
	Item TBizObject `json:"item"`
}

type EditRefObjectRequest[TObjectKey ObjectKeyType, TBizObject any] struct {
	Key  TObjectKey `json:"key"`
	Item TBizObject `json:"item"`
}

type DeleteRefObjectRequest[TObjectKey ObjectKeyType] struct {
	Key     TObjectKey `json:"key"`
	ItemKey TObjectKey `json:"itemKey"`
}

type GetRefObjectRequest[TObjectKey ObjectKeyType] struct {
	Key     TObjectKey `json:"key"`
	ItemKey TObjectKey `json:"itemKey"`
}

type QueryRefObjectRequest[TObjectKey ObjectKeyType] struct {
	Key     TObjectKey          `json:"key"`
	Filters map[string]string   `json:"filters,omitempty"`
	List    *protobuf.ListParam `json:"list,omitempty"`
}

/* response */

type QueryRefObjectResponse[TBizObject any] struct {
	Total int64        `json:"total"`
	Items []TBizObject `json:"items"`
}
