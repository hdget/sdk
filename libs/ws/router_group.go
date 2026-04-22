package lib_ws

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	ginRouterGroup *gin.RouterGroup
	UrlPrefix      string
}

func (b baseServer) NewRouterGroup(urlPrefix string) *RouterGroup {
	return &RouterGroup{
		ginRouterGroup: b.engine.Group(urlPrefix),
		UrlPrefix:      urlPrefix,
	}
}

func (rg *RouterGroup) Use(middlewares ...gin.HandlerFunc) *RouterGroup {
	rg.ginRouterGroup.Use(middlewares...)
	return rg
}

func (rg *RouterGroup) AddRoute(routes ...*Route) error {
	routeMap := make(map[string]struct{})
	for _, route := range routes {
		// 如果是重复的路由忽略掉
		k := fmt.Sprintf("%s_%s", route.Method, route.Path)
		if _, exist := routeMap[k]; exist {
			continue
		}

		// 记录已添加的路由
		routeMap[k] = struct{}{}
		switch strings.ToUpper(route.Method) {
		case "GET":
			rg.ginRouterGroup.GET(route.Path, route.Handler)
		case "POST":
			rg.ginRouterGroup.POST(route.Path, route.Handler)
		default:
			return fmt.Errorf("invalid route method, method: %s", route.Method)
		}
	}
	return nil
}
