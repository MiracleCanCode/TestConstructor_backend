package constants

var (
	NotFoundUser             = "Пользователь не найден"
	LoginOrPasswordIncorrect = "Неправильный логин или пароль"
	LoginIsExist             = "Пользователь с таким логином уже существует"
	EmailExist               = "Пользователь с такой почтой уже существует"
	ErrorUpdateUserData      = "Не получилось обновить данные, попробуйте позже"
)

var (
	InternalServerError = "Ошибка сервера, попробуйте в другой раз"
)

var (
	GetTestByIdError      = "Ошибка, получения теста по id"
	ErrorChangeActiveTest = "Ошибка, изменения видимости теста"
	ErrorDeleteTest       = "Ошибка, удаления теста"
	ErrorCreateTest       = "Ошибка, создания теста"
	ErrorGetAllTests      = "Ошибка, получения тестов"
)
