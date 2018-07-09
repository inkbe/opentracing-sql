package sql

import (
	"context"
	"database/sql/driver"
)

// conn defines a tracing wrapper for driver.Stmt.
type stmt struct {
	stmt   driver.Stmt
	tracer *tracer
}

// Close implements driver.Stmt Close.
func (s *stmt) Close() error {
	return s.stmt.Close()
}

// NumInput implements driver.Stmt NumInput.
func (s *stmt) NumInput() int {
	return s.stmt.NumInput()
}

// Exec implements driver.Stmt Exec.
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.stmt.Exec(args)
}

// Query implements driver.Stmt Query.
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.stmt.Query(args)
}

// ExecContext implements driver.ExecerContext ExecContext.
func (s *stmt) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	span := s.tracer.newSpan(ctx)
	if s.tracer.saveQuery {
		span.SetTag(TagQuery, query)
	}
	defer span.Finish()
	if execerContext, ok := s.stmt.(driver.ExecerContext); ok {
		return execerContext.ExecContext(ctx, query, args)
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return s.Exec(values)
}

// QueryContext implements Driver.QueryerContext QueryContext.
func (s *stmt) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span := s.tracer.newSpan(ctx)
	if s.tracer.saveQuery {
		span.SetTag(TagQuery, query)
	}
	defer span.Finish()
	if queryerContext, ok := s.stmt.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		return rows, err
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return s.Query(values)
}
