package httperror

const (
	OK                  = "OK"
	UnknownError        = "1"
	UnknownMethod       = "2"
	BadParam            = "3"
	BadSession          = "4"
	AccessForbidden     = "5"
	LoginNotExist       = "10"
	WrongPassword       = "11"
	IncorrectName       = "12"
	IncorrectFamilyname = "13"
	IncorrectNickname   = "14"
	IncorrectBirthDate  = "15"
	IncorrectSex        = "16"
	IncorrectLogin      = "17"
	IncorrectPassword   = "18"
	LoginAlreadyExists  = "19"
	SameNewPass         = "20"
)

// func ConvertFromError(errI error) string {
// 	if errI == nil {
// 		return OK
// 	}
// 	err, _ := errI.(e.Error)
// 	switch err.Code {
// 	case e.IncorrectParams:
// 		return
// 	}
// }

// 1	Произошла неизвестная ошибка.
// 2	Передан неизвестный метод. Проверьте, правильно ли указано название вызываемого метода.
// 3	Один из необходимых параметров был не передан или неверен. Проверьте список требуемых параметров и их формат на странице с описанием метода
// 4	Авторизация пользователя не удалась. session_id недействителен.
// 5	Доступ запрещён.
// 10	Пользователя с таким логином не существует.
// 11	Неверный пароль.
// 12	Некорректное имя.
// 13	Некорректная фамилия.
// 14	Некорректный никнейм.
// 15	Некорректная дата рождения.
// 16	Некорректный пол.
// 17	Некорректный логин.
// 18	Некорректный пароль.
// 19	Пользователь с таким логином уже существует.
// 20	Новый пароль совпадает со старым.
