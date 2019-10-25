package repository

import "2019_2_Next_Level/internal/model"

type Postgres struct {
	conn model.Connection
}

func (d *Postgres) SetConnection(c *model.Connection) error {
	return nil
}
func (d *Postgres) RegisterNewSession(model.User) error {
	return nil
}
func (d *Postgres) CheckSession(model.UUID) error {
	return nil
}
func (d *Postgres) DiscardSession(string) error {
	return nil
}

func NewPostgres(c *model.Connection) *Postgres {
	return &Postgres{conn: *c}
}
