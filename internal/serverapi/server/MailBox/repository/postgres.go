package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"2019_2_Next_Level/pkg/sqlTools"
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
	query := queryGetEmailByCode

	mail := model.Email{}
	var when string
	id, _ := strconv.ParseInt(code.(string), 10, 8)
	err := r.DB.QueryRow(query, id).Scan(&mail.From, &mail.To, &when, &mail.Header.Subject, &mail.Body)
	if err != nil {
		return mail, e.Error{}.SetError(err)
	}
	mail.Header.WhenReceived, _ = time.Parse("2006/01/02 15:04:05", when)
	return mail, nil
}

func (r *PostgresRepository) GetEmailList(login string, folder string, sort interface{}, firstNumber int, count int) ([]model.Email, error) {
	query := queryGetEmailList
	var placeholder string

	if folder=="sent" || folder=="proceed"{
		placeholder = "Message.sender"
	}else{
		placeholder = "Receiver.email"
	}
	query = fmt.Sprintf(query, placeholder)

	row, err := r.DB.Query(query, login, count, firstNumber-1)
	list := make([]model.Email, 0)
	if err != nil {
		return list, e.Error{}.SetCode(e.NotExists)
	}

	for row.Next() {
		mail := model.Email{}
		var when string

		if err := row.Scan(&mail.Id, &mail.From, &mail.To, &when, &mail.Header.Subject, &mail.Body, &mail.IsRead); err != nil {
			return list, e.Error{}.SetError(err)
		}

		mail.Header.WhenReceived, _= time.Parse(time.RFC3339, when)
		list = append(list, mail)
	}
	return list, nil
}

func (r *PostgresRepository) GetMessagesCount(login string, folder string, flag interface{}) (int, error) {
	query := queryGetMessagesCount
	var count int
	err := r.DB.QueryRow(query, login).Scan(&count)
	return count, err
}

func (r *PostgresRepository) MarkMessages(login string, messagesID []models.MailID, mark interface{}) error {
	query := queryMarkMessage
	var placeholder string
	switch mark.(int) {
	case models.MarkMessageRead:
		placeholder = `isread=true`
	case models.MarkMessageUnread:
		placeholder = `isread=false`
	case models.MarkMessageDeleted:
		placeholder = `folder='trash'`
	default:
		return fmt.Errorf("Unknown mark")
	}
	query = fmt.Sprintf(query, placeholder)

	for _, id := range messagesID {
		_, err := r.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) PutSentMessage(email model.Email) error {
	task := func() error {
		saveMessage, err := r.DB.Prepare(queryWriteMessage)
		if err != nil {
			return err
		}

		var id int
		from, _ := email.Split(email.From)
		err = saveMessage.QueryRow(from, email.Header.Subject, email.Body, "out", "sent").Scan(&id)
		if err != nil {
			return err
		}

		recQuery := sqlTools.CreatePacketQuery(queryInflateReceivers, 2, len(email.Header.To))
		saveReceivers, err := r.DB.Prepare(recQuery)
		if err != nil {
			return err
		}

		params := make([]interface{}, 0, 2*len(email.Header.To))
		for _, addr := range email.Header.To {
			params = append(params, id, addr)
		}
		_, err = saveReceivers.Exec(params...)
		return err
	}

	return sqlTools.WithTransaction(r.DB, task)
}