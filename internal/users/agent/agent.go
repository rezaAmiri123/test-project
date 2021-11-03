package agent

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/rezaAmiri123/test-project/internal/users/adapters"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
	"github.com/rezaAmiri123/test-project/internal/users/ports"
)

type Config struct {
	HttpServerPort int
	HttpServerAddr string
	DBConfig       adapters.GORMConfig
}

func (c Config) HttpAddr() string {
	return fmt.Sprintf("%s:%d", c.HttpServerAddr, c.HttpServerPort)
}

type Agent struct {
	Config

	httpServer *http.Server
	repository user.Repository

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
		a.setupRepository,
		a.setupHttpServer,
	}
	for _, fn := range setupsFn {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupHttpServer() error {
	application := app.NewApplication(a.repository)
	httpServer, err := ports.NewHttpServer(a.HttpAddr(), application)
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
func (a *Agent) setupRepository() error {
	repository, err := adapters.NewGORMUserRepository(a.DBConfig)
	if err != nil {
		return err
	}
	a.repository = repository
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
