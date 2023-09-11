package main

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/aborgesrodrigues/to-do-api/internal/audit"
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

const (
	// Audit logging env vars
	envVarAuditLogS3Bucket    = "AUDITLOG_S3_BUCKET"
	envVarAuditLogS3Directory = "AUDITLOG_S3_DIRECTORY"
	envVarAuditLogS3Endpoint  = "AUDITLOG_S3_ENDPOINT"
	envVarAuditLogS3Region    = "AUDITLOG_S3_REGION"
)

func main() {
	// initiate viper
	viper.AutomaticEnv()

	logger := getLogger()
	var s3Endpoint *string
	if s3EndpointVal := viper.GetString(envVarAuditLogS3Endpoint); s3EndpointVal != "" {
		s3Endpoint = &s3EndpointVal
	}

	auditWriter, err := audit.NewS3Writer(audit.S3Config{
		Bucket:    requireENV(envVarAuditLogS3Bucket),
		Directory: requireENV(envVarAuditLogS3Directory),
		Endpoint:  s3Endpoint,
		Region:    requireENV(envVarAuditLogS3Region),
	})
	if err != nil {
		logger.Fatal("Unable to instantiate S3 audit writer.", zap.Error(err))
	}
	auditLogger, err := logging.NewHTTPAuditLogger(logging.HTTPAuditLogOptions{
		Writer: auditWriter,
	})
	if err != nil {
		logger.Fatal("Unable to instantiate audit logger.", zap.Error(err))
	}
	defer auditLogger.Close()

	hdl := handlers.New(logger, auditLogger)

	logger.Info("Server listening.", zap.String("addr", "8080"))
	if err := http.ListenAndServe(":8080", getRouter(hdl)); err != nil {
		logger.Error(err.Error())
	}
}

func getRouter(hdl *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(logging.RequestLogger(hdl.Logger))
		r.Use(handlers.AccessLogger(handlers.AccessLoggerOptions{
			HTTPAuditLogger: hdl.AuditLogger,
		}))

		// with JWT
		r.Route("/", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// no JWT
				r.Post("/token", hdl.AddUser)
				r.Post("/login", hdl.Login)

				// refresh token
				r.Route("/refresh_token", func(r chi.Router) {
					r.Use(hdl.VerifyRefreshJWT)
					r.Get("/", hdl.RefreshToken)
				})
			})
			r.Group(func(r chi.Router) {
				r.Use(hdl.VerifyJWT)
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
		})

	})

	return r
}

func getLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}

func requireENV(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(key + " not set")
	}
	return value
}
