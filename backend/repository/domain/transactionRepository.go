package dr

import "context"

type TransactionRepository interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
