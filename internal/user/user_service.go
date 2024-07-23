package user

import (
	"HomeWork5/util"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    10 * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, user *UserReq) (*UserRes, error) {
	const op = "user,.CreateUser"

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := User{
		Username: user.Username,
		Password: hashedPassword,
		Email:    user.Email,
	}

	r, err := s.Repository.CreateUser(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UserRes{
		strconv.FormatInt(r.ID, 10),
		r.Username,
		r.Email,
		"user was successfully created",
	}, nil
}

func (s *service) Login(c context.Context, user *UserReq) (*LoginUser, error) {
	const op = "user.Login"

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	dbUser, err := s.Repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	flag := util.CheckPasswordHash(user.Password, dbUser.Password)
	if !flag {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := NewToken(*dbUser)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &LoginUser{
		Token:    token,
		Username: dbUser.Username,
		ID:       dbUser.ID,
	}, nil
}

func NewToken(user User) (string, error) {
	const op = "user.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["uname"] = user.Username
	claims["uemail"] = user.Email
	claims["exp"] = time.Now().Add(24 * time.Hour)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
