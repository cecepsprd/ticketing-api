package service

import (
	"context"
	"time"

	"github.com/cecepsprd/ticketing-api/constans"
	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/repository"
	"github.com/cecepsprd/ticketing-api/utils"
	"github.com/cecepsprd/ticketing-api/utils/logger"
)

type ProductService interface {
	Read(context.Context) ([]model.Product, error)
	Create(ctx context.Context, product model.ProductRequest) error
	Update(ctx context.Context, product model.ProductRequest) error
	Delete(ctx context.Context, id int64) error
	ReadByID(ctx context.Context, productID int64) (*model.Product, error)
}

type product struct {
	repo           repository.ProductRepository
	contextTimeout time.Duration
}

func NewProductService(repo repository.ProductRepository, timeout time.Duration) ProductService {
	return &product{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (s *product) Read(ctx context.Context) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	products, err := s.repo.Read(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return products, nil
}

func (s *product) Create(ctx context.Context, request model.ProductRequest) error {
	var (
		err     error
		product model.Product
	)

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := utils.MappingRequest(request, &product); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if request.ImageURL == "" {
		product.ImageURL = constans.DefaultImage
	}

	err = s.repo.Create(ctx, product)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *product) Update(ctx context.Context, request model.ProductRequest) error {
	var (
		product model.Product
	)

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := utils.MappingRequest(request, &product)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	err = s.repo.Update(ctx, product)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *product) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *product) ReadByID(ctx context.Context, productID int64) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	product, err := s.repo.ReadByID(ctx, productID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return product, nil
}
