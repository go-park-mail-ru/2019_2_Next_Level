package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/pkg/Error"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type PostgresRepository struct {
	DB     *sql.DB
	dbName string
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
	r.dbName = conf.DBName
	return nil
}

func (r *PostgresRepository) AddNewSession(login, uuid string) error {
	if ok := r.checkUserExist(login); !ok {
		return e.Error{}.SetCode(e.NotExists)
	}
	checkExistQuery := "DELETE FROM session WHERE login=$1"
	query := fmt.Sprintf("INSERT INTO %s (login, token) VALUES ($1, $2);", "session")
	_, err := r.DB.Exec(checkExistQuery, login)

	_, err = r.DB.Exec(query, login, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteSession(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token=$1", "session")
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return nil
}

func (r *PostgresRepository) GetLoginBySession(uuid string) (string, error) {
	query := fmt.Sprintf("SELECT login FROM %s WHERE token=$1", "session")
	row := r.DB.QueryRow(query, uuid)
	if row == nil {
		return "", e.Error{}.SetCode(e.NotExists)
	}
	login := ""
	err := row.Scan(&login)
	if err != nil {
		return "", e.Error{}.SetCode(e.NotExists)
	}
	return login, nil
}

func (r *PostgresRepository) AddNewUser(user *model.User) error {
	query := `INSERT INTO users (login, password, sault, firstname, secondname, avatar)
				VALUES($1, $2, $3, $4, $5, $6);`

	_, err := r.DB.Exec(query, user.Email, []byte(user.Password), []byte(user.Sault), user.Name, user.Sirname, user.Avatar)
	if err != nil {
		return e.Error{}.SetCode(e.AlreadyExists).SetError(err)
	}
	return nil
}

func (r *PostgresRepository) GetUserCredentials(login string) ([]string, error) {
	query := "SELECT password, sault FROM users WHERE login=$1;"
	row := r.DB.QueryRow(query, login)

	var pass, sault string
	err := row.Scan(&pass, &sault)
	if err != nil {
		return nil, e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return []string{pass, sault}, nil
}

func (r *PostgresRepository) checkUserExist(login string) bool {
	row := r.DB.QueryRow("SELECT COUNT(login)>0 FROM users WHERE login=$1", login)
	var isExist bool
	err := row.Scan(&isExist)
	if err != nil {
		isExist = false
	}
	return isExist
}
