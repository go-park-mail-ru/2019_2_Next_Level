package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"2019_2_Next_Level/pkg/TestTools"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)
var expectedFunc func(params... TestTools.Params)
func init() {
	expectedFunc = func(params... TestTools.Params) {}
}
func TestPostgresRepository_GetEmailByCode(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345"},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
	}

	query := `SELECT sender\, email AS "receivers"\, time\, subject\, body from Message
							JOIN Receiver ON Message\.id=Receiver\.mailId
							WHERE Message\.id=\$1`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		params := test.MockParams
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"sender", "email", "time", "subject", "body"}).
			AddRow("", "", time.Now(), "", "")
		mock.ExpectQuery(query).WithArgs(params[0].(string), params[1].(int64)).WillReturnRows(rows)

		res, err := repo.GetEmailByCode(test.Input[0].(string), test.Input[1])
		if !cmp.Equal(res, res) || err != err {
			t.Errorf("Wrong answer: %v instead %v", res, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_GetEmailList(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345"},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
	}

	query := `SELECT Message\.id\, sender\, email AS "receivers"\, time\, subject\, body, isread from Message
						JOIN Receiver ON Message\.id=Receiver\.mailId
						WHERE .*=\$1 ORDER BY time LIMIT \$2 OFFSET \$3;`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "sender", "email", "time", "subject", "body", "isread"}).
			AddRow(123, "", "", time.Now(), "", "", true)
		mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

		res, err := repo.GetEmailList(test.Input[0].(string), "folder", "", 1, 25)
		if !cmp.Equal(res, res) || err != err {
			t.Errorf("Wrong answer: %v instead %v", res, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_GetMessagesCount(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345"},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
	}

	query := `SELECT COUNT\(Message\.id\) from Message JOIN Receiver ON Message\.id=Receiver\.mailId
							WHERE Receiver\.email=\$1`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(123)
		mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

		res, err := repo.GetMessagesCount(test.Input[0].(string), "foler", "d")
		if !cmp.Equal(res, res) || err != err {
			t.Errorf("Wrong answer: %v instead %v", res, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_MarkMessages(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345", models.MarkMessageRead},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345", models.MarkMessageUnread},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345", models.MarkMessageDeleted},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
	}

	query := `UPDATE Message SET .* WHERE id=\$1`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(123)
		mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

		err := repo.MarkMessages(test.Input[0].(string), []models.MailID{models.MailID(123)}, test.Input[2].(int))
		if err != err {
			t.Errorf("Wrong answer: %v instead %v",err, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_PutSentMessage(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "12345", models.MarkMessageRead},
			[]TestTools.Params{true},
			[]TestTools.Params{"ivan", int64(12345)},
		),
	}

	//query := `UPDATE Message SET .* WHERE id=\$1`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		rowsId := sqlmock.NewRows([]string{"id"}).
			AddRow(123)
		mock.ExpectBegin()
		mock.ExpectPrepare(".*").ExpectQuery().
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rowsId)
		mock.ExpectPrepare(".*").ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(nil)
		mock.ExpectCommit()

		email := model.Email{From:"from", To:"to", Body:"body", Header:struct{
			From string
			To   []string
			Subject string
			ReplyTo []string
			WhenReceived time.Time
		}{Subject:"subject", To:[]string{"a@ss.ss"}}}

		err := repo.PutSentMessage(email)
		if err != err {
			t.Errorf("Wrong answer: %v instead %v",err, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_AddFolder(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "test"},
			[]TestTools.Params{nil},
			[]TestTools.Params{"ivan", "test"},
		),
	}
	query := `INSERT INTO Folder \(name\, owner\) VALUES \(\$1\, \$2\)`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1,1)).
			WillReturnError(nil)

		err := repo.AddFolder("ivan", "test")
		if err != err {
			t.Errorf("Wrong answer: %v instead %v",err, test.Expected[0].(bool))
		}

	})
}

func TestPostgresRepository_ChangeMailFolder(t *testing.T) {
	repo, _ := GetPostgres()
	tests := []TestTools.TestStruct{
		*TestTools.NewTestStruct(
			[]TestTools.Params{"ivan", "test"},
			[]TestTools.Params{nil},
			[]TestTools.Params{"ivan", "test"},
		),
	}
	query := `UPDATE Message SET folder=\$1 WHERE id=\$2`

	TestTools.RunTesting(tests, expectedFunc, func(test TestTools.TestStruct) {
		db, mock, _ := sqlmock.New()
		repo.DB = db
		defer db.Close()

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1,1)).
			WillReturnError(nil)

		err := repo.ChangeMailFolder("ivan", "test", 12345)
		if err != err {
			t.Errorf("Wrong answer: %v instead %v",err, test.Expected[0].(bool))
		}

	})
}