package agent

import (
	"context"
	"net/http"
	"sync"

	"github.com/rezaAmiri123/test-project/internal/articles/adapters"
	"github.com/rezaAmiri123/test-project/internal/articles/app"
	"github.com/rezaAmiri123/test-project/internal/articles/app/command"
	"github.com/rezaAmiri123/test-project/internal/articles/app/query"
	"github.com/rezaAmiri123/test-project/internal/articles/ports"
)

type Config struct {
	HttpServerConfig ports.HttpConfig
	DBConfig         adapters.GORMConfig
}

type Agent struct {
	Config

	httpServer  *http.Server
	Application *app.Application

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func NewAgent(config Config) (*Agent, error) {
	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}

	setupsFn := []func() error{
		a.setupApplication,
		a.setupHttpServer,
	}

	for _, fn := range setupsFn {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupApplication() error {
	repo, err := adapters.NewGORMArticleRepository(a.DBConfig)
	if err != nil {
		return err
	}
	application := &app.Application{
		Commands: app.Commands{
			CreateArticle: command.NewCreateArticleHandler(repo),
		},
		Queries: app.Queries{
			GetArticle: query.NewGetArticleHandler(repo),
		},
	}
	a.Application = application
	return nil
}

func (a *Agent) setupHttpServer() error {
	httpServer, err := ports.NewHttpServer(a.HttpServerConfig, a.Application)
	if err != nil {
		return err
	}
	a.httpServer = httpServer
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			_ = a.Shutdown()
		}
	}()

	return nil
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)
	shutdown := []func() error{
		func() error {
			return a.httpServer.Shutdown(context.Background())
		},
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
