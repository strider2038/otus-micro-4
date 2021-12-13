package di

import (
	"net/http"

	"echo-service/internal/api"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() (http.Handler, error) {
	service := api.NewEchoApiService()
	controller := api.NewEchoApiController(service)

	apiRouter := api.NewRouter(controller)
	metrics := api.NewMetrics("echo_service")
	apiRouter.Use(func(handler http.Handler) http.Handler {
		return api.MetricsMiddleware(handler, metrics)
	})

	router := mux.NewRouter()

	router.PathPrefix("/api").Handler(apiRouter)
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"status":"ok"}`))
	})

	router.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"status":"ok"}`))
	})

	return router, nil
}
