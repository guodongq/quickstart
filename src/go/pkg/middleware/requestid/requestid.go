package requestid

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

const (
	metadataKey = "x-request-id"
)

var requestIDKeys = []string{
	"Request-ID", "X-Request-ID",
}

// Get reads the Request-ID and X-Request-ID HTTP header from an `*http.Request`
// If no header is set, an empty string is returned
func Get(r *http.Request) string {
	for _, try := range requestIDKeys {
		if id := r.Header.Get(try); id != "" {
			return id
		}
	}
	return ""
}

// FromContext returns a request ID from gRPC metadata if available in ctx.
func FromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	ids, ok := md[metadataKey]
	if !ok || len(ids) == 0 {
		return "", false
	}

	return ids[0], true
}

// NewMetadata constructs gRPC metadata with the request ID set.
func NewMetadata(id string) metadata.MD {
	return metadata.Pairs(metadataKey, id)
}

// AppendToOutgoingContext returns a context with the request-id added to the gRPC metadata.
func AppendToOutgoingContext(ctx context.Context) context.Context {
	id, ok := FromContext(ctx)
	if !ok {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, metadataKey, id)
}
