package repository

import (
	"2019_2_Next_Level/internal/model"
	"fmt"
)

type Postgres struct {
	conn model.Connection
}

func (d *Postgres) UserExists(username string) bool {
	return true
}

func (d *Postgres) AddEmail(email model.Email) error {
	fmt.Println(email.Stringify())
	return nil
}

func NewRepository(conn *model.Connection) Postgres {
	return Postgres{conn: *conn}
}
