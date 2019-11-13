package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"strconv"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

func GetPostgres() (PostgresRepository, error) {
	r := PostgresRepository{}
	if r.DB == nil {
		err := r.Init()
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (r *PostgresRepository) Init() error {
	conf := &config.Conf.DB
	// dsn = "postgres://postgres:postgres@localhost:5432/test?sslmode=disable"
	dsnTemplate := "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	user := conf.User
	password := conf.Password
	host := conf.Host
	port := conf.Port
	dbname := conf.DBName
	dsn := fmt.Sprintf(dsnTemplate, user, password, host, port, dbname)
	var err error
	r.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	// db.SetMaxOpenConns(10)

	err = r.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetEmailByCode(login string, code interface{}) (model.Email, error) {
	query := `SELECT sender, email AS "receivers", time, body from Message JOIN Receiver ON Message.id=Receiver.mailId
				WHERE Message.id=$1`
	mail := model.Email{}
	var when string
	id, _ := strconv.ParseInt(code.(string), 10, 8)
	err := r.DB.QueryRow(query, id).Scan(&mail.From, &mail.To, &when, &mail.Body)
	if err != nil {
		return mail, e.Error{}.SetError(err)
	}
	mail.Header.WhenReceived, _ = time.Parse("2006/01/02 15:04:05", when)
	return mail, nil
}

func (r *PostgresRepository) GetEmailList(login string, folder string, sort interface{}, firstNumber int, count int) ([]model.Email, error) {
	query := `SELECT Message.id, sender, email AS "receivers", time, body from Message JOIN Receiver ON Message.id=Receiver.mailId
				WHERE Receiver.email=$1 ORDER BY time LIMIT $2 OFFSET $3;`

	row, err := r.DB.Query(query, login, count, firstNumber-1)
	list := make([]model.Email, 0)
	if err != nil {
		return list, e.Error{}.SetCode(e.NotExists)
	}
	for row.Next() {
		mail := model.Email{}
		var when string
		err := row.Scan(&mail.Id, &mail.From, &mail.To, &when, &mail.Body)
		if err != nil {
			return list, e.Error{}.SetError(err)
		}
		mail.Header.WhenReceived, _ = time.Parse("2006/01/02 15:04:05", when)
		list = append(list, mail)
	}
	return list, nil
}

func (r *PostgresRepository) GetMessagesCount(login string, folder string, flag interface{}) (int, error) {
	query := `SELECT COUNT(Message.id) from Message JOIN Receiver ON Message.id=Receiver.mailId
				WHERE Receiver.email=$1`
	var count int
	err := r.DB.QueryRow(query, login).Scan(&count)
	return count, err
}

func (r *PostgresRepository) MarkMessages(login string, messagesID []models.MailID, mark interface{}) error {
	return nil
}