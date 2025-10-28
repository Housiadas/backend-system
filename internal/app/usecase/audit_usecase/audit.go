// Package audit_usecase maintains the app layer api for the audit domain.
package audit_usecase

import (
	"context"

	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
)

type App struct {
	AuditCore *auditcore.Core
}

func NewApp(core *auditcore.Core) *App {
	return &App{
		AuditCore: core,
	}
}

func (a *App) Query(ctx context.Context, qp AppQueryParams) (page.Result[Audit], error) {
	p, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return page.Result[Audit]{}, validation.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Result[Audit]{}, err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, user.DefaultOrderBy)
	if err != nil {
		return page.Result[Audit]{}, validation.NewFieldErrors("order", err)
	}

	adts, err := a.AuditCore.Query(ctx, filter, orderBy, p)
	if err != nil {
		return page.Result[Audit]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.AuditCore.Count(ctx, filter)
	if err != nil {
		return page.Result[Audit]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewResult(toAppAudits(adts), total, p), nil
}
