package repo

import (
	"context"
	"crm/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type authRepo struct {
	pool *pgxpool.Pool
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx, `
	SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func NewAuthRepo(pool *pgxpool.Pool) *authRepo {
	return &authRepo{
		pool: pool,
	}
}

func (r *authRepo) CreateUser(ctx context.Context, user *model.User) (int, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
	SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, user.Email).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("Почта %s уже занята", user.Email)
	}
	var id int
	err = r.pool.QueryRow(ctx, `
	INSERT INTO users (email, password_hash, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id
	`, user.Email, user.PasswordHash, user.Role, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
