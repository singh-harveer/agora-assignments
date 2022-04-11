package gql

import (
	"agora/assignments/agora"
	"context"
)

type graphQL struct {
	client *Client
}

var (
	_ agora.UserManager = (*graphQL)(nil)
)

// GetUserByID retrieves user for given userID.
func (g *graphQL) GetUserByID(ctx context.Context, id string) (agora.User, error) {
	return g.client.GetUserByID(ctx, id)
}

// GetUsers retrives all users.
func (g *graphQL) GetUsers(ctx context.Context, page, perPage int64) ([]agora.User, error) {
	return g.client.GetUsers(ctx, page, perPage)
}

func NewGraphQLClientUserManager(url string, appID string) (agora.UserManager, error) {
	var client, err = newClient(url, appID)
	if err != nil {
		return nil, err
	}

	return &graphQL{
		client: client,
	}, nil
}
