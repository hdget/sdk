package types

import "github.com/hdget/sdk/common/protobuf"

type ObjectKeyType interface {
	int64 | int32 | int | string
}

/* request */

type OperateObjectRequest[TObjectKey ObjectKeyType] struct {
	Key TObjectKey `json:"key"`
}

type BulkOperateObjectRequest[TObjectKey ObjectKeyType] struct {
	Keys []TObjectKey `json:"keys"`
}

type QueryObjectRequest struct {
	Filters map[string]string   `json:"filters,omitempty"`
	List    *protobuf.ListParam `json:"list,omitempty"`
}

/* response */

type CreateObjectResponse[TObjectId ObjectKeyType] struct {
	Id TObjectId `json:"id"`
}

type QueryObjectResponse[TBizObject any] struct {
	Total int64        `json:"total"`
	Items []TBizObject `json:"items"`
}
