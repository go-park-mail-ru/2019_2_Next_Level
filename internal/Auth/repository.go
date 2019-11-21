package Auth

type Repository interface {
	GetLoginBySession(uuid string) (login string, err error)
	AddNewSession(login string, uuid string) error
	DeleteSession(uuiв string) error
	DeleteUserSessions(login string) error
	GetUserCredentials(login string) (pass, sault []byte, err error)
	UpdateUserPassword(login string, newPassword []byte, sault []byte) error
}
