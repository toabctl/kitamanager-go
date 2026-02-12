package store

import (
	"context"

	"gorm.io/gorm"
)

// txKey is the context key for storing a transaction.
type txKey struct{}

// Transactor provides transaction management for service-layer operations
// that span multiple store calls.
type Transactor interface {
	// InTransaction executes fn within a database transaction.
	// If fn returns an error, the transaction is rolled back.
	// Nested calls to InTransaction reuse the outer transaction.
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// GormTransactor implements Transactor using GORM.
type GormTransactor struct {
	db *gorm.DB
}

// NewTransactor creates a new GormTransactor.
func NewTransactor(db *gorm.DB) *GormTransactor {
	return &GormTransactor{db: db}
}

// InTransaction executes fn within a database transaction.
// The transaction is stored in the context so that stores can pick it up
// via DBFromContext.
func (t *GormTransactor) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// If already inside a transaction, reuse it (nested call).
	if _, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return fn(ctx)
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// DBFromContext returns the transaction DB from context if present,
// otherwise returns the provided default DB. Stores should call this
// at the start of each operation to participate in service-layer transactions.
func DBFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
