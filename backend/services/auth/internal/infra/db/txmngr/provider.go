package txmngr

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/db"
)

func NewProvider(dbpool db.PgxPool) func() (*manager.Manager, *trmpgx.CtxGetter) {
	return func() (*manager.Manager, *trmpgx.CtxGetter) {
		return New(dbpool)
	}
}
