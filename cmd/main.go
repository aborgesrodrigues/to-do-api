package main

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Error creating logger")
	}

	svc, err := service.New(service.Config{Logger: logger})
	if err != nil {
		panic(err)
	}

	logger.Info("Server listening.", zap.String("addr", "8080"))
	if err := http.ListenAndServe(":8080", getRouter(svc)); err != nil {
		logger.Error(err.Error())
	}
}

func getRouter(svc *service.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(svc.LoggerMiddleware)
		r.Get("/hello", svc.HelloWorld)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", svc.ListUsers)
			r.Post("/", svc.AddUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(svc.IdMiddleware)
				r.Get("/", svc.GetUser)
				r.Put("/", svc.UpdateUser)
				r.Delete("/", svc.DeleteUser)
				r.Get("/tasks", svc.GetUserTasks)
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
