package sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"runtime"

	"github.com/opentracing/opentracing-go"
)

// ErrUnsupported is an error returned when the underlying driver doesn't provide a given function.
var ErrUnsupported = errors.New("operation unsupported by the underlying driver")

// spanTagQuery is a span tag for SQL queries.
var spanTagQuery = "query"

// SetTagQuery sets a span tag for SQL queries.
// Use in init funcs only.
func SetTagQuery(tag string) {
	spanTagQuery = tag
}

// SpanNameFunc defines a function which returns a name for the span which is being created on traceable operations.
// Passing span naming function is optional, however it gives the user a way to use a custom naming strategy. To allow
// getting some more information related to the current call, the context, which is passed with the call, is propagated
// to the naming function.
type SpanNameFunc func(context.Context) string

// tracer defines a set of instances for collecting spans.
type tracer struct {
	t            opentracing.Tracer
	nameFunc     SpanNameFunc
	observerFunc func(context.Context, opentracing.Span)
	saveQuery    bool
}

// newSpan creates a new opentracing.Span instance from the given context.
func (t *tracer) newSpan(ctx context.Context) opentracing.Span {
	name := t.nameFunc(ctx)
	var opts []opentracing.StartSpanOption
	parent := opentracing.SpanFromContext(ctx)
	if parent != nil {
		opts = append(opts, opentracing.ChildOf(parent.Context()))
	}
	span := t.t.StartSpan(name, opts...)
	if t.observerFunc != nil {
		t.observerFunc(ctx, span)
	}
	return span
}

// defaultNameFunc defines a default span naming function.
// Call stack at the moment of call to the function has the following frames (digits represent the depth from the top):
// 0 - name function itself.
// 1 - newSpan.
// 2 - wrapper function in this package, e.g. QueryContext.
func defaultNameFunc(ctx context.Context) string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
	return f.Name()
}

// namedValueToValue converts driver arguments of NamedValue format to Value format. Implemented in the same way as in
// database/sql ctxutil.go.
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
