package Auth

import (
	pb "2019_2_Next_Level/generated/Auth/service"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"2019_2_Next_Level/tests/mock/serverapi/auth"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Params interface{}

type TestStruct struct {
	input []Params
	expected []Params
	mockParams []Params
}

type ExpectedFunc func(params ...Params)

func initTest(t *testing.T) (*auth.MockRepository, *AuthServer) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := auth.NewMockRepository(mockCtrl)
	s := NewAuthServer(mockRepo)
	return mockRepo, s
}

func runTesting (tests []TestStruct, f ExpectedFunc, testFunc func(TestStruct)) {
	for _, test := range tests{
		f(test.mockParams...)
		testFunc(test)
	}
}

func TestAuthServer_LoginBySession(t *testing.T) {
	t.Parallel()
	mockRepo, s := initTest(t)

	tests := []TestStruct{
		{
			input: []Params{&pb.String{Data: "session_token"}},
			expected: []Params{&pb.StringResult{Result: "admin", Code: e.OK}, nil},
			mockParams: []Params{"session_token", "admin", nil},
		},
		{
			input:[]Params{&pb.String{Data: "session_token"}},
			expected: []Params{
				&pb.StringResult{Result: "", Code: e.NotExists}, nil,
			},
			mockParams: []Params{"session_token", "", e.Error{}.SetCode(e.NotExists)},
		},
	}
	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			GetLoginBySession(params[0].(string)).
			Return(params[1].(string), params[2])
	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.LoginBySession(context.Background(), test.input[0].(*pb.String))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans, test.expected[0].(*pb.StringResult)) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestAuthServer_StartSession(t *testing.T) {
	t.Parallel()
	mockRepo, s := initTest(t)

	tests := []TestStruct{
		{
			input: []Params{&pb.String{Data: "admin"}},
			expected: []Params{&pb.StringResult{Result: "session_token", Code: e.OK}, nil},
			mockParams: []Params{"admin", nil},
		},
		{
			input:[]Params{&pb.String{Data: "admin"}},
			expected: []Params{
				&pb.StringResult{Result: "", Code: e.NotExists},
				nil,
			},
			mockParams: []Params{"admin", e.Error{}.SetCode(e.NotExists)},
		},
	}
	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			AddNewSession(params[0].(string), gomock.Any()).
			Return(params[1])
	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.StartSession(context.Background(), test.input[0].(*pb.String))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans.Code, test.expected[0].(*pb.StringResult).Code) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestAuthServer_DestroySession(t *testing.T) {
	t.Parallel()
	mockRepo, s := initTest(t)

	tests := []TestStruct{
		{
			input: []Params{&pb.String{Data: "session_token"}},
			expected: []Params{&pb.StatusResult{Code: e.OK}, nil},
			mockParams: []Params{"session_token", nil},
		},
		{
			input: []Params{&pb.String{Data: "session_token"}},
			expected: []Params{&pb.StatusResult{Code: e.NotExists},
				nil},
			mockParams: []Params{"session_token", e.Error{}.SetCode(e.NotExists)},
		},
	}
	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			DeleteSession(params[0].(string)).
			Return(params[1])
	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.DestroySession(context.Background(), test.input[0].(*pb.String))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans.Code, test.expected[0].(*pb.StatusResult).Code) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestAuthServer_DestroyUserSessions(t *testing.T) {
	t.Parallel()
	mockRepo, s := initTest(t)

	tests := []TestStruct{
		{
			input: []Params{&pb.String{Data: "admin"}},
			expected: []Params{&pb.StatusResult{Code: e.OK}, nil},
			mockParams: []Params{"admin", nil},
		},
		{
			input: []Params{&pb.String{Data: "admin"}},
			expected: []Params{&pb.StatusResult{Code: e.NotExists},
				nil},
			mockParams: []Params{"admin", e.Error{}.SetCode(e.NotExists)},
		},
	}
	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			DeleteUserSessions(params[0].(string)).
			Return(params[1])
	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.DestroyUserSessions(context.Background(), test.input[0].(*pb.String))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans.Code, test.expected[0].(*pb.StatusResult).Code) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestAuthServer_ChangePassword(t *testing.T) {
	t.Parallel()
	login:="admin"
	mockRepo, s := initTest(t)
	sault := GenSault(login)
	oldPassHash := PasswordPBKDF2([]byte("oldpass"), sault)
	newPassHash := PasswordPBKDF2([]byte("newpass"), sault)

	tests := []TestStruct{
		// OK
		{
			input: []Params{&pb.ChangePasswordMessage{Login: "admin", OldPass:"oldpass", NewPass:"newpass"}},
			expected: []Params{&pb.StatusResult{Code: e.OK}, nil},
			mockParams: []Params{
				// GetCredentials
					// Params
				"admin",
					// Returns
				oldPassHash,
				sault,
				nil,
				//UpdateUserPassword
					// Params
				"admin",
				newPassHash,
				sault,
					// Returns
				nil,
			},
		},
		// Wrong old pass
		{
			input: []Params{&pb.ChangePasswordMessage{Login: "admin", OldPass:"wrongOldpass", NewPass:"newpass"}},
			expected: []Params{&pb.StatusResult{Code: e.WrongPassword}, nil},
			mockParams: []Params{
				// GetCredentials
					// Params
				"admin",
					// Returns
				oldPassHash,
				sault,
				nil,
			},
		},
		// User not exists
		{
			input: []Params{&pb.ChangePasswordMessage{Login: "admin", OldPass:"oldpass", NewPass:"newpass"}},
			expected: []Params{&pb.StatusResult{Code: e.NotExists}, nil},
			mockParams: []Params{
				// GetCredentials
				// Params
				"admin",
				// Returns
				[]byte{},
				[]byte{},
				e.Error{}.SetCode(e.NotExists),
			},
		},
	}
	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			GetUserCredentials(params[0].(string)).
			Return(params[1].([]byte), params[2].([]byte), params[3])
		if len(params) > 4 {
			mockRepo.EXPECT().
				UpdateUserPassword(params[4].(string), params[5].([]byte), params[6].([]byte)).
				Return(params[7])
		}
	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.ChangePassword(context.Background(), test.input[0].(*pb.ChangePasswordMessage))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans.Code, test.expected[0].(*pb.StatusResult).Code) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestAuthServer_CheckCredentials(t *testing.T) {
	t.Parallel()
	login:="admin"
	mockRepo, s := initTest(t)
	sault := GenSault(login)
	pass := "password"
	passHash := PasswordPBKDF2([]byte(pass), sault)

	tests := []TestStruct{
		// OK
		{
			input:    []Params{&pb.CredentialsMessage{Login: "admin", Password: pass}},
			expected: []Params{&pb.StatusResult{Code: e.OK}, nil},
			mockParams: []Params{
				// GetCredentials
				// Params
				"admin",
				// Returns
				passHash,
				sault,
				nil,
			},
		},
		// Wrong pass
		{
			input:    []Params{&pb.CredentialsMessage{Login: "admin", Password: "wrongPassword"}},
			expected: []Params{&pb.StatusResult{Code: e.NotExists}, nil},
			mockParams: []Params{
				// GetCredentials
				// Params
				"admin",
				// Returns
				[]byte{},
				[]byte{},
				e.Error{}.SetCode(e.NotExists),
			},
		},
	}

	expectedFunc := func(params ...Params) {
		mockRepo.EXPECT().
			GetUserCredentials(params[0].(string)).
			Return(params[1].([]byte), params[2].([]byte), params[3])

	}

	runTesting(tests, expectedFunc, func(test TestStruct) {
		ans, err := s.CheckCredentials(context.Background(), test.input[0].(*pb.CredentialsMessage))
		if !cmp.Equal(err, test.expected[1]) {
			t.Errorf("Wrong error: got %v instead %v", err, test.expected[1])
		}
		if !cmp.Equal(ans.Code, test.expected[0].(*pb.StatusResult).Code) {
			t.Errorf("Wrong answer got")
		}
	})
}

func TestClient(t *testing.T) {
	client := AuthClient{}
	client.Init("df", "12")
	client.Destroy()
	//client.LoginBySession("123")
	//recover()
}