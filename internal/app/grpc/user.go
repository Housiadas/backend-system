package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"github.com/Housiadas/backend-system/internal/app/service/userapp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetUserById(
	ctx context.Context,
	req *userV1.GetUserByIdRequest,
) (*userV1.GetUserByIdResponse, error) {

	appUsr, err := s.App.User.Query(ctx, toUserQueryParams(req))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if len(appUsr.Data) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	protoUsr, err := toProtoUser(appUsr.Data[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %s", err)
	}

	return &userV1.GetUserByIdResponse{
		User: protoUsr,
	}, nil
}

func toUserQueryParams(req *userV1.GetUserByIdRequest) userapp.QueryParams {
	return userapp.QueryParams{
		ID: req.Id,
	}
}

func toProtoUser(user userapp.User) (*userV1.User, error) {
	dateCreated, err := time.Parse(time.RFC3339, user.DateCreated)
	if err != nil {
		return &userV1.User{}, err
	}

	dateUpdated, err := time.Parse(time.RFC3339, user.DateUpdated)
	if err != nil {
		return &userV1.User{}, err
	}

	return &userV1.User{
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
