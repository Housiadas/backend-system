package handler

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

func (h *Handler) userCreate(ctx context.Context, _ http.ResponseWriter, r *http.Request) (web.Encoder, error) {
	var app userapp.NewUser
	if err := web.Decode(r, &app); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	usr, err := h.App.User.Create(ctx, app)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (h *Handler) userUpdate(ctx context.Context, _ http.ResponseWriter, r *http.Request) (web.Encoder, error) {
	var app userapp.UpdateUser
	if err := web.Decode(r, &app); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	usr, err := h.App.User.Update(ctx, app)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (h *Handler) updateRole(ctx context.Context, _ http.ResponseWriter, r *http.Request) (web.Encoder, error) {
	var app userapp.UpdateUserRole
	if err := web.Decode(r, &app); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	usr, err := h.App.User.UpdateRole(ctx, app)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (h *Handler) userDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (web.Encoder, error) {
	if err := h.App.User.Delete(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) userQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) (web.Encoder, error) {
	qp, err := userParseQueryParams(r)
	if err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	usr, err := h.App.User.Query(ctx, qp)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (h *Handler) userQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (web.Encoder, error) {
	usr, err := h.App.User.QueryByID(ctx)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
