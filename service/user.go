package service

import (
	"context"
	"reflect"
	"time"

	"github.com/cecepsprd/ticketing-api/constans"
	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/repository"
	"github.com/cecepsprd/ticketing-api/utils"
	"github.com/cecepsprd/ticketing-api/utils/logger"
)

type UserService interface {
	Create(context.Context, model.User) error
	Read(ctx context.Context) (users []model.User, err error)
	ReadByUsername(ctx context.Context, username string) (*model.User, error)
}

type user struct {
	repo           repository.UserRepository
	contextTimeout time.Duration
}

func NewUserService(urepo repository.UserRepository, timeout time.Duration) UserService {
	return &user{
		repo:           urepo,
		contextTimeout: timeout,
	}
}

func (s *user) Read(ctx context.Context) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	users, err := s.repo.Read(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (s *user) Create(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	userCounted, err := s.repo.CountUser(ctx, user)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if !reflect.ValueOf(userCounted).IsZero() {
		return constans.ErrConflict
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *user) ReadByUsername(ctx context.Context, username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.ReadByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
