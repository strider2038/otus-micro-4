package di

import (
	"context"
	"encoding/json"
	"net/http"

	"identity-service/internal/api"
	"identity-service/internal/postgres"
	"identity-service/internal/postgres/database"
	"identity-service/internal/security/argon"
	"identity-service/internal/security/jwks"
	"identity-service/internal/security/jwt"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(connection *pgxpool.Pool, config Config) (http.Handler, error) {
	db := database.New(connection)

	issuer := jwt.NewIssuer(config.PrivateKey, config.SessionLifeTime)
	hasher := argon.Hasher{}

	userRepository := postgres.NewUserRepository(db)
	userApiService := api.NewIdentityApiService(userRepository, hasher, issuer)
	userApiController := api.NewIdentityApiController(userApiService)

	apiRouter := api.NewRouter(userApiController)
	metrics := api.NewMetrics("identity_service")
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
		err := connection.Ping(context.Background())
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`{"status":"ok"}`))
		} else {
			writer.WriteHeader(http.StatusServiceUnavailable)
			writer.Write([]byte(`{"status":"not available"}`))
		}
	})

	router.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(struct {
			Environment        string `json:"environment"`
			ApplicationVersion string `json:"application_version"`
		}{
			Environment:        config.Environment,
			ApplicationVersion: config.Version,
		})
		writer.Write(response)
	})

	jwksHandler, err := jwks.NewHandler(config.PublicKey)
	if err != nil {
		return nil, err
	}

	router.Handle("/.well-known/jwks.json", jwksHandler)

	router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(struct {
			ApplicationName string `json:"application_name"`
			URL             string `json:"url"`
		}{
			ApplicationName: "IdentityService",
			URL:             request.URL.String(),
		})
		writer.Write(response)
	})

	return router, nil
}
