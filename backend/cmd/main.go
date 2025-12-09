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
	problemRepo := postgres.NewProblemRepository(db)
	resultRepo := postgres.NewResultRepository(db)
	tokenRepo := redisrepo.NewTokenRepository(redisClient)

	authService := services.NewAuthService(userRepo, tokenRepo, cfg.JWT)
	subjectService := services.NewSubjectService(subjectRepo, roleRepo, userRepo)
	roleService := services.NewRoleService(roleRepo)
	taskService := services.NewTaskService(taskRepo, roleRepo, subjectRepo)
	projectService := services.NewProjectService(projectRepo, taskRepo, roleRepo, problemRepo)
	resultService := services.NewResultService(resultRepo, problemRepo, projectRepo)
	problemService := services.NewProblemService(problemRepo, projectRepo)

	authHandler := handlers.NewAuthHandler(authService)
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	roleHandler := handlers.NewRoleHandler(roleService)
	taskHandler := handlers.NewTaskHandler(taskService)
	projectHandler := handlers.NewProjectHandler(projectService)
	problemHandler := handlers.NewProblemHandler(problemService, resultService)
	resultHandler := handlers.NewResultHandler(resultService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := mux.NewRouter()
	r.Use(CORSMiddleware)

	// Глобальный обработчик для всех OPTIONS запросов
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// auth
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/refresh", authHandler.Refresh).Methods("POST", "OPTIONS")

	protected := r.PathPrefix("/api/auth").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")

	// предмет паблик
	r.HandleFunc("/api/subjects", subjectHandler.GetAllSubjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/subjects/{id}", subjectHandler.GetSubject).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/subjects/{id}/tasks", taskHandler.GetSubjectTasks).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/subjects/{id}/members", roleHandler.GetSubjectMembers).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/tasks/{taskId}", taskHandler.GetTask).Methods("GET", "OPTIONS")

	// защищенные
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)

	// предметы
	protectedRouter.HandleFunc("/subjects/get/my", subjectHandler.GetMySubjects).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/subjects/join", subjectHandler.JoinSubject).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/subjects", subjectHandler.CreateSubject).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/subjects/{id}", subjectHandler.UpdateSubject).Methods("PUT", "OPTIONS")
	protectedRouter.HandleFunc("/subjects/{id}", subjectHandler.DeleteSubject).Methods("DELETE", "OPTIONS")

	// роли
	protectedRouter.HandleFunc("/subjects/{id}/roles/{roleId}/change", roleHandler.ChangeRole).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/subjects/{id}/roles/{roleId}", roleHandler.RemoveFromSubject).Methods("DELETE", "OPTIONS")

	// задания
	protectedRouter.HandleFunc("/subjects/{id}/tasks", taskHandler.CreateTask).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/tasks/{taskId}", taskHandler.UpdateTask).Methods("PUT", "OPTIONS")

	// проекты
	protectedRouter.HandleFunc("/projects/my", projectHandler.GetMyProjects).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/projects/join", projectHandler.JoinProject).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/projects/{projectId}/users", projectHandler.GetProjectUsers).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/tasks/{taskId}/projects", projectHandler.GetTaskProjects).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/tasks/{taskId}/projects", projectHandler.CreateProject).Methods("POST", "OPTIONS")

	// проблемы
	protectedRouter.HandleFunc("/projects/{projectId}/problems", problemHandler.GetProjectProblems).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/problems/{parentId}/subproblems", problemHandler.CreateSubproblem).Methods("POST", "OPTIONS")
	protectedRouter.HandleFunc("/problems/{parentId}/subproblems", problemHandler.GetSubproblems).Methods("GET", "OPTIONS")

	protectedRouter.HandleFunc("/problems/{problemId}", problemHandler.GetProblem).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/problems/{problemId}", problemHandler.UpdateProblem).Methods("PUT", "OPTIONS")

	// результаты
	protectedRouter.HandleFunc("/problems/{problemId}/result", resultHandler.GetResult).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/problems/{problemId}/result", resultHandler.CreateResult).Methods("POST", "OPTIONS")

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(w, r)
	})
}
