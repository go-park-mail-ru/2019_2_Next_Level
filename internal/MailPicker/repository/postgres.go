package repository

import (
	"2019_2_Next_Level/internal/MailPicker/config"
	log "2019_2_Next_Level/internal/MailPicker/logger"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/pkg/sqlTools"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository() *PostgresRepository {
	repo := PostgresRepository{}
	err := repo.init()
	if err != nil {
		log.Log().E(err)
		return nil
	}
	return &repo
}

func (r *PostgresRepository) init() error {
	conf := &config.Conf.DB
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

	err = r.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

// UserExists : checks if the user exists
func (r *PostgresRepository) UserExists(login string) bool {
	query := `SELECT COUNT(login)>0 FROM users WHERE login=$1`
	row := r.DB.QueryRow(query, login)
	if row == nil {
		return false
	}
	var isExist bool
	err := row.Scan(&isExist)
	if err != nil {
		isExist = false
	}
	return isExist
}

// AddEmail : Inserts the new email to database
func (r *PostgresRepository) AddEmail(email *model.Email) error {
	task := func() error {
		saveMessage, err := r.DB.Prepare(`INSERT INTO Message (sender, time, body) VALUES ($1, $2, $3) RETURNING id;`)
		if err != nil {
			return err
		}

		var id int
		whenReceived := sqlTools.FormatDate(sqlTools.BDPostgres, email.Header.WhenReceived)
		err = saveMessage.QueryRow(email.From, whenReceived, email.Body).Scan(&id)
		if err != nil {
			return err
		}

		recQuery := sqlTools.CreatePacketQuery(`INSERT INTO Receiver (mailId, email) VALUES`, 2, len(email.Header.To))
		params := make([]interface{}, 0, 2*len(email.Header.To))
		for _, addr := range email.Header.To {
			params = append(params, id, addr)
		}

		saveReceivers, err := r.DB.Prepare(recQuery)
		if err != nil {
			return err
		}
		_, err = saveReceivers.Exec(params...)
		return err
	}

	return sqlTools.WithTransaction(r.DB, task)
}