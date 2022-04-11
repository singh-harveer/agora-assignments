package service

import (
	"agora/assignments/agora"
	"agora/assignments/userpb"
	"context"
	"log"

	"google.golang.org/grpc"
)

// Client represents agora's GRPC Client.
type Client struct {
	grpcConn userpb.UserServiceClient
}

// GetUserByID retrieves user for given userID.
func (c *Client) GetUserByID(ctx context.Context, id string) (agora.User, error) {
	var user, err = c.grpcConn.GetUserByID(ctx, &userpb.GetUserByIDRequest{Id: id})
	if err != nil {
		// TODO- handle error using GRPC error code.
		return agora.User{}, err
	}

	return agora.User{
		ID:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}, nil
}

// GetUsers retrives all users.
func (c *Client) GetUsers(ctx context.Context, page, perPage int64) ([]agora.User, error) {
	var users, err = c.grpcConn.GetUsers(ctx, &userpb.GetUsersRequest{})
	if err != nil {
		// TODO- handle error using GRPC error code.
		return nil, err
	}

	var result []agora.User
	for _, user := range users.Users {
		result = append(result, agora.User{
			ID:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Picture:   user.Picture,
		})
	}

	return result, nil
}

// NewUserClient creates new User grpc service client.
func NewUserClient(ctx context.Context, addr string) (agora.UserManager, error) {
	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &Client{
		grpcConn: userpb.NewUserServiceClient(conn),
	}, nil
}
