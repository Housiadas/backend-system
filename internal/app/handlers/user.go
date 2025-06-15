package handlers

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/internal/app/service/userapp"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

// User godoc
// @Summary      Crete User
// @Description  Create a new user
// @Tags 		 User
// @Accept       json
// @Produce      json
// @Param        request body userapp.NewUser true "User data"
// @Success      200  {object}  userapp.User
// @Failure      500  {object}  errs.Error
// @Router       /user [post]
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

// User godoc
// @Summary      Update User
// @Description  Update an existing user
// @Tags 		 User
// @Accept       json
// @Produce      json
// @Param        request body userapp.UpdateUser true "User data"
// @Success      200  {object}  userapp.User
// @Failure      500  {object}  errs.Error
// @Router       /user/{user_id} [put]
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

// User godoc
// @Summary      Update User's role
// @Description  Update user's role
// @Tags 		 User
// @Accept       json
// @Produce      json
// @Param        request body userapp.UpdateUserRole true "User data"
// @Success      200  {object}  userapp.User
// @Failure      500  {object}  errs.Error
// @Router       /user/role/{user_id} [put]
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

// User godoc
// @Summary      Delete a user
// @Description  Delete a user
// @Tags 		 User
// @Accept       json
// @Produce      json
// @Success      204
// @Failure      500  {object}  errs.Error
// @Router       /user/{user_id} [delete]
func (h *Handler) userDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	if err := h.App.User.Delete(ctx); err != nil {
		return errs.NewError(err)
	}

	return nil
}

// User godoc
// @Summary      Query Users
// @Description  Search users in database based on criteria
// @Tags		 User
// @Accept       json
// @Produce      json
// @Success      200  {object}  userapp.UserPageResult
// @Failure      500  {object}  errs.Error
// @Router       /user [get]
func (h *Handler) userQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	qp := userParseQueryParams(r)

	usr, err := h.App.User.Query(ctx, qp)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

// User godoc
// @Summary      Find User by id
// @Description  Search user in database by id
// @Tags		 User
// @Accept       json
// @Produce      json
// @Success      200  {object}  userapp.User
// @Failure      500  {object}  errs.Error
// @Router       /user/{user_id} [get]
func (h *Handler) userQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	usr, err := h.App.User.QueryByID(ctx)
	if err != nil {
		return errs.NewError(err)
	}

	return usr
}

func userParseQueryParams(r *http.Request) userapp.AppQueryParams {
	values := r.URL.Query()

	return userapp.AppQueryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("rows"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("user_id"),
		Name:             values.Get("name"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}
}
