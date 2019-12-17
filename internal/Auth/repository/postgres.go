package repository

import (
	"2019_2_Next_Level/internal/model"
	e "2019_2_Next_Level/pkg/Error"
	"fmt"
)
const (
	queryGetLoginBySession = `SELECT login FROM %s WHERE token=$1`
)

type PostgresRepo struct {
	model.PostgresRepoTemplate
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{}
}

func (r *PostgresRepo) GetLoginBySession(uuid string) (string, error) {
	query := fmt.Sprintf(queryGetLoginBySession, "session")
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

func (r *PostgresRepo) AddNewSession(login string, uuid string) error {
	if ok := r.checkUserExist(login); !ok {
		return e.Error{}.SetCode(e.NotExists)
	}
	//clearExistSessionQuery := `DELETE FROM session WHERE login=$1`
	query := fmt.Sprintf("INSERT INTO %s (login, token) VALUES ($1, $2);", "session")

	//r.DB.Exec(clearExistSessionQuery, login)
	_, err := r.DB.Exec(query, login, uuid)
	return err
}

func (r *PostgresRepo) DeleteSession(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token=$1", "Session")
	_, err := r.DB.Exec(query, token)
	if err != nil {
		return e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return nil
}

func (r *PostgresRepo) DeleteUserSessions(login string) error {
	query := `DELETE FROM Session WHERE login=$1`
	_, err := r.DB.Exec(query, login)
	if err != nil {
		return e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return nil
}

func (r *PostgresRepo) GetUserCredentials(login string) ([]byte, []byte, error) {
	query := "SELECT password, sault FROM users WHERE login=$1;"
	row := r.DB.QueryRow(query, login)

	var pass, sault []byte
	err := row.Scan(&pass, &sault)
	if err != nil {
		return nil, nil, e.Error{}.SetCode(e.NotExists).SetError(err)
	}
	return pass, sault, nil
}

func (r *PostgresRepo) UpdateUserPassword(login string, newPassword []byte, sault []byte) error {
	query := `UPDATE users SET password=$1, sault=$2 WHERE login=$3`
	_, err := r.DB.Exec(query, []byte(newPassword), []byte(sault), login)
	return err
}

func (r *PostgresRepo) checkUserExist(login string) bool {
	row := r.DB.QueryRow("SELECT COUNT(login)>0 FROM users WHERE login=$1", login)
	var isExist bool
	err := row.Scan(&isExist)
	if err != nil {
		isExist = false
	}
	return isExist
}
