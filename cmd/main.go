package main

import (
	"net/http"
	"os"

	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func main() {
	viper.AutomaticEnv()

	logger := getLogger()

	hdl := handlers.New(logger)

	logger.Info("Server listening.", zap.String("addr", "8080"))
	if err := http.ListenAndServe(":8080", getRouter(hdl)); err != nil {
		logger.Error(err.Error())
	}
}

func getRouter(svc *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(svc.LoggerMiddleware)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", svc.ListUsers)
			r.Post("/", svc.AddUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(svc.IdMiddleware)
				r.Get("/", svc.GetUser)
				r.Put("/", svc.UpdateUser)
				r.Delete("/", svc.DeleteUser)
				r.Get("/tasks", svc.ListUserTasks)
			})
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", svc.ListTasks)
			r.Post("/", svc.AddTask)
			r.Route("/{taskId}", func(r chi.Router) {
				r.Use(svc.IdMiddleware)
				r.Get("/", svc.GetTask)
				r.Put("/", svc.UpdateTask)
				r.Delete("/", svc.DeleteTask)
			})
		})
	})

	return r
}

func getLogger() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger = logger.With(zap.String("app", "myapp")).With(zap.String("environment", "psm"))

	return logger
}
