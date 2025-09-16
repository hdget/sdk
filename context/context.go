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
