package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AcquireAdvisoryLock acquires a session-level PostgreSQL advisory lock.
// The caller must call conn.Release() when done to return the connection to the pool,
// which also releases the lock automatically.
// Use this for distributed workers that must not run concurrently (e.g. OutboxWorker).
func AcquireAdvisoryLock(ctx context.Context, pool *pgxpool.Pool, lockID int64) (*pgxpool.Conn, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockID)
	if err != nil {
		conn.Release()
		return nil, err
	}

	return conn, nil
}

// TryAdvisoryLock attempts to acquire a PostgreSQL advisory lock without blocking.
// Returns (true, conn) if the lock was acquired, (false, nil) if already held by another session.
// The caller must call conn.Release() when done.
func TryAdvisoryLock(ctx context.Context, pool *pgxpool.Pool, lockID int64) (bool, *pgxpool.Conn, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return false, nil, err
	}

	var acquired bool
	err = conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, lockID).Scan(&acquired)
	if err != nil {
		conn.Release()
		return false, nil, err
	}

	if !acquired {
		conn.Release()
		return false, nil, nil
	}

	return true, conn, nil
}

// ReleaseAdvisoryLock explicitly releases a session-level advisory lock and returns
// the connection to the pool.
func ReleaseAdvisoryLock(ctx context.Context, conn *pgxpool.Conn, lockID int64) error {
	_, err := conn.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockID)
	conn.Release()
	return err
}
