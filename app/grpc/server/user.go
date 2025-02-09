package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	userv1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetUserById(
	ctx context.Context,
	req *userv1.GetUserRequest,
) (*userv1.GetUserResponse, error) {

	appUsr, err := s.App.User.Query(ctx, toUserQueryParams(req))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	protoUsr, err := toProtoUser(appUsr.Data[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %s", err)
	}

	return &userv1.GetUserResponse{
		User: protoUsr,
	}, nil
}

func toUserQueryParams(req *userv1.GetUserRequest) userapp.QueryParams {
	return userapp.QueryParams{
		ID: req.Id,
	}
}

func toProtoUser(user userapp.User) (*userv1.User, error) {
	dateCreated, err := time.Parse(time.RFC3339, user.DateCreated)
	if err != nil {
		return &userv1.User{}, err
	}

	dateUpdated, err := time.Parse(time.RFC3339, user.DateUpdated)
	if err != nil {
		return &userv1.User{}, err
	}

	return &userv1.User{
		Id:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		Department:   user.Department,
		Enabled:      user.Enabled,
		DateCreated:  timestamppb.New(dateCreated),
		DateUpdated:  timestamppb.New(dateUpdated),
	}, nil
}
