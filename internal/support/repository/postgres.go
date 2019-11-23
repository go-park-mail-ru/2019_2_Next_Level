package repository

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/support/models"
	"database/sql"
	"fmt"
)

type PostgresRepository struct {
	model.PostgresRepoTemplate
}

func NewPostgresRepository(DB *sql.DB) *PostgresRepository {
	return &PostgresRepository{}
}


func (r *PostgresRepository) StartChat(user string, theme string, supportNick string) (chatId interface{}, err error) {
	query := `INSERT INTO Chat (userNick, theme, supportNick) VALUES ($1, $2, $3) RETURNING id`
	err = r.DB.QueryRow(query, user, theme, supportNick).Scan(&chatId)
	return
}

func (r *PostgresRepository) GetChat(user string, id interface{}) (models.Chat, error) {
	queryChat := `SELECT id, theme, isOpen, startDate FROM Chat WHERE id=$1`
	queryMessages := `SELECT id, sent, wasRead, author, body FROM ChatMessage WHERE chatId=$1`

	//var date string
	var chat models.Chat
	err := r.DB.QueryRow(queryChat, id).Scan(&chat.Id, &chat.Theme, &chat.IsOpen, &chat.StartDate)
	if err != nil {
		return chat, err
	}
	rows, err := r.DB.Query(queryMessages, id)
	if err != nil {
		return chat, err
	}
	messages := make([]models.Message, 0)
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Id, &msg.Sent, &msg.WasRead, &msg.UserNick, &msg.Body)
		if err != nil {
			return chat, err
		}
		messages = append(messages, msg)
	}
	chat.Messages = messages
	return chat, nil
}
func (r *PostgresRepository) GetChats(user string) ([]models.Chat, error) {
	query := `SELECT id, theme, isopen, startDate FROM Chat WHERE %s=$1`
	query0 := `SELECT get_chat_field($1)`
	var role string
	err := r.DB.QueryRow(query0, user).Scan(&role)
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf(query, role)
	rows, err := r.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	chats := make([]models.Chat, 0)
	for rows.Next(){
		var chat models.Chat
		err := rows.Scan(&chat.Id, &chat.Theme, &chat.IsOpen, &chat.StartDate)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *PostgresRepository) CloseChat(user string, chatId interface{}) error {
	query := `UPDATE Chat SET isOpen=false WHERE (userNick=$2 OR supportNick=$2) AND id=$1`
	_, err := r.DB.Exec(query, chatId, user)
	return err
}
func (r *PostgresRepository) SendMessage(user string, chatId interface{}, message models.Message) error {
	query := `INSERT INTO ChatMessage (chatId, body, author) VALUES ($1, $2, $3)  `
	_, err := r.DB.Exec(query, chatId, message.Body, user)
	return err
}
func (r *PostgresRepository) GetMessage(user string, chatId interface{}, messageId interface{}) (models.Message, error) {
	query := `SELECT id, chatId, sent, wasRead, body, author FROM ChatMessage WHERE id=$1`
	var msg models.Message
	err := r.DB.QueryRow(query, messageId).Scan(&msg.Id, &msg.ChatId, &msg.Sent, &msg.WasRead, &msg.Body, &msg.UserNick)
	return msg, err
}
