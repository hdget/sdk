package context

import (
	"context"

	"github.com/hdget/common/meta"
	"google.golang.org/grpc/metadata"
)

var (
	awareMetaKeys = []string{
		meta.KeyTid,
		meta.KeyUid,
		meta.KeyAppId,
		meta.KeyTsn,
	}
)

func New(kvs ...string) context.Context {
	if len(kvs) == 0 || len(kvs)%2 == 1 {
		return context.Background()
	}
	
	md := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		md[kvs[i]] = kvs[i+1]
	}
	return meta.AddMetaToContext(context.Background(), md)
}

func FromIncomingGrpcContext(ctx context.Context) context.Context {
	// 尝试从grpc context中获取meta信息
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return context.Background()
	}

	metas := make(map[string]string)
	for _, key := range awareMetaKeys {
		if values := md.Get(key); len(values) > 0 {
			metas[key] = values[0]
		}
	}
	return meta.AddMetaToContext(context.Background(), metas)
}

func NewOutgoingGrpcContext(ctx context.Context) context.Context {
	metas := meta.GetMetaFromContext(ctx)
	if len(metas) == 0 {
		return context.Background()
	}

	md := make(map[string][]string, len(metas))
	for k, v := range metas {
		md[k] = []string{v}
	}
	return metadata.NewOutgoingContext(ctx, md)
}
