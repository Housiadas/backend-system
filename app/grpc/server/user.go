package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/protogen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateUser(ctx context.Context, req *protogen.CreateUserRequest) (*protogen.CreateUserResponse, error) {
	usr, err := s.App.User.Create(ctx, toAppNewUser(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %s", err)
	}

	protogenUser, err := toProtogenUser(usr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %s", err)
	}

	res := &protogen.CreateUserResponse{
		User: protogenUser,
	}

	return res, nil
}

func toAppNewUser(req *protogen.CreateUserRequest) userapp.NewUser {
	return userapp.NewUser{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Roles:           req.GetRoles(),
		Department:      req.GetDepartment(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func toProtogenUser(user userapp.User) (*protogen.User, error) {
	dateCreated, err := time.Parse(time.RFC3339, user.DateCreated)
	if err != nil {
		return &protogen.User{}, err
	}

	dateUpdated, err := time.Parse(time.RFC3339, user.DateUpdated)
	if err != nil {
		return &protogen.User{}, err
	}

	return &protogen.User{
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
