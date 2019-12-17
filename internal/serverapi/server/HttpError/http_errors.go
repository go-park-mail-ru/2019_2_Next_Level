package httperror

const (
	OK                  = "ok"
	UnknownError        = 1
	UnknownMethod       = 2
	BadParam            = 3
	BadSession          = 4
	AccessForbidden     = 5
	LoginNotExist       = 10
	WrongPassword       = 11
	IncorrectName       = 12
	IncorrectFamilyname = 13
	IncorrectNickname   = 14
	IncorrectBirthDate  = 15
	IncorrectSex        = 16
	IncorrectLogin      = 17
	IncorrectPassword   = 18
	LoginAlreadyExists  = 19
	SameNewPass         = 20
)

var messagesArray map[int]string
var DefaultResponse HttpResponse

func init() {
	DefaultResponse = HttpResponse{Status: OK, isEmpty: true}
	messagesArray = map[int]string{
		1:  "Unknown error",
		2:  "Unknown method",
		3:  "Invalid request",
		4:  "User is not authorized",
		5:  "User is not authorized",
		10: "Such login does not exists",
		11: "Wrong password",
		12: "Invalid first name",
		13: "Invalid second name",
		14: "Invalid nickname",
		15: "Invalid birthDate",
		16: "Invalid sex",
		17: "Invalid login",
		18: "Invalid password",
		19: "User with this login already exists",
		20: "Passwords are same",
	}
}

//easyjson:json
type HttpResponse struct {
	Status string `json:"status"`
	Err    HttpError  `json:"error,omitempty"`
	isEmpty bool
}

func (r HttpResponse) Error() string {
	return r.Status + " " + r.Err.Error()
}

//easyjson:json
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e HttpError) Error() string {
	return e.Message
}

func GetError(err int) *HttpResponse {
	return &HttpResponse{
		Status: "error",
		Err: HttpError{
			Code:    err,
			Message: messagesArray[err],
		}}
}

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
