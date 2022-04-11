package service

/**
This package will contains all GRPC service implementations.
**/

import (
	"agora/assignments/agora"
	"agora/assignments/userpb"
	"context"
)

// UserService implements UserServiceServicer.
type UserService struct {
	Manager agora.UserManager
	userpb.UnimplementedUserServiceServer
}

// GetUserByID retrieves user for given userID.
func (u *UserService) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.User, error) {
	var user, err = u.Manager.GetUserByID(context.Background(), req.Id)
	if err != nil {
		return &userpb.User{}, err
	}

	return &userpb.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}, nil
}

// GetUsers retrives all users.
func (u *UserService) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	var users, err = u.Manager.GetUsers(context.Background(), req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	var result []*userpb.User
	for _, user := range users {
		result = append(result, &userpb.User{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Picture:   user.Picture,
		})
	}

	return &userpb.GetUsersResponse{
		Users: result,
	}, nil
}
