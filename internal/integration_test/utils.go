package integrationtest

import (
	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/go-chi/chi/v5"
)

func getRouter(hdl *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", hdl.ListUsers)
			r.Post("/", hdl.AddUser)
			r.Route("/{Id}", func(r chi.Router) {
				r.Use(hdl.IdMiddleware)
				r.Get("/", hdl.GetUser)
				r.Put("/", hdl.UpdateUser)
				r.Delete("/", hdl.DeleteUser)
				r.Get("/tasks", hdl.ListUserTasks)
			})
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", hdl.ListTasks)
			r.Post("/", hdl.AddTask)
			r.Route("/{Id}", func(r chi.Router) {
				r.Use(hdl.IdMiddleware)
				r.Get("/", hdl.GetTask)
				r.Put("/", hdl.UpdateTask)
				r.Delete("/", hdl.DeleteTask)
			})
		})
	})

	return r
}
