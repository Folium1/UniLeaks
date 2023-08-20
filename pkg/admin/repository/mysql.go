package repository

import (
	"context"

	"leaks/pkg/config"
	"leaks/pkg/models"

	"database/sql"

	_ "github.com/lib/pq"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository() *UserRepo {
	db := config.PostgreSQLConn()
	return &UserRepo{db}
}

func (r *UserRepo) AllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT * FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.NickName, &user.IsBanned)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.NickName, &user.IsBanned)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetByNick(ctx context.Context, nick string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE nick_name = $1"
	row := r.db.QueryRow(query, nick)
	err := row.Scan(&user.ID, &user.NickName, &user.IsBanned)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepo) BannedUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT * FROM users WHERE is_banned = true"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.NickName, &user.IsBanned)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) BanUser(ctx context.Context, id string) error {
	query := "UPDATE users SET is_banned = true WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UnbanUser(ctx context.Context, id string) error {
	query := "UPDATE users SET is_banned = false WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
