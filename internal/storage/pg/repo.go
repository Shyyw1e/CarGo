package pg

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID		 	 string
	Email		 string
	PasswordHash string
	Role		 string
	CreatedAt	 time.Time
	UpdatedAt	 time.Time
}

type Repo struct {Conn *pgx.Conn}

func New(ctx context.Context, dsn string) (*Repo, error) {
	c, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Repo{Conn: c}, nil
}

func (r *Repo) Close(ctx context.Context) (error) {
	return r.Conn.Close(ctx)
}

func (r *Repo) CreateUser(ctx context.Context, email, hash string) (User, error) {
	var u User
	err := r.Conn.QueryRow(ctx,
		`INSERT INTO users(email,password_hash) VALUES($1, $2)
		 RETURNING id,email,password_hash,role,created_at,updated_at`, email, hash,
		 ).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.Conn.QueryRow(ctx,
		`SELECT id,email,password_hash,role,created_at,updated_at FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return User{}, pgx.ErrNoRows }
		return User{}, err
	}
	return u, nil

}

func (r *Repo) GetUserById(ctx context.Context, id string) (User, error) {
	var u User
	err := r.Conn.QueryRow(ctx,
		`SELECT id,email,password_hash,role,created_at,updated_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err

}