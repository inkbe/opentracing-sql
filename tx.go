package sql

import (
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

// conn defines a tracing wrapper for driver.Tx.
type tx struct {
	tx     driver.Tx
	tracer *tracer
	span   opentracing.Span
}

// Commit implements driver.Tx Commit.
func (t *tx) Commit() error {
	if t.span != nil {
		defer t.span.Finish()
	}
	return t.tx.Commit()
}

// Rollback implements driver.Tx Rollback.
func (t *tx) Rollback() error {
	if t.span != nil {
		defer t.span.Finish()
	}
	return t.tx.Rollback()
}
