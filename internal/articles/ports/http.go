package ports

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rezaAmiri123/test-project/internal/articles/app"
	"github.com/rezaAmiri123/test-project/internal/articles/domain/article"
	"github.com/rezaAmiri123/test-project/internal/common/server"
)

type HttpServer struct {
	app *app.Application
}

type HttpConfig struct {
	HttpServerPort int
	HttpServerAddr string

	AuthServer server.AuthConfig
}

func NewHttpServer(config HttpConfig, application *app.Application) (*http.Server, error) {
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.HttpServerAddr, config.HttpServerPort),
	}
	httpServer := &HttpServer{app: application}
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)
	authMiddleware, err := server.NewAuthMiddleware(config.AuthServer)
	if err != nil {
		return nil, err
	}
	apiRouter.Use(authMiddleware)
	apiRouter.Route("/articles", func(r chi.Router) {
		r.Post("/article", httpServer.CreateArticle)
		r.Get("/article/{slug}", httpServer.GetArticle)
	})
	rootRouter := chi.NewRouter()
	rootRouter.Mount("/api/v1", apiRouter)

	srv.Handler = rootRouter
	return srv, nil
}

func (h *HttpServer) CreateArticle(w http.ResponseWriter, r *http.Request) {
	a := &article.Article{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userUUID := "xxxxxxxxxxxxxxxxxxxxxxx"
	if err := h.app.Commands.CreateArticle.Handle(r.Context(), a, userUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpServer) GetArticle(w http.ResponseWriter, r *http.Request) {
	articleSlug := chi.URLParam(r, "slug")
	a, err := h.app.Queries.GetArticle.Handle(r.Context(), articleSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	//addCorsMiddleware(router)
	//addAuthMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}
