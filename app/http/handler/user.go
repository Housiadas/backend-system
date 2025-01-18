package handler

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

func (h *Handler) userCreate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app userapp.NewUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := h.App.User.Create(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

func (h *Handler) userUpdate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app userapp.UpdateUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := h.App.User.Update(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

func (h *Handler) updateRole(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app userapp.UpdateUserRole
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := h.App.User.UpdateRole(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

func (h *Handler) userDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	if err := h.App.User.Delete(ctx); err != nil {
		return errs.NewError(err)
	}

	return nil
}

func (h *Handler) userQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	qp := userapp.ParseQueryParams(r)

	usr, err := h.App.User.Query(ctx, qp)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

func (h *Handler) userQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	usr, err := h.App.User.QueryByID(ctx)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}
