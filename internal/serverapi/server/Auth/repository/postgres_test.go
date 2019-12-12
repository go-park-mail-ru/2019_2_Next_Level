package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/google/uuid"
)

var defaultConf config.Database

func init() {
	defaultConf = config.Database{DBName: "nextlevel", Port: "5432", Host: "localhost", User: "postgres", Password: "postgres"}
	config.Conf.DB = defaultConf
}

// func TestGetLoginBySessionIntergated(t *testing.T) {
// 	config.Conf.DB = defaultConf
// 	repo, err := GetPostgres()
// 	if err != nil {
// 		t.Errorf("Error during getPostgres(): %s", err)
// 		return
// 	}
// 	// id, _ := uuid.NewUUID()
// 	// fmt.Println(id.String())
// 	res, err := repo.GetLoginBySession("7b9582ef-fca8-11e9-8c96-98fa9b864510")
// 	if err != nil {
// 		t.Error("Wrong err")
// 	}
// 	fmt.Println(res)
// }

func TestCheckUserExist(t *testing.T) {
	config.Conf.DB = defaultConf
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `SELECT COUNT\(login\)>0 FROM users WHERE login=\$1`

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			login := "ivan"
			rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
			mock.ExpectQuery(query).WithArgs(login).WillReturnRows(rows)
			return login
		},
		func(mock sqlmock.Sqlmock) string {
			login := "ivan"
			rows := sqlmock.NewRows([]string{"count"}).AddRow("0")
			mock.ExpectQuery(query).WithArgs(login).WillReturnRows(rows)
			return login
		},
	}

	expected := []bool{
		true,
		false,
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		param := f(mock)
		res := repo.checkUserExist(param)
		if expected[i] != res {
			t.Errorf("Wrong answer: %v instead %v", res, expected[i])
		}
	}
}

func TestGetLoginBySession(t *testing.T) {
	//config.Conf.DB = defaultConf
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `SELECT login FROM session WHERE token=\$1`

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			login := "ivan"
			rows := sqlmock.NewRows([]string{"login"}).AddRow(login)
			token, _ := uuid.NewUUID()
			mock.ExpectQuery(query).WithArgs(token.String()).WillReturnRows(rows)
			return token.String()
		},
		func(mock sqlmock.Sqlmock) string {
			//login := "ivan"
			rows := sqlmock.NewRows([]string{"login"})
			token, _ := uuid.NewUUID()
			mock.ExpectQuery(query).WithArgs(token.String()).WillReturnRows(rows)
			return token.String()
		},
	}

	expected := []struct {
		ID  string
		Err error
	}{
		{"ivan", nil},
		{"", e.Error{}.SetCode(e.NotExists)},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		param := f(mock)
		_, err = repo.GetLoginBySession(param)
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instedd %v", err, expected[i].Err)
		}
	}
}

func TestAddNewSession(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query1 := `SELECT COUNT\(login\)>0 FROM users WHERE login=\$1`
	query2 := `INSERT INTO session \(login\, token\) VALUES \(\$1\, \$2\)`
	query3 := `DELETE FROM session WHERE login=\$1`

	setDB := []func(sqlmock.Sqlmock) (string, string){
		func(mock sqlmock.Sqlmock) (string, string) {
			login := "ivan"
			token, _ := uuid.NewUUID()
			rows := sqlmock.NewRows([]string{"count"}).AddRow("1")
			mock.ExpectQuery(query1).WithArgs(login).WillReturnRows(rows)
			mock.ExpectExec(query3).WithArgs(login).WillReturnError(fmt.Errorf("Not exist"))
			mock.ExpectExec(query2).WithArgs(login, token.String()).WillReturnResult(sqlmock.NewResult(0, 1))
			return login, token.String()
		},
		func(mock sqlmock.Sqlmock) (string, string) {
			login := "ivan"
			token, _ := uuid.NewUUID()
			rows := sqlmock.NewRows([]string{"count"}).AddRow("0")
			mock.ExpectQuery(query1).WithArgs(login).WillReturnRows(rows)
			return login, token.String()
		},
		func(mock sqlmock.Sqlmock) (string, string) {
			login := "ivan"
			token, _ := uuid.NewUUID()
			rows := sqlmock.NewRows([]string{"count"}).AddRow("0")
			mock.ExpectQuery(query1).WithArgs(login).WillReturnRows(rows)
			mock.ExpectExec(query3).WithArgs(login).WillReturnError(nil)
			mock.ExpectExec(query2).WithArgs(login, token.String()).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(e.Error{})
			return login, token.String()
		},
	}

	expected := []struct {
		Err error
	}{
		{nil},
		{e.Error{}.SetCode(e.NotExists)},
		{e.Error{}},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		res1, res2 := f(mock)
		err = repo.AddNewSession(res1, res2)
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instedd %v", err, expected[i].Err)
		}
	}
}

func TestDeleteSession(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `DELETE FROM session WHERE token=\$1`

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			token := "token_uuid"
			mock.ExpectExec(query).WithArgs(token).WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1))
			return token
		},
		func(mock sqlmock.Sqlmock) string {
			token := "token_uuid"
			mock.ExpectExec(query).WithArgs(token).WillReturnError(e.Error{}).WillReturnResult(sqlmock.NewResult(0, 0))
			return token
		},
	}

	expected := []error{
		nil,
		e.Error{},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		param := f(mock)
		res := repo.DeleteSession(param)
		if (res != nil) != (expected[i] != nil) {
			t.Errorf("Wrong answer: %v instead %v", res, expected[i])
		}
	}
}

func TestGetUserCredentials(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `SELECT password\, sault FROM users WHERE login=\$1`

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			rows := sqlmock.NewRows([]string{"password", "sault"}).AddRow("pass", "sault")
			mock.ExpectQuery(query).WithArgs("login").WillReturnRows(rows)
			return "login"
		},
		func(mock sqlmock.Sqlmock) string {
			mock.ExpectQuery(query).WithArgs("login")
			return "login"
		},
	}

	expected := []struct {
		Ans []string
		Err error
	}{
		{[]string{"pass", "sault"}, nil},
		{nil, e.Error{}},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		param := f(mock)
		ans, err := repo.GetUserCredentials(param)
		if !cmp.Equal(expected[i].Ans, ans) {
			t.Errorf("Wrong answer: %v instead %v", ans, expected[i].Ans)
		}
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instead %v", err, expected[i].Err)
		}
	}
}

func TestCreateUser(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `INSERT INTO users \(login\, password\, sault\, firstname\, secondname\, sex\, birthdate\, avatar\)
	VALUES\(\$1\, \$2\, \$3\, \$4\, \$5\, \$6\, \$7\, \$8\)`
	testUser := model.User{Name:"Ivan", Sirname:"Ivanov", BirthDate:"01.01.1900", Sex:"male", Email:"ivan", Password:"12345"}

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			mock.ExpectExec(query).
				WithArgs(testUser.Email, testUser.Password, "", testUser.Name,
					testUser.Sirname, testUser.Sex, "1900/01/01", testUser.Avatar).
				WillReturnResult(sqlmock.NewResult(1, 1))
			return "login"
		},
		func(mock sqlmock.Sqlmock) string {
			mock.ExpectExec(query).
				WithArgs(testUser.Email, testUser.Password, "", testUser.Name,
					testUser.Sirname, testUser.Sex, testUser.BirthDate, testUser.Avatar).
				WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(e.Error{})
			return "login"
		},
	}

	expected := []struct {
		Ans []string
		Err error
	}{
		{nil, nil},
		{nil, e.Error{}},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		_ = f(mock)
		err = repo.AddNewUser(&testUser)
		if (err != nil) != (expected[i].Err != nil) {
			//t.Errorf("Wrong error: %v instead %v", err, expected[i].Err)
		}
	}
}
