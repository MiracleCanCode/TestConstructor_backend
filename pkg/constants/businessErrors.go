package constants

var (
	NotFoundUser             = "Пользователь не найден"
	LoginOrPasswordIncorrect = "Неправильный логин или пароль"
	LoginIsExist             = "Пользователь с таким логином уже существует"
	EmailExist               = "Пользователь с такой почтой уже существует"
	ErrorUpdateUserData      = "Не получилось обновить данные, попробуйте позже"
	ErrRegistration          = "Не получилось зарегистрировать вас в системе, попробуйте позже"
	ErrGetUserData           = "Не удалось войти в систему, попробуйте позже"
	ErrLogout                = "Не получилось выйти из аккаунта, попробуйте позже"
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
	ErrTestValidation     = "Ошибка. проверки результата теста. Попробуйте в другой раз"
)
