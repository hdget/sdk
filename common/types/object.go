package types

import "github.com/hdget/sdk/common/protobuf"

type ObjectIdentifier interface {
	int64 | int32 | int | string
}

/* request */

type OperateObjectRequest[TObjectId ObjectIdentifier] struct {
	Id TObjectId `json:"id"`
}

type BulkOperateObjectRequest[TObjectId ObjectIdentifier] struct {
	Ids []TObjectId `json:"ids"`
}

type QueryObjectRequest struct {
	Filters map[string]string   `json:"filters,omitempty"`
	List    *protobuf.ListParam `json:"list,omitempty"`
}

/* response */

type CreateObjectResponse[TObjectId ObjectIdentifier] struct {
	Id TObjectId `json:"id"`
}

type QueryObjectResponse[TBizObject any] struct {
	Total int64        `json:"total"`
	Items []TBizObject `json:"items"`
}
