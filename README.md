# OpenTracing-compatible SQL driver wrapper for Golang

This package is an implementation of a SQL driver wrapper with tracing capabilities, compatible with Opentracing API.

## Usage

Register a new database driver by passing an instance created by calling NewTracingDriver:

```
var driver *sql.Driver
var tracer opentracing.Tracer
// init driver, tracer.
...
sql.Register("opentracing-sql", otsql.NewTracingDriver(driver, tracer))
db, err := sql.Open("opentracing-sql", ...)
// use db handle as usual.
```

By default, runtime-based naming function will be used, which will set the span name according to the name of the
function being called (e.g. conn.QueryContext).

It's also possible to specify your own naming function:

```
otsql.NewTracingDriver(driver, tracer, otsql.SpanNameFunction(customNameFunction))
```

Name function format:

```
type SpanNameFunc func(context.Context) string
```

Note that only calls to context-aware DB functions will be traced (e.g. db.QueryContext).

## Comparison with existing packages

There is an existing package https://github.com/ExpansiveWorlds/instrumentedsql which uses the same approach by wrapping
an existing driver with a tracer, however the current implementation provides the following features:
- Pass custom naming function to name spans according to your needs.
- Option to enable/disable logging of SQL queries.

The following features from instrumentedsql package are not supported:
- Passing a custom logger.
- Support of cloud.google.com/go/trace.
- Don't log exact query args.
- Creating spans for LastInsertId, RowsAffected.

## References

[OpenTracing project](http://opentracing.io)
