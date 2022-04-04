package handler

import (
	"fmt"
	"net/http"

	"github.com/cecepsprd/ticketing-api/constans"
	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/service"
	"github.com/cecepsprd/ticketing-api/utils"
	"github.com/cecepsprd/ticketing-api/utils/convert"

	"github.com/labstack/echo"
)

type product struct {
	productService service.ProductService
	trxService     service.TransactionService
}

func NewProductHandler(e *echo.Echo, ps service.ProductService, ts service.TransactionService) {
	handler := &product{
		productService: ps,
		trxService:     ts,
	}

	e.POST("/api/products", handler.Create, auth(), isAdmin)
	e.GET("/api/products", handler.Read, auth())
	e.PUT("/api/products/:id", handler.Update, auth(), isAdmin)
	e.DELETE("/api/products/:id", handler.Delete, auth(), isAdmin)
	e.GET("/api/products/:id", handler.ReadByID, auth())
	e.POST("/api/products/checkout", handler.Checkout, auth())
}

func (p *product) Create(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req model.ProductRequest
	)

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	err = p.productService.Create(ctx, req)
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.APIResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf(constans.MessageSuccessCreate, constans.ProductEntity),
		Data:    nil,
	})
}

func (p *product) Read(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	data, err := p.productService.Read(ctx)
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf(constans.MessageSuccessReadAll, constans.ProductEntity),
		Data:    data,
	})
}

func (p *product) Update(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		id  = c.Param("id")
		req model.ProductRequest
	)

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	err = p.productService.Update(ctx, req)
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf(constans.MessageSuccessUpdate, constans.ProductEntity, id),
		Data:    nil,
	})
}

func (p *product) Delete(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		id  = c.Param("id")
	)

	err := p.productService.Delete(ctx, convert.Atoi(id))
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf(constans.MessageSuccessDelete, constans.ProductEntity, id),
		Data:    nil,
	})
}

func (p *product) ReadByID(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		id  = c.Param("id")
	)

	data, err := p.productService.ReadByID(ctx, convert.Atoi(id))
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf(constans.MessageSuccessReadByID, constans.ProductEntity, id),
		Data:    data,
	})
}

func (h *product) Checkout(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req model.CreateTransactionRequest
	)

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	req.User = utils.GetUserByContext(c)

	transaction, err := h.trxService.Checkout(ctx, req)
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf(constans.MessageSuccessCheckoutItem),
		Data:    transaction,
	})
}
