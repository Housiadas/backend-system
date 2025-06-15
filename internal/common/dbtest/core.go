package dbtest

import (
	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/internal/app/repository/auditrepo"
	"github.com/Housiadas/backend-system/internal/app/repository/productrepo"
	"github.com/Housiadas/backend-system/internal/app/repository/userrepo"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/logger"
)

// Core represents all the internal core apis needed for testing.
type Core struct {
	Audit   *auditcore.Core
	User    *usercore.Core
	Product *productcore.Core
}

func newCore(log *logger.Logger, db *sqlx.DB) Core {
	auditCore := auditcore.NewCore(log, auditrepo.NewStore(log, db))
	userBus := usercore.NewCore(log, userrepo.NewStore(log, db))
	productBus := productcore.NewCore(log, userBus, productrepo.NewStore(log, db))

	return Core{
		Audit:   auditCore,
		User:    userBus,
		Product: productBus,
	}
}
