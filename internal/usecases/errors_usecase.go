package usecases

import "errors"

var (
	ErrInvalidCredentials = errors.New("Невалидные данные")
	ErrLoginAlreadyTaken  = errors.New("Пользователь с таким логином уже существует")
	ErrUserNotFound       = errors.New("Пользователь не найден")
	ErrInternalServer     = errors.New("Ошибка сервера, попробуйте в другой раз")
)

var (
	ErrInvalidID         = errors.New("invalid ID format")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrTestNotFound      = errors.New("test not found")
	ErrUserNotAuthorized = errors.New("user is not authorized to perform this action")
)
