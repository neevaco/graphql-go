package trace

import (
	"context"

	"github.com/graph-gophers/graphql-go/errors"
)

// ValidationStats represents metrics collected during validation.
type ValidationStats struct {
	CacheHit bool // Was this a cache hit?
	CacheLen int  // Number of valid queries in the cache.
}

type TraceValidationFinishFunc = func([]*errors.QueryError, ValidationStats)

type ValidationTracer interface {
	TraceValidation(ctx context.Context) (context.Context, TraceValidationFinishFunc)
}

type NoopValidationTracer struct{}

func (NoopValidationTracer) TraceValidation(ctx context.Context) (context.Context, TraceValidationFinishFunc) {
	return ctx, func([]*errors.QueryError, ValidationStats) {}
}
