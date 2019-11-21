package Auth

import (
	pb "2019_2_Next_Level/internal/Auth/service"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"context"
	"google.golang.org/grpc"
)

type IAuthClient interface {
	Init(string, string) error
	Destroy()
	LoginBySession(string) (string, int32)
	StartSession(string) (string, int32)
	DestroySession(string) (int32)
	DestroyUserSessions(string) (int32)
	ChangePassword(login string, oldPass string, newPass string) (int32)
	CheckCredentials(login string, password string) (int32)
	GetError(int32) error
}

type AuthClient struct {
	connection *grpc.ClientConn
	client pb.AuthClient
}

func (c *AuthClient) Init(host, port string) error{
	var err error
	c.connection, err = grpc.Dial(host+port, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.client = pb.NewAuthClient(c.connection)
	return nil
}

func (c *AuthClient) Destroy() {
	c.connection.Close()
}
func (c *AuthClient) LoginBySession(session string) (string, int32) {
	message := &pb.String{Data:session}
	res, _ := c.client.LoginBySession(context.Background(), message)
	return res.Result, res.Code
}

func (c *AuthClient) StartSession(data string) (string, int32) {
	message := &pb.String{Data:data}
	res, _ := c.client.StartSession(context.Background(), message)
	return res.Result, res.Code
}

func (c *AuthClient) DestroySession(data string) (int32) {
	message := &pb.String{Data:data}
	res, _ := c.client.DestroySession(context.Background(), message)
	return res.Code
}

func (c *AuthClient) DestroyUserSessions(login string) (int32) {
	message := &pb.String{Data:login}
	res, _ := c.client.DestroyUserSessions(context.Background(), message)
	return res.Code
}

func (c *AuthClient) ChangePassword(login string, oldPass string, newPass string) (int32) {
	message := &pb.ChangePasswordMessage{Login:login, OldPass:oldPass, NewPass:newPass}
	res, _ := c.client.ChangePassword(context.Background(), message)
	return res.Code
}

func (c *AuthClient) CheckCredentials(login string, password string) (int32) {
	message := &pb.CredentialsMessage{Login:login, Password:password}
	res, _ := c.client.CheckCredentials(context.Background(), message)
	return res.Code
}

func (c *AuthClient) GetError(code int32) error {
	var err error
	if code != e.OK{
		err = e.Error{}.SetCode(int(code))
	}
	return err
}