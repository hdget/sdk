package hddapr

import (
	"encoding/json"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"strings"
)

type RouteAnnotation struct {
	*moduleInfo
	Handler       string   // dapr method
	Endpoint      string   // endpoint
	HttpMethods   []string // http methods
	Origin        string   // 请求来源
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共路由
	Permissions   []string // 所属权限列表
	Comments      []string // 备注
}

type rawRouteAnnotation struct {
	Endpoint      string   `json:"endpoint"`      // endpoint
	Methods       []string `json:"methods"`       // http methods
	Origin        string   `json:"origin"`        // 请求来源
	IsRawResponse bool     `json:"isRawResponse"` // 是否返回原始消息
	IsPublic      bool     `json:"isPublic"`      // 是否是公共路由
	Permissions   []string `json:"permissions"`   // 所属权限列表
}

// 注解的前缀
const annotationPrefix = "@hd."
const annotationRoute = annotationPrefix + "route"

// GetRouteAnnotations 获取路由注解
func (b *invocationModuleImpl) GetRouteAnnotations(srcPath string, args ...HandlerMatch) ([]*RouteAnnotation, error) {
	return b.parseRouteAnnotations(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"any", "error"},
		args...,
	)
}

// parseRouteAnnotations 从源代码的注解中解析路由注解
func (b *invocationModuleImpl) parseRouteAnnotations(srcPath, annotationPrefix string, fnParams, fnResults []string, args ...HandlerMatch) ([]*RouteAnnotation, error) {
	matchFn := b.defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := hdutils.AST().InspectFunction(srcPath, fnParams, fnResults, annotationPrefix)
	if err != nil {
		return nil, err
	}

	routeAnnotations := make([]*RouteAnnotation, 0)
	for _, fnInfo := range funcInfos {
		modInfo, err := toModuleInfo(fnInfo.Receiver)
		if err != nil {
			return nil, err
		}

		// 忽略掉不是本模块的备注
		if modInfo.Module != b.Module {
			continue
		}

		// 无路由注解忽略
		ann := fnInfo.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		// 忽略不匹配的函数
		_, matched := matchFn(fnInfo.Function)
		if !matched {
			continue
		}

		// 通过handlerId获取注册时候的invocationHandler
		h := b.handlers[genHandlerId(b.Name, fnInfo.Function)]
		if h == nil {
			continue
		}

		routeAnnotation, err := b.toRouteAnnotation(h.alias, fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routeAnnotations = append(routeAnnotations, routeAnnotation)
	}

	return routeAnnotations, nil
}

// toRouteAnnotation alias为register或者discover handler时候使用的别名
func (b *invocationModuleImpl) toRouteAnnotation(alias string, fnInfo *hdutils.AstFunction, ann *hdutils.AstAnnotation) (*RouteAnnotation, error) {
	// 设置初始值
	raw := &rawRouteAnnotation{
		Endpoint:      "",
		Methods:       []string{"GET"},
		Origin:        "",
		IsRawResponse: false,
		IsPublic:      false,
		Permissions:   []string{},
	}

	// 尝试将注解后的值进行jsonUnmarshal
	if strings.HasPrefix(ann.Value, "{") && strings.HasSuffix(ann.Value, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(hdutils.StringToBytes(ann.Value), &raw)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	}

	// 处理特殊情况, 设置缺省值
	if len(raw.Methods) == 0 {
		raw.Methods = []string{"GET"}
	}

	return &RouteAnnotation{
		moduleInfo:    b.moduleInfo,
		Handler:       alias,
		Endpoint:      raw.Endpoint,
		HttpMethods:   raw.Methods,
		Origin:        raw.Origin,
		IsRawResponse: raw.IsRawResponse,
		IsPublic:      raw.IsPublic,
		Permissions:   raw.Permissions,
		Comments:      fnInfo.PlainComments,
	}, nil
}