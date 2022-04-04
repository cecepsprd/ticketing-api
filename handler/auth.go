package handler

import (
	"net/http"

	"github.com/cecepsprd/ticketing-api/model"
	"github.com/cecepsprd/ticketing-api/service"
	"github.com/cecepsprd/ticketing-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(e *echo.Echo, as service.AuthService) {
	handler := &AuthHandler{
		authService: as,
	}

	e.POST("/api/auth/signin", handler.Login)
}

func (ah *AuthHandler) Login(c echo.Context) error {
	var (
		req = model.LoginRequest{}
		ctx = c.Request().Context()
	)

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	response, err := ah.authService.Login(ctx, req)
	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, response)
}

func auth() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString("APP_JWT_SECRET")),
	})
}

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		roles := claims["roles"].(string)
		if roles != "admin" {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
