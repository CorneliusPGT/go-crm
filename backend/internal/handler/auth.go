package handler

import (
	"crm/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Почта и пароль обязательны", http.StatusBadRequest)
			return
		}

		retUser, err := h.service.Register(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]string{"token": retUser.Token, "role": retUser.Role}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
			return
		}
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Почта и пароль обязательны", http.StatusBadRequest)
			return
		}

		retUser, err := h.service.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if err.Error() == "Неверные данные" {
				http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{"token": retUser.Token, "role": retUser.Role}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
			return
		}

	}
}
