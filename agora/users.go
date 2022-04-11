package agora

/**
This package will contains all interface and entities definations.
**/

import "context"

// User represents User.
type User struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	DateOfBirth  string
	RegisterDate string
	Picture      string `json:"picture"`
	Location     Location
}

type Location struct {
	Street   string
	City     string
	State    string
	Country  string
	TimeZone string
}

// UserManager manages users.
type UserManager interface {
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context, page, limit int64) ([]User, error)
}
