package repository

import (
	"os"
	"testing"

	"leaks/pkg/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var (
	r = New()
)

func TestRepository_CreateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		newUser models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"CreateUser ", args{context.Background(), models.User{NickName: "test", Email: "test" + os.Getenv("EMAIL_DOMAIN"), Password: "test"}}, false},
		{"CreateUser with used nickname and mail", args{context.Background(), models.User{NickName: "test", Email: "test" + os.Getenv("EMAIL_DOMAIN"), Password: "test"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.CreateUser(tt.args.ctx, tt.args.newUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_UserById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"UserById", args{context.Background(), 1}, false},
		{"UserById", args{context.Background(), 231232131231242341}, true},
	}
	for _, tt := range tests {
		_, err := r.UserById(tt.args.ctx, tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("Repository.UserById() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestRepository_BannedMails(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}
			_, err := r.BannedMailHashes(tt.args.ctx)
			if err != nil {
				t.Errorf("Repository.BannedMails() error = %v, wantErr", err)
				return
			}
		})
	}
}
