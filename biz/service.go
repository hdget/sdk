package biz

import (
	"context"

	"github.com/hdget/sdk/dapr/meta"
)

type Service interface {
	GetTid() int64
	GetUid() int64
	Context() context.Context
}

type bizSvcImpl struct {
	ctx context.Context
}

func NewService(ctx context.Context) Service {
	return &bizSvcImpl{ctx: ctx}
}

func (s bizSvcImpl) Context() context.Context {
	return s.ctx
}

func (s bizSvcImpl) GetTid() int64 {
	return meta.New().GetTid(s.ctx)
}

func (s bizSvcImpl) GetUid() int64 {
	return meta.New().GetUid(s.ctx)
}
