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
	md := toMap(kvs...)
	if md == nil {
		return context.Background()
	}
	return meta.AddMetaToContext(context.Background(), md)
}

func FromIncomingGrpcContext(ctx context.Context) context.Context {
	// 尝试从grpc context中获取meta信息
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return context.Background()
	}

	metaMap := make(map[string]string)
	for _, key := range awareMetaKeys {
		if values := md.Get(key); len(values) > 0 {
			metaMap[key] = values[0]
		}
	}
	return meta.AddMetaToContext(context.Background(), metaMap)
}

func NewOutgoingGrpcContext(ctx context.Context, kvs ...string) context.Context {
	metaMap := meta.GetMetaFromContext(ctx)

	md := make(map[string][]string, len(metaMap)+len(kvs)/2)
	for k, v := range metaMap {
		md[k] = []string{v}
	}

	// last one override first
	kvMap := toMap(kvs...)
	for k, v := range kvMap {
		md[k] = []string{v}
	}
	return metadata.NewOutgoingContext(context.Background(), md)
}

func toMap(kvs ...string) map[string]string {
	if len(kvs) == 0 || len(kvs)%2 == 1 {
		return nil
	}

	result := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		result[kvs[i]] = kvs[i+1]
	}
	return result
}
