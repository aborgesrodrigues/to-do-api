package main

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	svc    *service.Service
}

func newHandler() *handler {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Error creating logger")
	}

	svc, err := service.New(service.Config{Logger: logger})
	if err != nil {
		panic(err)
	}

	return &handler{
		logger: logger,
		svc:    svc,
	}
}

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Error creating logger")
	}

	hdl := newHandler()

	logger.Info("Server listening.", zap.String("addr", "8080"))
	if err := http.ListenAndServe(":8080", getRouter(hdl)); err != nil {
		logger.Error(err.Error())
	}
}

func getRouter(svc *handler) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(svc.loggerMiddleware)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", svc.listUsers)
			r.Post("/", svc.addUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(svc.idMiddleware)
				r.Get("/", svc.getUser)
				r.Put("/", svc.updateUser)
				r.Delete("/", svc.deleteUser)
				r.Get("/tasks", svc.getUserTasks)
			})
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", svc.listTasks)
			r.Post("/", svc.addTask)
			r.Route("/{taskId}", func(r chi.Router) {
				r.Use(svc.idMiddleware)
				r.Get("/", svc.getTask)
				r.Put("/", svc.updateTask)
				r.Delete("/", svc.deleteTask)
			})
		})
	})

	return r
}
