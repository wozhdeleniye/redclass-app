package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wozhdeleniye/redclass-app/internal/config"
	"github.com/wozhdeleniye/redclass-app/internal/handlers"
	"github.com/wozhdeleniye/redclass-app/internal/middleware"
	"github.com/wozhdeleniye/redclass-app/internal/migrations"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
	redisrepo "github.com/wozhdeleniye/redclass-app/internal/repositories/redis"
	"github.com/wozhdeleniye/redclass-app/internal/services"
	"github.com/wozhdeleniye/redclass-app/pkg/database"
	"github.com/wozhdeleniye/redclass-app/pkg/redis"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к PostgreSQL
	db, err := database.NewPostgresConnection(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	migrator := migrations.NewGormMigrator(db)
	migrator.Migrate()

	// Подключение к Redis
	redisClient, err := redis.NewRedisConnection(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := redisrepo.NewTokenRepository(redisClient)

	// Инициализация сервисов
	authService := services.NewAuthService(userRepo, tokenRepo, cfg.JWT)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Настройка маршрутов
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/refresh", authHandler.Refresh).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/auth").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
