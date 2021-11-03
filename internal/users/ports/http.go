package ports

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

type HttpServer struct {
	app *app.Application
}

func NewHttpServer(addr string, application *app.Application) (*http.Server,error) {
	httpServer := &HttpServer{app: application}
	router := newRouter(httpServer)
	return &http.Server{
		Addr:    addr,
		Handler: router,
	},nil
}

func (h *HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := &user.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.app.Commands.CreateUser.Handle(r.Context(), u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	u, err := h.app.Queries.GetProfile.Handler(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newRouter(httpServer *HttpServer) chi.Router {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)
	apiRouter.Route("/users", func(r chi.Router) {
		r.Get("/{username}", httpServer.GetProfile)
		r.Post("/register", httpServer.CreateUser)
	})

	rootRouter := chi.NewRouter()
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api/v1", apiRouter)
	return rootRouter
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
