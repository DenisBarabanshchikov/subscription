package port

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

type Customer interface {
	CreateCustomer(ctx context.Context, email string) model.Customer
}
