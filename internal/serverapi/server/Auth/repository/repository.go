package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"database/sql"
	"fmt"
	"time"

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
	query := `INSERT INTO users (login, password, sault, firstname, secondname, sex, birthdate, avatar)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8);`
	// birthDateSeparated := strings.Split(user.BirthDate, ".")
	// if len(birthDateSeparated) == 3 {
	// 	user.BirthDate = fmt.Sprintf("%s/%s/%s", birthDateSeparated[2], birthDateSeparated[1], birthDateSeparated[0])
	// 	user.BirthDate = birthDateSeparated[2] + "/" + birthDateSeparated[1] + "/" + birthDateSeparated[0]
	// }
	parsedDate, err0 := time.Parse("02.01.2006", user.BirthDate)
	if err0 != nil {
		return e.Error{}.SetCode(e.InvalidParams).SetError(err0)
	}
	user.BirthDate = parsedDate.Format("2006/01/02")
	_, err := r.DB.Exec(query, user.Email, user.Password, "", user.Name, user.Sirname, user.Sex, user.BirthDate, user.Avatar)
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
		fmt.Println(err)
		isExist = false
	}
	return isExist
}
