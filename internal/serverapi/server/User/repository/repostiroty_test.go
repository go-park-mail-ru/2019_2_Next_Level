package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/pkg/Error"
	"2019_2_Next_Level/pkg/TestTools"
	"2019_2_Next_Level/pkg/sqlTools"
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

func TestPostgresRepository_GetUserFolders(t *testing.T) {

	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `SELECT name\, count FROM Folder WHERE owner=\$1`
	login := "ian"

	setDB := []func(sqlmock.Sqlmock) []model.Folder{
		func(mock sqlmock.Sqlmock) []model.Folder {
			rows := sqlmock.NewRows([]string{"name", "count"}).
				AddRow("folder", 1)
			mock.ExpectQuery(query).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
			return []model.Folder{model.Folder{
				Name:         "folder",
				MessageCount: 1,
			}}
		},
		func(mock sqlmock.Sqlmock) []model.Folder {
			rows := sqlmock.NewRows([]string{"name", "count"}).
				AddRow("folder", 1)
			mock.ExpectQuery(query).WithArgs(login).WillReturnRows(rows).WillReturnError(e.Error{})
			return []model.Folder{}
		},
	}

	expected := []struct {
		Folder []model.Folder
		Err error
	}{
		{[]model.Folder{model.Folder{
			Name:         "folder",
			MessageCount: 1,
		}}, nil},
		{[]model.Folder{}, e.Error{}},
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		defer db.Close()
		repo.DB = db
		f(mock)
		gotUser, err := repo.GetUserFolders(login)
		if !cmp.Equal(gotUser, expected[i].Folder) {
			t.Errorf("Wrong answer: %v instead %v", gotUser, expected[i].Folder)
		}
		if (err != nil) != (expected[i].Err != nil) {
			t.Errorf("Wrong error: %v instead %v", err, expected[i].Err)
		}
	}
}

func TestPostgresRepository_UpdateUserData(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `UPDATE users SET avatar=\$1\, firstName=\$2\, secondname=\$3\, sex=\$4\, birthdate=\$5 WHERE login=\$6`
	user := model.User{Email: "ivanovivan", Name: "Anonim", Sirname: "Noone",
		Sex: "male", Avatar: "my_ava.png", BirthDate: "01.01.1274"}
	parsedDate, _ := time.Parse("02.01.2006", user.BirthDate)

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
		//params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		mock.ExpectExec(query).WithArgs(user.Avatar, user.Name, user.Sirname, user.Sex,
			sqlTools.FormatDate(sqlTools.BDPostgres, parsedDate), user.Email)

		_ = repo.UpdateUserData(&user)
	})
}

func TestPostgresRepository_UpdateUserPassword(t *testing.T) {
	repo, err := GetPostgres()
	if err != nil {
		t.Errorf("Error during getPostgres(): %s", err)
		return
	}
	query := `UPDATE users SET password=\$1\, sault=\$2 WHERE login=\$3`
	user := model.User{Email: "ivanovivan", Name: "Anonim", Sirname: "Noone",
		Sex: "male", Avatar: "my_ava.png", BirthDate: "01.01.1274"}

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"error": nil, "pass": []byte("pass"),"salt": []byte("salt")},
			map[string]TestTools.Params{
				"err":nil,
			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()
		var err error
		if params["err"] != nil {
			err = params["err"].(error)
		}
		mock.ExpectExec(query).WithArgs([]byte("login"), []byte("password"), []byte("sault")).
			WillReturnResult(sqlmock.NewResult(1,1)).
			WillReturnError(err)

		_ = repo.UpdateUserData(&user)



	})
}