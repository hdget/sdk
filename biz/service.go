package biz

import (
	"context"

	"github.com/hdget/common/meta"
	sdkContext "github.com/hdget/sdk/context"
)

type Service interface {
	Context() context.Context
	GetTid() int64
	GetUid() int64
	GetAppId() string
	GetTsn() string
	GetClient() string
}

type bizSvcImpl struct {
	ctx context.Context
}

func NewService(ctx context.Context) Service {
	return &bizSvcImpl{ctx: sdkContext.FromIncomingGrpcContext(ctx)}
}

func (s bizSvcImpl) Context() context.Context {
	return s.ctx
}

func (s bizSvcImpl) GetTid() int64 {
	return meta.FromServiceContext(s.ctx).GetTid()
}

func (s bizSvcImpl) GetUid() int64 {
	return meta.FromServiceContext(s.ctx).GetUid()
}

func (s bizSvcImpl) GetAppId() string {
	return meta.FromServiceContext(s.ctx).GetAppId()
}

func (s bizSvcImpl) GetTsn() string {
	return meta.FromServiceContext(s.ctx).GetTsn()
}

func (s bizSvcImpl) GetClient() string {
	return meta.FromServiceContext(s.ctx).GetClient()
}
