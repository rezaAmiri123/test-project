package ports_test

import (
	"context"
	"fmt"
	api "github.com/rezaAmiri123/test-project/internal/common/genproto/users"
	"github.com/rezaAmiri123/test-project/internal/users/adapters"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
	"github.com/rezaAmiri123/test-project/internal/users/ports"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

type Config struct {
}

func TestGRPCServer(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		client api.UsersServiceClient,
		config *Config,
		application *app.Application,
	){
		"login test":
		testGrpcLoginUser,
		"verify token":
		testGrpcServer_VerifyTokenUser,
	} {
		t.Run(scenario, func(t *testing.T) {
			client, config, application, teardown := setupGRPCServerTest(t, nil)
			defer teardown()
			fn(t, client, config, application)
		})
	}
}

func setupGRPCServerTest(t *testing.T, fn func(config *Config)) (
	client api.UsersServiceClient,
	cfg *Config,
	application *app.Application,
	teardown func(),
) {
	t.Helper()
	repo := adapters.NewMemoryUserRepository()
	application = app.NewApplication(repo)
	httpPorts := dynaport.Get(1)
	bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", httpPorts[0])
	serverConfig := ports.GRPCConfig{application}
	var opts []grpc.ServerOption

	grpcServer, err := ports.NewGRPCServer(&serverConfig, opts...)
	ln, err := net.Listen("tcp", bindAddr)
	require.NoError(t, err)
	go func() {
		grpcServer.Serve(ln)
	}()

	clientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc, err := grpc.Dial(bindAddr, clientOptions...)
	require.NoError(t, err)
	client = api.NewUsersServiceClient(cc)
	return client, cfg,application, func() {
		grpcServer.Stop()
		cc.Close()
		ln.Close()
	}
}

func testGrpcLoginUser(t *testing.T, client api.UsersServiceClient, config *Config, application *app.Application) {
	ctx := context.Background()
	want := &user.User{
		Username: "username",
		Password: "password",
		Email: "email@example.com",
	}
	err := application.Commands.CreateUser.Handle(ctx, want)
	require.NoError(t,err)
	token, err := client.Login(ctx, &api.LoginRequest{
		Username: want.Username,
		Password: want.Password,
	})
	require.NoError(t, err)
	log.Println(token.GetToken())
}

func testGrpcServer_VerifyTokenUser(t *testing.T, client api.UsersServiceClient, config *Config, application *app.Application) {
	ctx := context.Background()
	want := &user.User{
		Username: "username",
		Password: "password",
		Email: "email@example.com",
	}
	err := application.Commands.CreateUser.Handle(ctx, want)
	require.NoError(t,err)
	token, err := client.Login(ctx, &api.LoginRequest{
		Username: want.Username,
		Password: want.Password,
	})
	require.NoError(t, err)
	got ,err :=client.VerifyToken(ctx, &api.VerifyTokenRequest{
		Token: token.GetToken(),
	})
	require.NoError(t,err)
	require.Equal(t,got.Username,want.Username)
}
