package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"late/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

type SQLDB interface {
	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// SQLQuery contains information for the exeqution of the postgresql query.
type SQLQuery struct {
	name string
	sql  string
}

// NewSQLQuery returns new sql query.
func NewSQLQuery(name, sqlQuery string) *SQLQuery {
	return &SQLQuery{name: name, sql: sqlQuery}
}

// String returns name of the query.
// Implements fmt.Stringer interface
func (q *SQLQuery) String() string {
	return fmt.Sprintf("postgres.SQLQuery(%s)", q.name)
}

// GoString returns SQL query text
func (q *SQLQuery) GoString() string { return q.sql }

// QueryRow executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (q *SQLQuery) QueryRow(ctx context.Context, src SQLDB, args ...any) *sql.Row {
	defer measureTime(time.Now(), q.name, "query_row")
	return src.QueryRowContext(ctx, q.sql, args...)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (q *SQLQuery) Query(ctx context.Context, src SQLDB, args ...any) (*sql.Rows, error) {
	defer measureTime(time.Now(), q.name, "query")
	return src.QueryContext(ctx, q.sql, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (q *SQLQuery) Exec(ctx context.Context, src SQLDB, args ...any) (sql.Result, error) {
	defer measureTime(time.Now(), q.name, "exec")
	return src.ExecContext(ctx, q.sql, args...)
}

func measureTime(start time.Time, name, kind string) {
	metrics.DatabaseQueryRequestDuration.With(prometheus.Labels{
		"query": name,
		"kind":  kind,
	}).Observe(time.Since(start).Seconds())
}
