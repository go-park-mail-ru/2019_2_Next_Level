package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"2019_2_Next_Level/pkg/sqlTools"
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

func (r *PostgresRepository) GetUser(login string) (model.User, error) {
	user := model.User{}
	query := `SELECT login, firstname, secondname, sex, avatar, birthdate FROM users WHERE login=$1`
	res := r.DB.QueryRow(query, login)
	if res == nil {
		return user, e.Error{}.SetCode(e.ProcessError)
	}
	var birthDateRaw time.Time
	err := res.Scan(&user.Email, &user.Name, &user.Sirname, &user.Sex, &user.Avatar, &birthDateRaw)
	if err != nil {
		return user, e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	user.BirthDate = birthDateRaw.Format("02.01.2006")
	return user, nil
}

func (r *PostgresRepository) UpdateUserData(user *model.User) error {
	query := `UPDATE users SET avatar=$1, firstName=$2, secondname=$3, sex=$4, birthdate=$5 WHERE login=$6;`
	parsedDate, err0 := time.Parse("02.01.2006", user.BirthDate)
	if err0 != nil {
		return e.Error{}.SetCode(e.InvalidParams).SetError(err0)
	}
	_, err := r.DB.Exec(query, user.Avatar, user.Name, user.Sirname, user.Sex, sqlTools.FormatDate(sqlTools.BDPostgres, parsedDate), user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateUserPassword(login string, newPassword string, sault string) error {
	query := `UPDATE users SET password=$1, sault=$2 WHERE login=$3`
	_, err := r.DB.Exec(query, []byte(newPassword), []byte(sault), login)
	return err
}

func (r *PostgresRepository) GetUserCredentials(login string) (string, string, error) {
	query := "SELECT password, sault FROM users WHERE login=$1;"
	row := r.DB.QueryRow(query, login)

	var pass, sault string
	err := row.Scan(&pass, &sault)
	if err != nil {
		return "", "", e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return pass, sault, nil
}
