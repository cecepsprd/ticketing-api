package service

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/cecepsprd/ticketing-api/constans"
	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/repository"
	"github.com/cecepsprd/ticketing-api/utils/convert"
	"github.com/cecepsprd/ticketing-api/utils/logger"
	"github.com/spf13/viper"
	"github.com/veritrans/go-midtrans"
)

type TransactionService interface {
	Checkout(ctx context.Context, request model.CreateTransactionRequest) (*model.Transaction, error)
}

type transaction struct {
	transactionRepo repository.TransactionRepository
	productRepo     repository.ProductRepository
	contextTimeout  time.Duration
}

func NewTransactionService(transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, timeout time.Duration) TransactionService {
	return &transaction{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		contextTimeout:  timeout,
	}
}

func (s *transaction) Checkout(ctx context.Context, req model.CreateTransactionRequest) (*model.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	product, err := s.productRepo.ReadByID(ctx, req.ProductID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	if reflect.ValueOf(product).IsNil() {
		return nil, errors.New(constans.ErrNotFound.Error())
	}

	if product.Stock == 0 {
		return nil, errors.New("ticket has run out")
	}

	trx, err := s.transactionRepo.Create(ctx, model.Transaction{
		ProductID: req.ProductID,
		UserID:    req.User.ID,
		Status:    constans.PENDING,
		Amount:    product.Price,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	trx.PaymentURL, err = s.GetPaymentURL(trx, req.User)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	err = s.transactionRepo.Update(ctx, *trx)
	if err != nil {
		logger.Log.Warn(err.Error())
		return nil, err
	}

	return trx, nil
}

func (s *transaction) GetPaymentURL(trx *model.Transaction, user model.User) (paymentURL string, err error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = viper.GetString("MIDTRANS_SERVER_KEY")
	midclient.ClientKey = viper.GetString("MIDTRANS_CLIENT_KEY")
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Username,
			Phone: user.Phone,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(trx.ID)),
			GrossAmt: int64(trx.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

func (s *transaction) Update(ctx context.Context, request model.UpdateTransactionRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	transaction, err := s.transactionRepo.ReadByID(ctx, convert.Atoi(request.OrderID))
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if request.PaymentType == "credit_card" && request.TransactionStatus == "capture" && request.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "deny" || request.TransactionStatus == "expire" || request.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	err = s.transactionRepo.UpdateStatus(ctx, convert.Atoi(request.OrderID), transaction.Status)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if transaction.Status == "paid" {
		err = s.productRepo.UpdateStock(ctx, transaction.ProductID, -1)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	return nil
}
