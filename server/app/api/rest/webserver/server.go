package rest

import (
	"context"
	"fmt"

	"github.com/IgorAndrade/analytics-twitter/server/app/api"
	"github.com/IgorAndrade/analytics-twitter/server/app/api/rest/middleware"
	route "github.com/IgorAndrade/analytics-twitter/server/app/api/rest/routes"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/labstack/echo/v4"
)

//Server struct
type server struct {
	ctx    context.Context
	server *echo.Echo
	cancel context.CancelFunc
}

//NewServer Create a new Rest server
func NewServer(ctx context.Context, cancel context.CancelFunc) api.Server {
	e := echo.New()
	middleware.ApplyMiddleware(e)
	route.ApplyRoutes(e)
	return &server{
		server: e,
		ctx:    ctx,
		cancel: cancel,
	}
}

//Start a rest server
func (s server) Start() error {
	c := config.Container.Get(config.NAME).(*config.Config)
	defer s.cancel()
	return s.server.Start(c.Rest.Port)
}

//Stop a rest server
func (s server) Stop() error {
	fmt.Println("Stopping rest Server")
	return s.server.Close()
}
