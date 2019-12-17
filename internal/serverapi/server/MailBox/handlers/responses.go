package handlers

import "2019_2_Next_Level/internal/serverapi/server/MailBox/models"

//easyjson:json
type GetMessagesList struct{
	Status string `json:"status"`
	Messages []int64 `json:"messages"`
}

//easyjson:json
type GetFolderMessagesResponse struct {
	Status string `json:"status"`
	PagesCount int `json:"pagesCount"`
	Page int `json:"page"`
	Messages []models.MailToGet `json:"messages"`
}

//easyjson:json
type GetMessageResponse struct {
	Status  string           `json:"status"`
	Message models.MailToGet `json:"message"`
}

type GetMessagesListResponse struct {
	Status  string           `json:"status"`
	Messages []models.MailToGet `json:"messages"`
}

//easyjson:json
type GetMessagesCountResponse struct {
	Status string
	Count int
}