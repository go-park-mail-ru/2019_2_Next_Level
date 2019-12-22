package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/pkg/Error"
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
	local := "User.Repository.GetUser"
	user := model.User{}
	query := `SELECT login, firstname, secondname, sex, avatar, birthdate FROM users WHERE login=$1`
	res := r.DB.QueryRow(query, login)
	if res == nil {
		return user, e.Error{}.SetCode(e.ProcessError).SetPlace(local)
	}
	var birthDateRaw time.Time
	err := res.Scan(&user.Email, &user.Name, &user.Sirname, &user.Sex, &user.Avatar, &birthDateRaw)
	if err != nil {
		return user, e.Error{}.SetCode(e.NotExists).SetError(err).SetPlace(local)
	}
	user.BirthDate = birthDateRaw.Format("02.01.2006")
	return user, nil
}

func (r *PostgresRepository) GetUserFolders(login string) ([]model.Folder, error) {
	local := "User.Repository.GetUserFolders"
	folders :=make([]model.Folder, 0)
	query := `SELECT name, count, isSystem FROM Folder WHERE owner=$1 ORDER BY id`
	rows, err := r.DB.Query(query, login)
	if err != nil {
		return folders, e.Error{}.SetCode(e.ProcessError).SetPlace(local).SetError(err)
	}
	for rows.Next(){
		var folder model.Folder
		err := rows.Scan(&folder.Name, &folder.MessageCount, &folder.IsSystem)
		if err != nil {
			return folders, e.Error{}.SetCode(e.ProcessError).SetError(err).SetPlace(local)
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (r *PostgresRepository) UpdateUserData(user *model.User) error {
	//query := `UPDATE users SET avatar=$1, firstName=$2, secondname=$3 WHERE login=$4;`
	f := func (name, value string) error{
		query := `UPDATE users SET %s=$1 WHERE login=$2`
		if value != ""{
			_, err := r.DB.Exec(fmt.Sprintf(query, name), value, user.Email)
			//log.Log().L("Result user.repo:99 ", err);
			return err
		}
		return nil
	}
	inflate := func () error {
		err := f("firstName", user.Name)
		if err != nil {
			return err
		}
		err = f("secondName", user.Sirname)
		if err != nil {
			return err
		}
		err = f("avatar", user.Avatar)
		if err != nil {
			return err
		}
		return nil
	}
	if err:=inflate(); err != nil {
		return e.Error{}.SetError(err).SetPlace("User.Repository.UpdateUserData")
	}
	//_, err := r.DB.Exec(query, user.Avatar, user.Name, user.Sirname, user.Email)
	return nil
}
