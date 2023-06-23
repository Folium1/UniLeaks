package testing

import (
	"context"
	errHandler "leaks/err"
	"leaks/models"
	repository "leaks/user/repository/mysql"
	user "leaks/user/service"
	"os"
	"testing"
)

var (
	repo        = repository.New()
	userService = user.New(repo)
)

func TestUserUseCase(t *testing.T) {
	t.Run("CreateUser", createUserTest)
	t.Run("GetById", getByIdTest)
	t.Run("IsBanned", isBannedTest)
}

func createUserTest(t *testing.T) {
	testUsers := []struct {
		data     models.User
		testName string
		wantErr  bool
	}{
		{
			models.User{
				NickName: RandomString(8),
				Email:    RandomString(10) + os.Getenv("MAIL_DOMAIN"),
				Password: "mockpass",
			},
			"Valid User register",
			false,
		},
	}
	for _, user := range testUsers {
		t.Run(user.testName, func(t *testing.T) {
			got, err := userService.CreateUser(context.Background(), user.data)
			if err != nil && !errHandler.IsDuplicateEntryError(err) {
				t.Error("Couldn't create user,err: ", err)
			}
			if got == 0 {
				t.Error("Couldn't create user, zero value return")
			}
		})
	}
}

func getByIdTest(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		testName string
		args     args
		wantErr  bool
	}{
		{
			testName: "Get existed user", args: args{context.Background(), 1}, wantErr: false,
		},
		{
			testName: "Get user that doesn't exist", args: args{context.Background(), 50000}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := userService.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) && !tt.wantErr && got.ID == 0 {
				t.Errorf("UserUseCase.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func isBannedTest(t *testing.T) {
	testData := []struct {
		testName string
		email    string
		want     error
	}{
		{"Existing banned mock user", "2", nil},
		{"User that doesn't exist", "212312312321", errHandler.UserIsBannedErr},
	}
	for _, tt := range testData {
		t.Run(tt.testName, func(t *testing.T) {
			err := userService.IsBanned(context.Background(), tt.email)
			if err != nil && tt.want != nil && tt.want != err {
				t.Error("Err:", err)
			}
		})
	}
}
