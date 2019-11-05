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
	Wrong
)

var Ok = Error{Code: OK}

type Error struct {
	Value string
	Code  int
	Place string
	Sub   error
}

func (e Error) Error() string {
	return fmt.Sprintf("Error[%d] in %s: %s", e.Code, e.Place, e.Value)
}
func (e Error) SetCode(code int) Error {
	e.Code = code
	return e
}
func (e Error) SetError(value error) Error {
	e.Sub = value
	return e
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

func CompareErrors(e1, e2 error, comparator func(Error, Error) bool) bool {
	// if nil or any equal type
	if e1 == e2 {
		return true
	}
	E1, ok1 := e1.(Error)
	E2, ok2 := e2.(Error)
	if !ok1 || !ok2 {
		return false
	}
	return comparator(E1, E2)
}

func CompareByCode(e1, e2 Error) bool {
	return e1.Code == e2.Code
}
