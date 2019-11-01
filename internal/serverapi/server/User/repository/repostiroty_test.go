package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/DATA-DOG/go-sqlmock"
)

var defaultConf config.Database

func init() {
	defaultConf = config.Database{DBName: "nextlevel", Port: "5432", Host: "localhost", User: "postgres", Password: "postgres"}
	config.Conf.DB = defaultConf
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
			rows := sqlmock.NewRows([]string{"password", "sault"}).AddRow("password", "sault")
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
		{[]string{"password", "sault"}, nil},
		{[]string{"", ""}, e.Error{}},
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
		pass, saut, err := repo.GetUserCredentials(param)
		if expected[i].Ans[0] != pass || expected[i].Ans[1] != saut {
			t.Errorf("Wrong answer: %v instead %v", []string{pass, saut}, expected[i].Ans)
		}
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instead %v", err, expected[i].Err)
		}
	}
}

func TestGetUser(t *testing.T) {

	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `SELECT login\, firstname\, secondname\, sex\, avatar\, birthdate FROM users WHERE login=\$1`
	user := model.User{Email: "ivanovivan", Name: "Anonim", Sirname: "Noone",
		Sex: "male", Avatar: "my_ava.png", BirthDate: "01.01.1274"}

	setDB := []func(sqlmock.Sqlmock) string{
		func(mock sqlmock.Sqlmock) string {
			parsed, _ := time.Parse("02.01.2006", user.BirthDate)
			rows := sqlmock.NewRows([]string{"login", "firstname", "secondname", "sex", "avatar", "birthdate"}).
				AddRow(user.Email, user.Name, user.Sirname, user.Sex, user.Avatar, parsed)
			mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)
			return user.Email
		},
		func(mock sqlmock.Sqlmock) string {
			rows := sqlmock.NewRows([]string{"login", "firstname", "secondname", "sex", "avatar", "birthdate"})
			// AddRow(user.Email, user.Name, user.Sirname, user.Sex, user.Avatar, user.BirthDate)
			mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)
			return user.Email
		},
	}

	expected := []struct {
		Ans model.User
		Err error
	}{
		{user, nil},
		{model.User{}, e.Error{}.SetCode(e.NotExists)},
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
		gotUser, err := repo.GetUser(param)
		if !cmp.Equal(gotUser, expected[i].Ans) {
			t.Errorf("Wrong answer: %v instead %v", gotUser, expected[i].Ans)
		}
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instead %v", err, expected[i].Err)
		}
	}
}
