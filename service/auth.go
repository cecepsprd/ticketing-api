package service

import (
	"context"
	"errors"
	"time"

	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/utils/logger"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userService UserService
	JWTSecret   string
}

type AuthService interface {
	Login(context.Context, model.LoginRequest) (*model.LoginResponse, error)
}

func NewAuthService(us UserService, JWTSecret string) AuthService {
	return &authService{
		userService: us,
		JWTSecret:   JWTSecret,
	}
}

func (s *authService) Login(ctx context.Context, request model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userService.ReadByUsername(ctx, request.Username)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New("password incorrect")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["phone"] = user.Phone
	claims["roles"] = user.Roles
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: t,
	}, nil
}
