package handler

import (
	"crm/internal/middleware"
	"crm/internal/model"
	"crm/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type TaskHandler struct {
	service service.TaskService
}

func NewATaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

func (h *TaskHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Не авторизирован", http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "Запрещено создавать задачи", http.StatusForbidden)
			return
		}
		users, err := h.service.GetUsers(r.Context())
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
		if users == nil {
			users = []*model.User{}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
			return
		}
	}
}

func (h *TaskHandler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Не авторизирован", http.StatusUnauthorized)
			return
		}
		var req model.Task
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		if req.Title == "" || req.Description == "" || req.Status == "" {
			http.Error(w, "title, description и status обязательны", http.StatusBadRequest)
			return
		}
		claimID := claims.UserID
		req.CreatedBy = claimID
		req.CreatedAt = time.Now()

		if claims.Role != "admin" {
			http.Error(w, "Запрещено создавать задачи", http.StatusForbidden)
			return
		}

		err = h.service.CreateTask(r.Context(), &req)
		if err != nil {
			http.Error(w, "Ошибка создания задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *TaskHandler) AssignTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Не авторизирован", http.StatusUnauthorized)
			return
		}
		var req struct {
			TaskID int `json:"taskID"`
			UserID int `json:"userID"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "Запрещено назначать задачи", http.StatusForbidden)
			return
		}

		err = h.service.AssignTask(r.Context(), req.TaskID, req.UserID)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *TaskHandler) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id", http.StatusBadRequest)
			return
		}
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Не авторизирован", http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "Запрещено удалять задачи", http.StatusForbidden)
			return
		}

		err = h.service.DeleteTask(r.Context(), id)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *TaskHandler) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.service.GetList(r.Context())
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
		if tasks == nil {
			tasks = []*model.Task{}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
			return
		}

	}
}

func (h *TaskHandler) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			http.Error(w, "Не авторизирован", http.StatusUnauthorized)
			return
		}
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id", http.StatusBadRequest)
			return
		}
		var req struct {
			Status string `json:"status"`
		}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		updatedTask, err := h.service.UpdateStatus(r.Context(), id, claims.UserID, req.Status, claims.Role)
		if err != nil {
			if err.Error() == "Запрещено" {
				http.Error(w, "Запрещено", http.StatusForbidden)
				return
			}
			http.Error(w, "Ошибка сервиса: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTask)

	}
}
