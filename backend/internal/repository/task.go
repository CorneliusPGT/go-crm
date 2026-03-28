package repo

import (
	"context"
	"crm/internal/model"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo interface {
	GetUsers(ctx context.Context) ([]*model.User, error)
	CreateTask(ctx context.Context, task *model.Task) error
	GetList(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id int) (*model.Task, error)
	DeleteTask(ctx context.Context, id int) error
	AssignTask(ctx context.Context, taskID, userID int) error
	UpdateStatus(ctx context.Context, taskID int, status string) (*model.Task, error)
}

type taskRepo struct {
	pool *pgxpool.Pool
}

func (t *taskRepo) GetUsers(ctx context.Context) ([]*model.User, error) {
	rows, err := t.pool.Query(ctx, `
	SELECT id FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User

	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (t *taskRepo) AssignTask(ctx context.Context, taskID int, userID int) error {
	_, err := t.pool.Exec(ctx, `
	UPDATE tasks SET assigned_to = $1 WHERE id = $2
	`, userID, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskRepo) CreateTask(ctx context.Context, task *model.Task) error {
	_, err := t.pool.Exec(ctx, `
		INSERT INTO tasks (title, description, status, assigned_to, created_by, created_at) VALUES (
		$1, $2, $3, $4, $5, $6
		)
	`, task.Title, task.Description, task.Status, task.AssignedTo, task.CreatedBy, task.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskRepo) DeleteTask(ctx context.Context, id int) error {
	res, err := t.pool.Exec(ctx, `
	DELETE FROM tasks WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	rowsAff := res.RowsAffected()
	if rowsAff == 0 {
		return errors.New("Не найдено")
	}
	return nil
}

func (t *taskRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
	var task model.Task
	row := t.pool.QueryRow(ctx,
		`SELECT id, title, description, status, assigned_to, created_by, created_at FROM tasks WHERE id = $1
	`, id)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Не найдено")
		}
		return nil, err
	}

	return &task, nil

}

func (t *taskRepo) GetList(ctx context.Context) ([]*model.Task, error) {
	rows, err := t.pool.Query(ctx,
		`
	SELECT id, title, description, status, assigned_to, created_by, created_at FROM tasks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task

	for rows.Next() {
		var t model.Task

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.AssignedTo,
			&t.CreatedBy,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (t *taskRepo) UpdateStatus(ctx context.Context, taskID int, status string) (*model.Task, error) {
	var task model.Task
	err := t.pool.QueryRow(ctx, `
		UPDATE tasks
		SET status = $1
		WHERE id = $2
		RETURNING id, title, description, status, assigned_to, created_by, created_at
	`, status, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Не найдено")
		}
		return nil, err
	}
	return &task, nil
}

func NewTaskRepo(pool *pgxpool.Pool) *taskRepo {
	return &taskRepo{
		pool: pool,
	}
}
