package sql

import (
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

// conn defines a tracing wrapper for driver.Driver.
type tracingDriver struct {
	driver driver.Driver
	tracer *tracer
}

// TracingDriver creates and returns a new SQL driver with tracing capabilities.
func NewTracingDriver(d driver.Driver, t opentracing.Tracer, options ...func(*tracingDriver)) driver.Driver {
	td := &tracingDriver{driver: d, tracer: &tracer{t: t}}
	for _, option := range options {
		option(td)
	}
	if td.tracer.nameFunc == nil {
		td.tracer.nameFunc = defaultNameFunc
	}
	return td
}

// Open implements driver.Driver Open.
func (d *tracingDriver) Open(name string) (driver.Conn, error) {
	c, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &conn{conn: c, tracer: d.tracer}, nil
}
