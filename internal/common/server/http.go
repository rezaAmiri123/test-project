package server

import (
	"fmt"
	"net/http"

	"github.com/rezaAmiri123/test-project/internal/common/auth"
	UserApi "github.com/rezaAmiri123/test-project/internal/common/genproto/users"
	"google.golang.org/grpc"
)

type AuthConfig struct {
	GRPCUserAddr string
	GRPCUserPort int
}

func NewAuthMiddleware(config AuthConfig) (func(http.Handler) http.Handler, error){
	addr := fmt.Sprintf("%s:%d", config.GRPCUserAddr, config.GRPCUserPort)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err!= nil{
		return nil, err
	}
	authClient := UserApi.NewUsersServiceClient(conn)
	return auth.UserHttpMiddleware{AuthClient: authClient}.Middleware, nil
}
