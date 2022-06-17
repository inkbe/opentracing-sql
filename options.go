package sql

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// SpanNameFunction is an option for using a custom span naming function.
func SpanNameFunction(f SpanNameFunc) func(*tracingDriver) {
	return func(d *tracingDriver) {
		d.tracer.nameFunc = f
	}
}

// SaveQuery is an option for saving SQL queries.
func SaveQuery(f SpanNameFunc) func(*tracingDriver) {
	return func(d *tracingDriver) {
		d.tracer.saveQuery = true
	}
}

// WithSpanObserver allows you to modify the span's tags for every span created.
func WithSpanObserver(obsFunc func(context.Context, opentracing.Span)) func(*tracingDriver) {
	return func(d *tracingDriver) {
		d.tracer.observerFunc = obsFunc
	}
}
