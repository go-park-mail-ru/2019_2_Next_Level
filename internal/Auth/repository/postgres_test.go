package repository

import (
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"2019_2_Next_Level/pkg/TestTools"
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	//"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	//"github.com/golang/mock/gomock"
	//"github.com/google/go-cmp/cmp"
	//"github.com/google/uuid"
	"testing"
)

const (
	dbuser = "user"
	dbpass = "pass"
	dbhost = "0.0.0.0"
	dbport = "12345"
	dbname = "dbtpark"
	port = "2000"
)
var expectedFunc func(params... TestTools.Params)
func init() {
	expectedFunc = func(params... TestTools.Params) {}
}

func TestCheckUserExist(t *testing.T) {
	repo := NewPostgresRepo()
	query := `SELECT COUNT\(login\)>0 FROM users WHERE login=\$1`
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
				[]TestTools.Params{"ivan"},
				[]TestTools.Params{true},
				[]TestTools.Params{"1", "ivan"},
			),
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan"},
			[]TestTools.Params{false},
			[]TestTools.Params{"0", "ivan"},
		),
	}

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(params[0].(string))
		mock.ExpectQuery(query).WithArgs(params[1].(string)).WillReturnRows(rows)

		res := repo.checkUserExist(test.Input[0].(string))
		if test.Expected[0].(bool) != res {
			t.Errorf("Wrong answer: %v instead %v", res, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepo_GetLoginBySession(t *testing.T) {
	repo := NewPostgresRepo()
	query := `SELECT login FROM session WHERE token=\$1`
	//login := "ivan"
	//sault := Auth.GenSault(login)
	//passHash := Auth.PasswordPBKDF2([]byte("oldpass"), sault)
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"session_token"},
			[]TestTools.Params{"ivan", nil},
			[]TestTools.Params{"ivan", "session_token"},
		),
		*TestTools.NewTestStruct(
			[]TestTools.Params{"session_token"},
			[]TestTools.Params{"", e.Error{}.SetCode(e.NotExists)},
			[]TestTools.Params{"", "session_token"},
		),
	}

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"login"})
		if params[0].(string) != "" {
			rows = rows.AddRow(params[0].(string))
		}
		mock.ExpectQuery(query).WithArgs(params[1].(string)).WillReturnRows(rows)

		login, err := repo.GetLoginBySession(test.Input[0].(string))
		//if !bytes.Equal(test.Expected[0].([]byte), pass) {
		//	t.Errorf("Wrong pass")
		//}
		//if !bytes.Equal(test.Expected[1].([]byte), salt) {
		//	t.Errorf("Wrong salt")
		//}
		if !cmp.Equal(test.Expected[0].(string), login) {
			t.Errorf("Wrong login")
		}
		if !cmp.Equal(test.Expected[1], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected[1])
		}

	})
}

func TestPostgresRepo_AddNewSession(t *testing.T) {
	repo := NewPostgresRepo()
	query1 := `SELECT COUNT\(login\)>0 FROM users WHERE login=\$1`
	query2 := `INSERT INTO session \(login\, token\) VALUES \(\$1\, \$2\)`
	query3 := `DELETE FROM session WHERE login=\$1`

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "token":"session_token"},
			map[string]TestTools.Params{"error":nil},
			map[string]TestTools.Params{
				"user_count": "1",
				"login":"login",
				"token":"session_token",
				"delete_err": fmt.Errorf("Not exist"),
				"add_err": nil,
			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "token":"session_token"},
			map[string]TestTools.Params{"error": e.Error{}.SetCode(e.NotExists)},
			map[string]TestTools.Params{
				"user_count": "0",
				"login":"login",
			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "token":"session_token"},
			map[string]TestTools.Params{"error": e.Error{}},
			map[string]TestTools.Params{
				"user_count": "1",
				"login":      "login",
				"token":      "session_token",
				"delete_err": nil,
				"add_err":    e.Error{},
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(params["user_count"].(string))
		mock.ExpectQuery(query1).WithArgs(params["login"].(string)).WillReturnRows(rows)
		if params["user_count"].(string) != "0"{
			var err error
			if params["delete_err"] == nil {
				err = nil
			} else {
				err = params["delete_err"].(error)
			}
			mock.ExpectExec(query3).WithArgs(params["login"].(string)).WillReturnError(err)
			if params["add_err"] == nil {
				err = nil
			} else {
				err = params["add_err"].(error)
			}
			mock.ExpectExec(query2).WithArgs(params["login"].(string), params["token"].(string)).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(err)
		}

		err := repo.AddNewSession(test.Input["login"].(string), test.Input["token"].(string))

		if !cmp.Equal(test.Expected["error"], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected["error"])
		}

	})
}

func TestPostgresRepo_DeleteSession(t *testing.T) {
	repo := NewPostgresRepo()
	query := `DELETE FROM Session WHERE token=\$1`

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"token":"session_token"},
			map[string]TestTools.Params{"error":nil},
			map[string]TestTools.Params{
				"token":"session_token",
				"id":int64(1),
				"rows":int64(1),
				"error":nil,
			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"token":"session_token"},
			map[string]TestTools.Params{"error": e.Error{Code: e.NotExists}.SetError(e.Error{})},
			map[string]TestTools.Params{
				"token": "session_token",
				"id":    int64(0),
				"rows":  int64(0),
				"error": e.Error{},
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()
		var err error
		if params["error"] != nil {
			err = params["error"].(error)
		}
		mock.ExpectExec(query).WithArgs(params["token"].(string)).
			WillReturnError(err).WillReturnResult(sqlmock.NewResult(params["id"].(int64), params["rows"].(int64)))

		err = repo.DeleteSession(params["token"].(string))

		if !cmp.Equal(test.Expected["error"], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected["error"])
		}

	})
}

func TestPostgresRepo_DeleteUserSessions(t *testing.T) {
	repo := NewPostgresRepo()
	query := `DELETE FROM Session WHERE login=\$1`

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"error":nil},
			map[string]TestTools.Params{
				"login":"login",
				"id":int64(1),
				"rows":int64(1),
				"error":nil,
			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"error": e.Error{}.SetCode(e.NotExists).SetError(e.Error{})},
			map[string]TestTools.Params{
				"login": "login",
				"id":    int64(0),
				"rows":  int64(0),
				"error": e.Error{},
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()
		var err error
		if params["error"] != nil {
			err = params["error"].(error)
		}
		mock.ExpectExec(query).WithArgs(params["login"].(string)).
			WillReturnError(err).WillReturnResult(sqlmock.NewResult(params["id"].(int64), params["rows"].(int64)))

		err = repo.DeleteUserSessions(params["login"].(string))

		if !cmp.Equal(test.Expected["error"], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected["error"])
		}

	})
}

func TestPostgresRepo_GetUserCredentials(t *testing.T) {
	repo := NewPostgresRepo()
	query := `SELECT password\, sault FROM users WHERE login=\$1`

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"error": nil, "pass": []byte("pass"),"salt": []byte("salt")},
			map[string]TestTools.Params{
				"login":"login",
				"pass": []byte("pass"),
				"salt": []byte("salt"),
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()
		var err error
		if params["error"] != nil {
			err = params["error"].(error)
		}
		rows := sqlmock.NewRows([]string{"password", "sault"})
		if params["pass"] != nil {
			rows.AddRow(params["pass"].([]byte), params["salt"].([]byte))
		}
		mock.ExpectQuery(query).WithArgs(params["login"].(string)).WillReturnRows(rows)

		pass, salt, err := repo.GetUserCredentials(params["login"].(string))

		if !bytes.Equal(pass, test.Expected["pass"].([]byte)) {
			t.Errorf("Wrong pass")
		}
		if !bytes.Equal(salt, test.Expected["salt"].([]byte)) {
			t.Errorf("Wrong salt")
		}
		if !cmp.Equal(test.Expected["error"], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected["error"])
		}

	})
}

func TestPostgresRepo_UpdateUserPassword(t *testing.T) {
	repo := NewPostgresRepo()
	query := `UPDATE users SET password=\$1, sault=\$2 WHERE login=\$3`

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "pass": []byte("pass"), "salt":[]byte("salt")},
			map[string]TestTools.Params{"error":nil},
			map[string]TestTools.Params{
				"login":"login",
				"pass": []byte("pass"),
				"salt": []byte("salt"),
				"res": int64(0),
				"error": nil,
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()
		var err error
		if params["error"] != nil {
			err = params["error"].(error)
		}

		mock.ExpectExec(query).
			WithArgs(params["pass"].([]byte), params["salt"].([]byte), params["login"].(string)).
			WillReturnResult(sqlmock.NewResult(params["res"].(int64), params["res"].(int64))).
			WillReturnError(err)

		err = repo.UpdateUserPassword(params["login"].(string), params["pass"].([]byte), params["salt"].([]byte))

		if !cmp.Equal(test.Expected["error"], err) {
			t.Errorf("Wrong error: %v instead %v", err, test.Expected["error"])
		}

	})
}
