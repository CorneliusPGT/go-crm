package main

import (
	"context"
	"crm/internal/handler"
	md "crm/internal/middleware"
	migrator "crm/internal/migrations"
	repo "crm/internal/repository"
	"crm/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
"github.com/go-chi/cors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/crm?sslmode=disable"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET не задан")
	}

	if err := migrator.Run(dsn, "./internal/migrations"); err != nil {
		log.Fatalf("Миграции провалились: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("pgxpool не запустился: %v", err)
	}

	authRepo := repo.NewAuthRepo(pool)
	taskRepo := repo.NewTaskRepo(pool)

	authService := service.NewAuthService(authRepo, secret)
	taskService := service.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewATaskHandler(taskService)

	fmt.Println(taskHandler)

	defer pool.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register())
		r.Post("/login", authHandler.Login())
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(md.JWTMiddleware(secret))
		r.Use(md.AdminMiddleware)
		r.Post("/task", taskHandler.CreateTask())
		r.Patch("/assign", taskHandler.AssignTask())
		r.Delete("/task/{id}", taskHandler.DeleteTask())
		r.Get("/users", taskHandler.GetUsers())
	})

	r.Group(func(r chi.Router) {
		r.Use(md.JWTMiddleware(secret))
		r.Get("/tasks", taskHandler.GetList())
		r.Patch("/task/{id}", taskHandler.UpdateStatus())
	})

	go func() {
		log.Println("Сервер запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Сервер упал: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Остановка сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера")
	}
	log.Println("Сервер остановлен")
}
