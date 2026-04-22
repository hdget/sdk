package lib_ws

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hdget/sdk/common/provider"
)

type Server interface {
	Start() error
	Stop() error
	SetMode(mode string)
	GracefulStop(ctx context.Context) error
	AddRoutes(routes []*Route) error
	NewRouterGroup(urlPrefix string) *RouterGroup
}

type baseServer struct {
	*http.Server
	engine                   *gin.Engine
	gracefulShutdownWaitTime time.Duration
	providers                map[provider.Category]provider.Provider
}

func newBaseServer(address string, options ...Option) *baseServer {
	s := &baseServer{
		Server: &http.Server{
			Addr: address,
		},
		gracefulShutdownWaitTime: defaultGracefulShutdownWaitTime,
		providers:                make(map[provider.Category]provider.Provider),
	}

	for _, apply := range options {
		apply(s)
	}

	// use gin engine
	engine := s.newDefaultGinEngine()
	s.Handler = engine
	s.engine = engine

	return s
}

func (b baseServer) Stop() error {
	if err := b.Close(); err != nil {
		return err
	}
	return nil
}

func (b baseServer) GracefulStop(ctx context.Context) error {
	if err := b.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (b baseServer) AddRoutes(routes []*Route) error {
	routeMap := make(map[string]struct{})
	for _, route := range routes {
		// 先检查是否有重复的路由
		k := fmt.Sprintf("%s_%s", route.Method, route.Path)
		if _, exist := routeMap[k]; exist {
			return fmt.Errorf("duplicate route, url: %s, method: %s", route.Path, route.Method)
		}

		// 添加到router group
		switch strings.ToUpper(route.Method) {
		case "GET":
			b.engine.GET(route.Path, route.Handler)
		case "POST":
			b.engine.POST(route.Path, route.Handler)
		}
	}
	return nil
}

// SetMode set ws to specific mode
func (b baseServer) SetMode(mode string) {
	gin.SetMode(mode)
}

func (b baseServer) newDefaultGinEngine() *gin.Engine {
	p := b.providers[provider.CategoryLogger]
	if p == nil {
		panic("logger provider is nil")
	}

	loggerProvider, ok := p.(provider.Logger)
	if !ok {
		panic("invalid logger provider")
	}

	// new route
	engine := gin.New()

	// set ws to logout to stdout and file
	gin.DefaultWriter = io.MultiWriter(loggerProvider.GetStdLogger().Writer())

	// add basic middleware
	engine.Use(
		gin.Recovery(),
	)

	return engine
}
