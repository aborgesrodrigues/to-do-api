package main

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/aborgesrodrigues/to-do-api/internal/elastic"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
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
	zapConfig := zap.NewProductionConfig()

	const endpointKey = "ELASTICSEARCH_ENDPOINT"
	endpoint := viper.GetString(endpointKey)

	const indexKey = "ELASTICSEARCH_INDEX"
	index := viper.GetString(indexKey)

	const usernameKey = "ELASTICSEARCH_USERNAME"
	username := viper.GetString(usernameKey)

	const passwordKey = "ELASTICSEARCH_PASSWORD"
	password := viper.GetString(passwordKey)

	elasticConfig := elastic.Config{
		Level:    &zapConfig.Level,
		Endpoint: endpoint,
		Username: &username,
		Password: &password,
		Index:    index,
	}

	elasticCore, err := elasticConfig.NewCore()
	if err != nil {
		panic(err)
	}

	logger, err := zapConfig.Build(
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, elasticCore)
		}),
		zap.Fields(zap.String("application", "to-do")),
	)
	if err != nil {
		panic(err)
	}

	zap.RedirectStdLog(logger)

	// Not logging elastic password, obviously.
	logger.Info("Logger configuration.",
		zap.String(endpointKey, endpoint),
		zap.String(indexKey, index),
		zap.String(usernameKey, username),
	)

	return logger
}
