package repository

import (
	"2019_2_Next_Level/internal/MailPicker/config"
	"2019_2_Next_Level/internal/model"
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
		fmt.Println(err)
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

func (r *PostgresRepository) AddEmail(email *model.Email) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	err = func() error {
		saveMessage, err := r.DB.Prepare(`INSERT INTO Message (sender, time, body) VALUES ($1, $2, $3) RETURNING id;`)
		if err != nil {
			return err
		}
		var id int
		whenReceived := email.Header.WhenReceived.Format("2006-01-02 15:04:05")
		err = tx.Stmt(saveMessage).QueryRow(email.From, whenReceived, email.Body).Scan(&id)
		if err != nil {
			return err
		}
		params := make([]interface{}, 0, 2*len(email.Header.To))
		recQuery := `INSERT INTO Receiver (mailId, email) VALUES `
		for i, addr := range email.Header.To {
			params = append(params, id, addr)
			recQuery = recQuery + fmt.Sprintf("($%d, $%d),", 2*i+1, 2*i+2)
		}
		recQuery = recQuery[:len(recQuery)-1] + ";"
		saveReceivers, err := r.DB.Prepare(recQuery)
		if err != nil {
			return err
		}
		_, err = tx.Stmt(saveReceivers).Exec(params...)
		return err
	}()

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}