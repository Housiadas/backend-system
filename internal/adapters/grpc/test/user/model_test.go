package user_test

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"github.com/Housiadas/backend-system/internal/adapters/domain/userapp"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
)

func toUserQueryParams(req *userV1.GetUserByIdRequest) userapp.QueryParams {
	return userapp.QueryParams{
		ID: req.Id,
	}
}

func toProtoUser(user userservice.User) *userV1.User {
	dateCreated, err := time.Parse(time.RFC3339, user.DateCreated.String())
	if err != nil {
		return &userV1.User{}
	}

	dateUpdated, err := time.Parse(time.RFC3339, user.DateUpdated.String())
	if err != nil {
		return &userV1.User{}
	}

	return &userV1.User{
		Id:           user.ID.String(),
		Name:         user.Name.String(),
		Email:        user.Email.String(),
		Roles:        role.ParseToString(user.Roles),
		PasswordHash: user.PasswordHash,
		Department:   user.Department.String(),
		Enabled:      user.Enabled,
		DateCreated:  timestamppb.New(dateCreated),
		DateUpdated:  timestamppb.New(dateUpdated),
	}
}
