package user

import (
	"context"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Message  string `json:"message"`
}

type LoginUser struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUser(ctx context.Context, user *UserReq) (*UserRes, error)
	Login(ctx context.Context, user *UserReq) (*LoginUser, error)
}
