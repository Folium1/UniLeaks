package testing

import (
	"context"
	errHandler "leaks/pkg/err"
	"leaks/pkg/models"
	repository "leaks/pkg/user/repository/mysql"
	user "leaks/pkg/user/service"
	"os"
	"testing"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
			got, err := userService.CreateUser(ctx, user.data)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
			testName: "Get existed user", args: args{ctx, 1}, wantErr: false,
		},
		{
			testName: "Get user that doesn't exist", args: args{ctx, 50000}, wantErr: true,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
			err := userService.IsBanned(ctx, tt.email)
			if err != nil && tt.want != nil && tt.want != err {
				t.Error("Err:", err)
			}
		})
	}
}
