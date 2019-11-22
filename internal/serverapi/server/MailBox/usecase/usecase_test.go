package usecase

import (
	"2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/repository"
	"testing"
)
func init() {
	defaultConf := config.Database{DBName: "nextlevel", Port: "5432", Host: "localhost", User: "postgres", Password: "postgres"}
	config.Conf.DB = defaultConf
}
func TestMailBoxUsecase_GetMailList(t *testing.T) {
	repo, err := repository.GetPostgres()
	if err != nil {
		return
	}
	_, err = repo.GetEmailList("aaa@nlmail.ddns.net", "incoming", "", 1, 100)
	if err != nil {
		return
	}
}

