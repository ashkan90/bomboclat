package interfaces

import "context"

type UserService interface {
	GetUser(ctx context.Context, id int) (*User, error)
	GetUserWithCache(ctx context.Context, id int) (*User, error)
	GetUsersBulk(ctx context.Context, ids []int) ([]*User, error)
	UpdateProfile(ctx context.Context, user *User, profile *Profile) error
	NotifyUser(ctx context.Context, user *User, message string) error
}

type ProfileService interface {
	GetProfile(ctx context.Context, userId int) (*Profile, error)
	ValidateProfile(ctx context.Context, profile *Profile) (bool, error)
}

type User struct {
	ID     int
	Name   string
	Status string
}

type Profile struct {
	ID      int
	UserID  int
	Address string
}
