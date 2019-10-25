package error

import "fmt"

const (
	InvalidParams = iota
	AlreadyExists
	NotExists
	OK
	ProcessError
	SessionExists
	UserExists
	WrongPassword
	WrongLogin
	NoPermission
)

var Ok = Error{Code: OK}

type Error struct {
	Value string
	Code  int
	Place string
	Sub   error
}

func (e Error) Error() string {
	return fmt.Sprintf("Error[%s] in %s: %s", e.Code, e.Place, e.Value)
}
func (e Error) SetCode(code int) Error {
	e.Code = code
	return e
}
func (e Error) SetError(value error) Error {
	return e.SetString(value.Error())
}
func (e Error) SetString(value string) Error {
	e.Value = value
	return e
}

func (e Error) SetPlace(place string) Error {
	e.Place = place
	return e
}

func (e Error) IncludeError(err error) Error {
	e.Sub = err
	return e
}

func (e Error) Unwrap() error {
	return e.Sub
}
