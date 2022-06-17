package sql

import (
	"context"
	"database/sql/driver"
)

// conn defines a tracing wrapper for driver.Conn.
type conn struct {
	conn   driver.Conn
	tracer *tracer
}

// Prepare implements driver.Conn Prepare.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	s, err := c.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &stmt{stmt: s, tracer: c.tracer}, nil
}

// Prepare implements driver.Conn Close.
func (c *conn) Close() error {
	return c.conn.Close()
}

// Prepare implements driver.Conn Begin.
func (c *conn) Begin() (driver.Tx, error) {
	t, err := c.conn.Begin()
	if err != nil {
		return nil, err
	}
	return &tx{tx: t, tracer: c.tracer}, nil
}

// BeginTx implements driver.ConnBeginTx BeginTx.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	s := c.tracer.newSpan(ctx)
	if connBeginTx, ok := c.conn.(driver.ConnBeginTx); ok {
		t, err := connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return &tx{tx: t, tracer: c.tracer, span: s}, nil
	}
	return c.conn.Begin()
}

// PrepareContext implements driver.ConnPrepareContext PrepareContext.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connPrepareContext, ok := c.conn.(driver.ConnPrepareContext); ok {
		s, err := connPrepareContext.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &stmt{stmt: s, tracer: c.tracer}, nil
	}
	return c.conn.Prepare(query)
}

// Exec implements driver.Execer Exec.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if execer, ok := c.conn.(driver.Execer); ok {
		return execer.Exec(query, args)
	}
	return nil, ErrUnsupported
}

// Exec implements driver.StmtExecContext ExecContext.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	s := c.tracer.newSpan(ctx)
	if c.tracer.saveQuery {
		s.SetTag(spanTagQuery, query)
	}
	defer s.Finish()
	if execerContext, ok := c.conn.(driver.ExecerContext); ok {
		r, err := execerContext.ExecContext(ctx, query, args)
		return r, err
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return c.Exec(query, values)
}

// Ping implements driver.Pinger Ping.
func (c *conn) Ping(ctx context.Context) error {
	if pinger, ok := c.conn.(driver.Pinger); ok {
		s := c.tracer.newSpan(ctx)
		defer s.Finish()
		return pinger.Ping(ctx)
	}
	return ErrUnsupported
}

// Query implements driver.Queryer Query.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if queryer, ok := c.conn.(driver.Queryer); ok {
		return queryer.Query(query, args)
	}
	return nil, ErrUnsupported
}

// QueryContext implements Driver.QueryerContext QueryContext.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	s := c.tracer.newSpan(ctx)
	if c.tracer.saveQuery {
		s.SetTag(spanTagQuery, query)
	}
	defer s.Finish()
	if queryerContext, ok := c.conn.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		return rows, err
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return c.Query(query, values)
}
