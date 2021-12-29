package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

func GetTextMapCarrierFromMetaDAta(ctx context.Context) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key := range md.Copy() {
			metadataMap.Set(key, md.Get(key)[0])
		}
	}
	return metadataMap
}
