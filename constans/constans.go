package constans

import "errors"

const (
	UserEntity        = `User`
	ProductEntity     = `Product`
	TransactionEntity = `Transaction`

	MessageSuccessReadAll      = "Success retrieve all data from %s"
	MessageSuccessReadByID     = "Success get %s with id %s"
	MessageSuccessCreate       = "Success create new %s"
	MessageSuccessUpdate       = "Success update %s with id %s"
	MessageSuccessDelete       = "Success delete %s with id %s"
	MessageSuccessUploadImage  = "Success upload %s image"
	MessageSuccessCheckoutItem = "Success checkout item"

	DefaultImage  = "image/default.jpg"
	BaseImagePath = "images/%d.%s"

	PENDING = "pending"
)

var (
	ErrInternalServerError  = errors.New("internal server error")
	ErrNotFound             = errors.New("your requested item is not found")
	ErrConflict             = errors.New("data already exist")
	ErrBadParamInput        = errors.New("given param is not valid")
	ErrWrongEmailOrPassword = errors.New("wrong email/password")
)
