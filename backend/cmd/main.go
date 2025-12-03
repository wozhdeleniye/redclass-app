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
	cfg := config.Load()

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

	userRepo := postgres.NewUserRepository(db)
	subjectRepo := postgres.NewSubjectRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	tokenRepo := redisrepo.NewTokenRepository(redisClient)

	authService := services.NewAuthService(userRepo, tokenRepo, cfg.JWT)
	subjectService := services.NewSubjectService(subjectRepo, roleRepo, userRepo)
	roleService := services.NewRoleService(roleRepo)
	taskService := services.NewTaskService(taskRepo, roleRepo, subjectRepo)
	projectService := services.NewProjectService(projectRepo, taskRepo, roleRepo)

	authHandler := handlers.NewAuthHandler(authService)
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	roleHandler := handlers.NewRoleHandler(roleService)
	taskHandler := handlers.NewTaskHandler(taskService)
	projectHandler := handlers.NewProjectHandler(projectService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := mux.NewRouter()

	// auth
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/refresh", authHandler.Refresh).Methods("POST")

	protected := r.PathPrefix("/api/auth").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// предмет паблик
	r.HandleFunc("/api/subjects", subjectHandler.GetAllSubjects).Methods("GET")
	r.HandleFunc("/api/subjects/{id}", subjectHandler.GetSubject).Methods("GET")
	r.HandleFunc("/api/subjects/{id}/tasks", taskHandler.GetSubjectTasks).Methods("GET")
	r.HandleFunc("/api/subjects/{id}/members", roleHandler.GetSubjectMembers).Methods("GET")
	r.HandleFunc("/api/tasks/{taskId}", taskHandler.GetTask).Methods("GET")

	// защищенные
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)

	// прдмет
	protectedRouter.HandleFunc("/subjects", subjectHandler.CreateSubject).Methods("POST")
	protectedRouter.HandleFunc("/subjects/{id}", subjectHandler.UpdateSubject).Methods("PUT")
	protectedRouter.HandleFunc("/subjects/{id}", subjectHandler.DeleteSubject).Methods("DELETE")
	protectedRouter.HandleFunc("/subjects/join", subjectHandler.JoinSubject).Methods("POST")
	protectedRouter.HandleFunc("/subjects/my", subjectHandler.GetMySubjects).Methods("GET")

	// роль
	protectedRouter.HandleFunc("/subjects/{id}/roles/{roleId}/change", roleHandler.ChangeRole).Methods("POST")
	protectedRouter.HandleFunc("/subjects/{id}/roles/{roleId}", roleHandler.RemoveFromSubject).Methods("DELETE")

	// задачи
	protectedRouter.HandleFunc("/subjects/{id}/tasks", taskHandler.CreateTask).Methods("POST")
	protectedRouter.HandleFunc("/tasks/{taskId}", taskHandler.UpdateTask).Methods("PUT")
	protectedRouter.HandleFunc("/tasks/{taskId}", taskHandler.DeleteTask).Methods("DELETE")

	// проекты
	r.HandleFunc("/api/tasks/{taskId}/projects", projectHandler.GetTaskProjects).Methods("GET")
	protectedRouter.HandleFunc("/tasks/{taskId}/projects", projectHandler.CreateProject).Methods("POST")
	protectedRouter.HandleFunc("/projects/join", projectHandler.JoinProject).Methods("POST")
	protectedRouter.HandleFunc("/projects/my", projectHandler.GetMyProjects).Methods("GET")

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
