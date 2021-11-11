package ports

import (
	"context"

	api "github.com/rezaAmiri123/test-project/internal/common/genproto/users"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCConfig struct {
	App *app.Application
}

var _ api.UsersServiceServer = (*grpcServer)(nil)

func newGRPCServer(config *GRPCConfig) (*grpcServer, error) {
	srv := &grpcServer{
		GRPCConfig: config,
	}
	return srv, nil
}

type grpcServer struct {
	*GRPCConfig
	api.UnimplementedUsersServiceServer
}

func NewGRPCServer(config *GRPCConfig, opts ...grpc.ServerOption) (*grpc.Server, error) {
	gsrv := grpc.NewServer(opts...)
	srv, err := newGRPCServer(config)
	if err != nil {
		return nil, err
	}
	api.RegisterUsersServiceServer(gsrv, srv)
	return gsrv, nil
}

func (s *grpcServer) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	u, err := s.App.Queries.GetProfile.Handle(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	token, err := u.GenerateJWTToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not generate token")
	}
	return &api.LoginResponse{Token: token}, nil
}

func (s *grpcServer) VerifyToken(ctx context.Context, req *api.VerifyTokenRequest) (*api.User, error) {
	u, err := s.App.Queries.GetUserToken.Handler(ctx, req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &api.User{
		Username: u.Username,
		Uuid: u.UUID,
	}, nil
}
