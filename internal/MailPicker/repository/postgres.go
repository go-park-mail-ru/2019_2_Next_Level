package repository

import (
	"2019_2_Next_Level/internal/model"
	"fmt"
)

type PostgresRepository struct {

}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) UserExists(login string) bool {
	fmt.Println("User exist", login)
	return true
}
func (r *PostgresRepository) AddEmail(email *model.Email) error {
	fmt.Println("Body", email.Body)
	return nil
}