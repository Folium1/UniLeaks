package repository

import (
	"fmt"
	"leaks/pkg/config"
	logg "leaks/pkg/logger"
	"leaks/pkg/models"
	"database/sql"

	_ "github.com/lib/pq"
)

var logger = logg.NewLogger()

type Repository struct {
	db *sql.DB
}

func New() *Repository {
	db := config.PostgreSQLConn()
	return &Repository{db}
}

func (r *Repository) CreateUser(newUser models.User) (int, error) {
	rows, err := r.db.Query(`INSERT INTO users (nick_name, email, password) VALUES ($1, $2, $3)`, newUser.NickName, newUser.Email, newUser.Password)
	for rows.Next() {
		if rows.Err() != nil {
			return -1, err
		}
	}
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return -1, err
	}
	return newUser.ID, nil
}

func (r *Repository) UserById(id int) (models.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}
	dbUser := models.User{ID: id}
	for rows.Next() {
		err = rows.Scan(dbUser)
		if err != nil {
			return models.User{}, err
		}
	}
	return dbUser, nil
}

func (r *Repository) UserByNick(nick string) (models.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE nick_name = $1", nick)
	if err != nil {
		return models.User{}, err
	}
	dbUser := models.User{NickName: nick}
	for rows.Next() {
		err = rows.Scan(dbUser)
		if err != nil {
			return models.User{}, err
		}
	}
	return dbUser, nil
}

func (r *Repository) BannedMailHashes() ([]string, error) {
	rows, err := r.db.Query("SELECT email FROM users WHERE is_banned = $1", true)
	if err != nil {
		return nil, err
	}
	var mails []string
	for rows.Next() {
		var mail string
		err = rows.Scan(mail)
		if err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}
	return mails, nil
}