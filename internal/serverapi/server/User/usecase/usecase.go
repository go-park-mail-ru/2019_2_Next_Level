package usecase

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/model"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	"2019_2_Next_Level/internal/serverapi/server/config"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

func NewUserUsecase(repo user.UserRepository, auth Auth.IAuthClient) UserUsecase {
	usecase := UserUsecase{repo:repo}
	usecase.auth = auth
	return usecase
}

type UserUsecase struct {
	repo user.UserRepository
	auth Auth.IAuthClient
}
func (u *UserUsecase) GetUser(login string) (model.User, error) {
	user, err := u.repo.GetUser(login)
	if err != nil {
		switch err.(type) {
		case e.Error:
			return user, err
		default:
			break
		}
		return user, e.Error{}.SetCode(e.ProcessError)
	}
	user.Login = user.Email
	//user.Avatar = config.Conf.HttpConfig.SelfURL + "avatar/"+user.Avatar
	//user.Avatar = "/static/images/icon/no-avatar.svg"
	if user.Avatar=="" {
		user.Avatar = config.Conf.DefaultAvatar
	}
	a := config.Conf
	fmt.Println(a)
	user.Avatar = config.Conf.Protocol+"://" + config.Conf.HostName+"/"+config.Conf.StaticDir+"/"+config.Conf.AvatarDir + "/" + user.Avatar;
	user.Sanitize()
	return user, nil
}

func (u *UserUsecase) GetUserFolders(login string) ([]model.Folder, error) {
	folders, err := u.repo.GetUserFolders(login)
	if err != nil {
		switch err.(type) {
		case e.Error:
			return folders, err
		default:
			break
		}
		return folders, e.Error{}.SetCode(e.ProcessError)
	}
	for i := range folders {
		folders[i].Sanitize()
	}
	return folders, nil
}

func (u *UserUsecase) EditAvatar(login string, file multipart.File, header *multipart.FileHeader) (string, error) {
	path := config.Conf.RootDir + "/" + config.Conf.StaticDir
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	path += config.Conf.AvatarDir+ "/"

	keyToHash := login + ":" + time.Now().String()
	filename := hex.EncodeToString(md5.New().Sum([]byte(keyToHash)))
	temp := strings.Split(header.Filename, ".")
	var typeFile string
	if len(temp) < 0 {
		typeFile = ""
	} else {
		typeFile = "." + temp[len(temp)-1]
	}
	filename += typeFile

	path +=  filename
	fmt.Println(file)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		//log.Log().E(fmt.Sprintf("Cannot create avatar file for %s with error: %v", login, err));
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	fmt.Println(err)
	return filename, err
}

func (u *UserUsecase) EditUser(user *model.User) error {
	user.Password = ""
	err := u.repo.UpdateUserData(user)
	if err != nil {
		switch err.(type) {
		case e.Error:
			return err
		default:
			break
		}
		return e.Error{}.SetCode(e.InvalidParams)
	}

	return nil
}
func (u *UserUsecase) EditPassword(login string, oldPass string, newPass string) error {
	return u.auth.GetError(u.auth.ChangePassword(login, oldPass, newPass))
}
