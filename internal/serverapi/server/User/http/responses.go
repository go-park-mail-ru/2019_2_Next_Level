package http

import "2019_2_Next_Level/internal/model"

type GetUserProfile struct {
	Status string `json:"status"`
	Answer GetUserProfileAnswer `json:"userInfo"`
}

type GetUserProfileAnswer struct {
	Name      string `json:"firstName"`
	Sirname   string `json:"secondName"`
	BirthDate string `json:"birthDate"`
	Sex       string `json:"sex"`
	Email     string `json:"login"`
	Avatar    string `json:"avatar"`
	Login 	  string `json:"nickName"`
	Folders   []model.Folder `json:"folders"`
}
