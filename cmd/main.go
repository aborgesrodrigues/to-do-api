package main

import (
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/aborgesrodrigues/to-do-api/internal/audit"
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		// DisableRequestAuditLogs:  true,
		// DisableResponseAuditLogs: true,
		Config: audit.Config{
			BufferSize: 1000000,
		},
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
		r.Route("/debug/pprof", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Post("/symbol", pprof.Symbol)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			// Rotas para todas as profiles específicas, ex: /debug/pprof/goroutine
			r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
			r.Get("/block", pprof.Handler("block").ServeHTTP)
			r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
			r.Get("/heap", pprof.Handler("heap").ServeHTTP)
			r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
			r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		})
		r.Route("/metrics", func(r chi.Router) {
			// Nenhum middleware aplicado aqui
			r.Get("/", promhttp.Handler().ServeHTTP)
		})

		// no JWT
		r.Route("/token", func(r chi.Router) {
			r.Post("/", hdl.AddUser)
		})

		// refresh token
		r.Route("/users/{Id}/refresh_token", func(r chi.Router) {
			r.Use(hdl.IdMiddleware)
			r.Use(hdl.VerifyRefreshJWT)
			r.Get("/", hdl.RefreshToken)
		})

		// with JWT
		r.Route("/", func(r chi.Router) {
			r.Use(logging.RequestLogger(hdl.Logger))
			r.Use(handlers.AccessLogger(handlers.AccessLoggerOptions{
				HTTPAuditLogger: hdl.AuditLogger,
			}))
			r.Use(hdl.VerifyJWT)
			r.Use(hdl.Metrics)

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
