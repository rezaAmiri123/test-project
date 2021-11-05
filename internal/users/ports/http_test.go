package ports_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/rezaAmiri123/test-project/internal/users/adapters"
	"github.com/rezaAmiri123/test-project/internal/users/app"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
	"github.com/rezaAmiri123/test-project/internal/users/ports"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
)

func TestServer(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name string
		Fn   func(t *testing.T, serverAddr string, application *app.Application, teardown func())
	}{
		{
			Name: "create user",
			Fn:   testCreateUser,
		},
		{
			Name: "get profile",
			Fn:   testGetProfile,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			serverAddr, application, teardown := setupTest(t)
			tc.Fn(t, serverAddr, application, teardown)
		})
	}
}

func testCreateUser(t *testing.T, serverAddr string, application *app.Application, teardown func()) {
	defer teardown()
	url := serverAddr + "/api/v1/users/register"
	userData := user.User{
		Username: "username",
		Password: "password",
		Email:    "email@example.com",
	}
	userJson, err := json.Marshal(userData)
	require.NoError(t, err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userJson))
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusCreated)
}

func testGetProfile(t *testing.T, serverAddr string, application *app.Application, teardown func()) {
	defer teardown()
	userData := setTestUser(t, application)
	url := serverAddr + "/api/v1/users/" + userData.Username
	resp, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	var respData user.User
	err = json.Unmarshal(body, &respData)
	require.NoError(t, err)
	require.Equal(t, respData.Username, userData.Username)
}

func setupTest(t *testing.T) (serverAddr string, application *app.Application, teardown func()) {
	repo := adapters.NewMemoryUserRepository()
	application = app.NewApplication(repo)
	httpPorts := dynaport.Get(1)
	bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", httpPorts[0])
	httpServer, err := ports.NewHttpServer(bindAddr, application)
	require.NoError(t, err)
	go func() {
		_ = httpServer.ListenAndServe()
		//fmt.Println(err)
		//require.NoError(t, err)
	}()
	time.Sleep(time.Millisecond)
	bindAddr = "http://" + bindAddr
	return bindAddr, application, func() {
		httpServer.Close()
	}
}

func setTestUser(t *testing.T, application *app.Application) *user.User {
	t.Helper()
	userData := user.User{
		Username: "username",
		Password: "password",
		Email:    "email@example.com",
	}
	err := application.Commands.CreateUser.Handle(context.Background(), &userData)
	require.NoError(t, err)
	return &userData
}
