package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cecepsprd/ticketing-api/constans"
	"github.com/cecepsprd/ticketing-api/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByContext(c echo.Context) model.User {
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return model.User{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		Phone:    claims["phone"].(string),
	}
}

func HashPassword(password string) (hashedPassword string, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func MappingRequest(request interface{}, model interface{}) error {
	// convert interface to json
	jsonRecords, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Error encode records: %s", err.Error())
	}

	// bind json to struct
	if err := json.Unmarshal(jsonRecords, model); err != nil {
		return fmt.Errorf("Error decode json to struct: %s", err.Error())
	}

	return nil
}

func SetHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case constans.ErrInternalServerError:
		return http.StatusInternalServerError
	case constans.ErrNotFound:
		return http.StatusNotFound
	case constans.ErrConflict:
		return http.StatusConflict
	case constans.ErrWrongEmailOrPassword:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
