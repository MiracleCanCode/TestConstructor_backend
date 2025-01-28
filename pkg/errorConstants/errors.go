package errorconstant

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid data")
	ErrLoginAlreadyTaken  = errors.New("user login is exist")
	ErrUserNotFound       = errors.New("user not found")
	ErrInternalServer     = errors.New("Ошибка сервера, попробуйте в другой раз")
)

var (
	ErrInvalidID         = errors.New("invalid ID format")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrTestNotFound      = errors.New("test not found")
	ErrUserNotAuthorized = errors.New("user is not authorized to perform this action")
	ErrRegisterUser      = errors.New("register user error")
)
