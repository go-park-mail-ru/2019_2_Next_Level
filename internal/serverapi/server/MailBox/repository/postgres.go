package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	e "2019_2_Next_Level/pkg/Error"
	"2019_2_Next_Level/pkg/sqlTools"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"strconv"
	"strings"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

func GetPostgres() (PostgresRepository, error) {
	r := PostgresRepository{}
	if r.DB == nil {
		//err := r.Init()
		//if err != nil {
		//	return r, err
		//}
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

func (r *PostgresRepository) GetEmailByCode(login string, code interface{}) ([]model.Email, error) {
	query := queryGetEmailByCode

	mails := make([]model.Email, 0)
	ids := make([]int64, 0)
	if codes, ok := code.([]string); ok {
		for _, elem := range codes {
			id, _ := strconv.ParseInt(elem, 10, 64)
			ids = append(ids, id)
		}
	} else {
		id, _ := strconv.ParseInt(code.(string), 10, 64)
		ids = append(ids, id)
	}
	for _, id := range ids {
		var mail model.Email
		var when string
		err := r.DB.QueryRow(query, id).Scan(&mail.From, &mail.To, &when, &mail.Header.Subject, &mail.Body, &mail.IsRead, &mail.Direction)
		if err != nil {
			return mails, e.Error{}.SetError(err)
		}
		mail.Id=int(id)
		mail.Header.WhenReceived, _= time.Parse(time.RFC3339, when)
		mails = append(mails, mail)
	}
	return mails, nil
}

func (r *PostgresRepository) GetEmailList(login string, folder string, sort interface{}, since int64, count int) ([]model.Email, error) {
	query := queryGetEmailList

	row, err := r.DB.Query(query, login, folder, count, since)
	list := make([]model.Email, 0)
	if err != nil {
		return list, e.Error{}.SetCode(e.NotExists).SetError(err)
	}

	for row.Next() {
		mail := model.Email{}
		var when string

		if err := row.Scan(&mail.Id, &mail.From, &mail.To, &when, &mail.Header.Subject, &mail.Body, &mail.IsRead, &mail.Direction); err != nil {
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
		err = saveMessage.QueryRow(from, email.Header.Subject, email.Body, "out", "sent", from).Scan(&id)
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

func (r *PostgresRepository) AddFolder(login string, foldername string) error {
	query := `INSERT INTO Folder (name, owner) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, foldername, login)
	return err
}
func (r *PostgresRepository) ChangeMailFolder(login string, foldername string, mailid []models.MailID) error {
	query := `UPDATE Message SET folder=$1 WHERE id IN (%s)`
	var args []string
	for i:=0; i<len(mailid); i++ {
		args = append(args, `$`+strconv.Itoa(i+2))
	}
	query = fmt.Sprintf(query, strings.Join(args, ", "))
	params := make([]interface{}, 0, len(mailid)+1)
	params = append(params, foldername)
	for _, val := range mailid{
		params = append(params, val)
	}
	_, err := r.DB.Exec(query, params...)
	return err
}

func (r *PostgresRepository) DeleteFolder(login string, folderName string) error {
	query := `DELETE FROM Folder WHERE owner=$1 AND name=$2`
	_, err := r.DB.Exec(query, login, folderName)
	return err;
}

func (r *PostgresRepository) FindMessages(login, request string) ([]int64, error) {
	query := `WITH
			own_messages AS (SELECT id, concat(subject, ' ', body) as data FROM Message
					WHERE owner=$1),
			fulltest_search AS (SELECT id, data FROM own_messages
					WHERE to_tsvector("data") @@ plainto_tsquery($2)),
			res AS (SELECT id, sum(p) as p FROM
					(SELECT id, p from
						(SELECT id, ts_rank(to_tsvector("data"), plainto_tsquery('t')) as p FROM fulltest_search) as b
						UNION
						(SELECT id, 0 as p FROM own_messages WHERE lower(data) LIKE concat('%', lower($2), '%'))
					) as FF
					GROUP BY id
				)
	SELECT res.id FROM res JOIN Message ON Message.id=res.id ORDER BY res.p desc, res.id;`

	row, err := r.DB.Query(query, login, request)
	list := make([]int64, 0)
	if err != nil {
		return list, e.Error{}.SetCode(e.NotExists).SetError(err)
	}

	for row.Next() {
		var id int64
		if err := row.Scan(&id); err != nil {
			return list, e.Error{}.SetError(err)
		}
		list = append(list, id)
	}
	return list, nil
}


func (r *PostgresRepository) GetUserData(login string) (string, string, error) {
	query := `SELECT firstname, secondname, avatar FROM users WHERE login=$1`
	row := r.DB.QueryRow(query, login)
	if row == nil {
		return "", "", e.Error{}.SetString("Empty row")
	}
	var name1, name2, avatar string
	err := row.Scan(&name1, &name2, &avatar)
	if err != nil {
		return "", "", err
	}
	return name1 + " " + name2, avatar, nil
}