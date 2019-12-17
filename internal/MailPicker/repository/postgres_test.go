package repository

import (
	"2019_2_Next_Level/internal/MailPicker/config"
	"2019_2_Next_Level/internal/MailPicker/log"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/pkg/sqlTools"
	"2019_2_Next_Level/tests/mock"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

var defaultConf config.Database

func init() {
	defaultConf = config.Database{DBName: "nextlevel", Port: "5432", Host: "localhost", User: "postgres", Password: "postgres"}
	config.Conf.DB = defaultConf
	log.SetLogger(&mock.MockLog{})
}

func TestPostgresRepository_UserExists(t *testing.T) {
	repo := NewPostgresRepository()
	if repo == nil {
		t.Errorf("Cannot init repository")
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
		repo.DB = db
		param := f(mock)
		res := repo.UserExists(param)
		if expected[i] != res {
			t.Errorf("Wrong answer: %v instead %v", res, expected[i])
		}
		db.Close()
	}
}

func TestPostgresRepository_AddEmail(t *testing.T) {
	repo := NewPostgresRepository()
	if repo == nil {
		t.Errorf("Cannot init repository")
		return
	}
	email := model.Email{
		From: "testsender",
		To: "testreceiver",
		Body: "testbody",
		Header: struct {
			From string
			To   []string
			Subject string
			ReplyTo []string
			WhenReceived time.Time
		}{From: "testsender", Subject:"S", To: []string{"testreceiver"}, WhenReceived: time.Now(),},
	}
	query := []string{
		`INSERT INTO Message \(sender\, time\, body\, subject\) VALUES \(\$1\, \$2\, \$3\, \$4\) RETURNING id`,
		`INSERT INTO Receiver \(mailId\, email\) VALUES \(\$1\, \$2\)`,
	}

	setDB := []func(sqlmock.Sqlmock) error {
		// Good work
		func(mock sqlmock.Sqlmock) error {
			id := sqlmock.NewRows([]string{"id"}).AddRow("123")
			mock.ExpectBegin()
			mock.ExpectPrepare(query[0]).ExpectQuery().
				WithArgs(email.From, sqlTools.FormatDate(sqlTools.BDPostgres, email.Header.WhenReceived), email.Body, email.Header.Subject).
				WillReturnRows(id)
			mock.ExpectPrepare(query[1]).ExpectExec().WithArgs(123, email.Header.To[0]).WillReturnResult(sqlmock.NewResult(0,1))
			mock.ExpectCommit()
			return nil
		},
		// Error preparing saveMessage query
		func(mock sqlmock.Sqlmock) error {
			mock.ExpectBegin()
			mock.ExpectPrepare(query[0]).WillReturnError(fmt.Errorf(""))
			mock.ExpectRollback()
			return nil
		},
		// Error quering saveMessage
		func(mock sqlmock.Sqlmock) error {
			mock.ExpectBegin()
			mock.ExpectPrepare(query[0]).ExpectQuery().
				WithArgs(email.From, sqlTools.FormatDate(sqlTools.BDPostgres, email.Header.WhenReceived), email.Body, email.Header.Subject).
				WillReturnError(fmt.Errorf(""))
			mock.ExpectRollback()
			return nil
		},
		// Error preparing saveReceivers query
		func(mock sqlmock.Sqlmock) error {
			id := sqlmock.NewRows([]string{"id"}).AddRow("123")
			mock.ExpectBegin()
			mock.ExpectPrepare(query[0]).ExpectQuery().
				WithArgs(email.From, sqlTools.FormatDate(sqlTools.BDPostgres, email.Header.WhenReceived), email.Body).
				WillReturnRows(id)
			mock.ExpectPrepare(query[1]).WillReturnError(fmt.Errorf(""))
			mock.ExpectRollback()
			return nil
		},
		// Error during exec saveReceivers
		func(mock sqlmock.Sqlmock) error {
			id := sqlmock.NewRows([]string{"id"}).AddRow("123")
			mock.ExpectBegin()
			mock.ExpectPrepare(query[0]).ExpectQuery().
				WithArgs(email.From, sqlTools.FormatDate(sqlTools.BDPostgres, email.Header.WhenReceived), email.Body, email.Header.Subject).
				WillReturnRows(id)
			mock.ExpectPrepare(query[1]).ExpectExec().WithArgs(123, email.Header.To[0]).WillReturnError(fmt.Errorf(""))
			mock.ExpectRollback()
			return nil
		},
	}
	expected := []error{
		nil,
		fmt.Errorf(""),
		fmt.Errorf(""),
		fmt.Errorf(""),
		fmt.Errorf(""),
	}

	for i, f := range setDB {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("Error doring init sqlmock")
			return
		}
		repo.DB = db
		f(mock)
		err = repo.AddEmail(&email)

		if (err == nil) != (expected[i] == nil) {
			t.Errorf("Wrong error returned: %v instead %v", err, expected[i])
		}

		db.Close()
	}
}