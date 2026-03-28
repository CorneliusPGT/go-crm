package service

import (
	"context"
	"crm/internal/model"
	repo "crm/internal/repository"
	"errors"
)

type TaskService interface {
	GetUsers(ctx context.Context) ([]*model.User, error)
	CreateTask(ctx context.Context, task *model.Task) error
	GetList(ctx context.Context) ([]*model.Task, error)
	DeleteTask(ctx context.Context, id int) error
	AssignTask(ctx context.Context, taskID, userID int) error
	UpdateStatus(ctx context.Context, taskID, userID int, status, role string) (*model.Task, error)
}

type taskService struct {
	repo repo.TaskRepo
}

func (t *taskService) GetUsers(ctx context.Context) ([]*model.User, error) {
	return t.repo.GetUsers(ctx)
}

func (t *taskService) AssignTask(ctx context.Context, taskID int, userID int) error {
	return t.repo.AssignTask(ctx, taskID, userID)
}

func (t *taskService) CreateTask(ctx context.Context, task *model.Task) error {
	return t.repo.CreateTask(ctx, task)
}

func (t *taskService) DeleteTask(ctx context.Context, id int) error {
	return t.repo.DeleteTask(ctx, id)
}

func (t *taskService) GetList(ctx context.Context) ([]*model.Task, error) {
	return t.repo.GetList(ctx)
}

func (t *taskService) UpdateStatus(ctx context.Context, taskID, userID int, status, role string) (*model.Task, error) {
	task, err := t.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if role != "admin" && userID != task.CreatedBy && (task.AssignedTo == nil || userID != *task.AssignedTo) {
		return nil, errors.New("Запрещено")
	}
	return t.repo.UpdateStatus(ctx, taskID, status)
}

func NewTaskService(repo repo.TaskRepo) *taskService {
	return &taskService{
		repo: repo,
	}
}
