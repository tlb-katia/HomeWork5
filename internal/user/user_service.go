package user

import (
	"HomeWork5/util"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
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
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	u := User{
		Username: user.Username,
		Password: hashedPassword,
		Email:    user.Email,
	}

	r, err := s.Repository.CreateUser(ctx, &u)
	if err != nil {
		return nil, err
	}

	return &UserRes{
		strconv.FormatInt(r.ID, 10),
		r.Username,
		r.Email,
	}, nil
}

func (s *service) Login(c context.Context, user *UserReq) (*LoginUser, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	dbUser, err := s.Repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	flag := util.CheckPasswordHash(user.Password, dbUser.Password)
	if !flag {
		return nil, errors.New("invalid password")
	}

	token, err := NewToken(*dbUser)
	if err != nil {
		return nil, err
	}

	return &LoginUser{
		Token:    token,
		Username: dbUser.Username,
		ID:       dbUser.ID,
	}, nil
}

// TODO: what is JWT

func NewToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["uname"] = user.Username
	claims["uemail"] = user.Email
	claims["exp"] = time.Now().Add(24 * time.Hour)

	//TODO change secret place

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
