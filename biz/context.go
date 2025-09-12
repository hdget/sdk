package biz

import (
	"context"

	"github.com/hdget/common/constant"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

type Contexter interface {
	Context() context.Context
	Set(key string, val any) Contexter
}

type ctxImpl struct {
	ctx context.Context
}

func NewTenantContext(tid int64) context.Context {
	c := NewContext()
	c.Set(constant.MetaKeyTid, tid)
	return c.Context()
}

func NewContext() Contexter {
	return &ctxImpl{
		ctx: context.Background(),
	}
}

func (c *ctxImpl) Set(key string, val any) Contexter {
	md := metadata.Pairs(key, cast.ToString(val))
	c.ctx = metadata.NewOutgoingContext(c.ctx, md)
	return c
}

func (c *ctxImpl) Context() context.Context {
	return c.ctx
}
