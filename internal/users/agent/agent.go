package agent

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/rezaAmiri123/test-project/internal/users/adapters"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
	"github.com/rezaAmiri123/test-project/internal/users/ports"
	"google.golang.org/grpc"
)

type Config struct {
	HttpServerPort int
	HttpServerAddr string
	GRPCServerPort int
	GRPCServerAddr string

	DBConfig adapters.GORMConfig
}

func (c Config) HttpAddr() string {
	return fmt.Sprintf("%s:%d", c.HttpServerAddr, c.HttpServerPort)
}

func (c Config) GRPCAddr() string {
	return fmt.Sprintf("%s:%d", c.GRPCServerAddr, c.GRPCServerPort)
}

type Agent struct {
	Config

	httpServer  *http.Server
	grpcServer  *grpc.Server
	repository  user.Repository
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
		a.setupRepository,
		a.setupApplication,
		a.setupHttpServer,
		a.setupGRPCServer,
	}
	for _, fn := range setupsFn {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupHttpServer() error {
	httpServer, err := ports.NewHttpServer(a.HttpAddr(), a.Application)
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

func (a *Agent) setupGRPCServer() error {
	serverConfig := ports.GRPCConfig{a.Application}
	var opts []grpc.ServerOption
	var err error
	a.grpcServer, err = ports.NewGRPCServer(&serverConfig, opts...)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", a.GRPCAddr())
	if err != nil {
		return err
	}
	go func() {
		if err := a.grpcServer.Serve(ln); err != nil {
			_ = a.Shutdown()
		}
	}()
	return err
}

func (a *Agent) setupRepository() error {
	repository, err := adapters.NewGORMUserRepository(a.DBConfig)
	if err != nil {
		return err
	}
	a.repository = repository
	return nil
}

func (a *Agent) setupApplication() error {
	application := app.NewApplication(a.repository)
	a.Application = application
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
		func() error {
			a.grpcServer.GracefulStop()
			return nil
		},
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
