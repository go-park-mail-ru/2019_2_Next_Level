package Auth

import (
	pb "2019_2_Next_Level/internal/Auth/service"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type AuthServer struct {
	repo Repository
}

func NewAuthServer(repo Repository) *AuthServer {
	return &AuthServer{repo: repo}
}

func (s *AuthServer) LoginBySession(ctx context.Context, sessionToken *pb.String) (*pb.StringResult, error) {
	login, err := s.repo.GetLoginBySession(sessionToken.Data)
	if err != nil {
		return &pb.StringResult{Result:"", Code:e.NotExists}, nil
	}
	return &pb.StringResult{Result:login, Code:e.OK}, nil
}

func (s *AuthServer) StartSession(ctx context.Context, login *pb.String) (*pb.StringResult, error) {
	token, _ := uuid.NewUUID()

	err := s.repo.AddNewSession(login.Data, token.String())
	if err != nil {
		return &pb.StringResult{Result:"", Code: int32(err.(e.Error).Code)}, nil
	}

	return &pb.StringResult{Result:token.String(), Code:e.OK}, nil
}

func (s *AuthServer) DestroySession(ctx context.Context, sessionToken *pb.String) (*pb.StatusResult, error) {
	err := s.repo.DeleteSession(sessionToken.Data)
	if err != nil {
		return &pb.StatusResult{Code:int32(err.(e.Error).Code)}, nil
	}
	return &pb.StatusResult{Code: int32(e.OK)}, nil
}

func (s *AuthServer) DestroyUserSessions(ctx context.Context, login *pb.String) (*pb.StatusResult, error) {
	err := s.repo.DeleteUserSessions(login.Data)
	if err != nil {
		return &pb.StatusResult{Code: int32(err.(e.Error).Code)}, nil
	}
	return &pb.StatusResult{Code: e.OK}, nil
}

func (s *AuthServer) ChangePassword(ctx context.Context, data *pb.ChangePasswordMessage) (*pb.StatusResult, error) {
	err := s.checkPassword(data.Login, data.OldPass)
	if err != nil {
		return &pb.StatusResult{Code:int32(err.(e.Error).Code)}, nil
	}

	sault := s.getSault(data.Login)
	newPassHash := PasswordPBKDF2([]byte(data.NewPass), []byte(sault))

	err = s.repo.UpdateUserPassword(data.Login, newPassHash, sault)
	return &pb.StatusResult{Code:e.OK}, nil
}

func (s *AuthServer) CheckCredentials(ctx context.Context, data *pb.CredentialsMessage) (*pb.StatusResult, error) {
	err := s.checkPassword(data.Login, data.Password)
	if err != nil {
		return &pb.StatusResult{Code:int32(err.(e.Error).Code)}, nil
	}
	return &pb.StatusResult{Code: e.OK}, nil
}

func (s *AuthServer) RegisterUser(ctx context.Context, data *pb.CredentialsMessage) (*pb.StatusResult, error) {
	sault := s.getSault(data.Login)
	newPassHash := PasswordPBKDF2([]byte(data.Password), []byte(sault))

	err := s.repo.UpdateUserPassword(data.Login, newPassHash, sault)
	fmt.Println(err)
	return &pb.StatusResult{Code:e.OK}, nil
}

func (s *AuthServer) checkPassword(login string, password string) error {
	currPass, sault, err := s.repo.GetUserCredentials(login)
	if err != nil {
		return e.Error{}.SetCode(e.NotExists)
	}
	if !CheckPassword([]byte(password), []byte(currPass), []byte(sault)) {
		return e.Error{}.SetCode(e.WrongPassword)
	}
	return nil
}

func (s *AuthServer) getSault(login string) []byte {
	return GenSault(login)
}