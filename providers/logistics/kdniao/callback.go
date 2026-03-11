package logistics_kdniao

import (
	"context"

	"github.com/hdget/sdk/common/types"
)

// CallbackServer 回调服务器
type CallbackServer struct {
	client   *kdniaoClient
	handlers map[string]types.LogisticsPushHandler // target -> handler
}

// NewCallbackServer 创建回调服务器
func NewCallbackServer(client types.LogisticsClient, handlers ...types.NamedLogisticsPushHandler) *CallbackServer {
	h := make(map[string]types.LogisticsPushHandler)
	for _, nh := range handlers {
		h[nh.Target] = nh.Handler
	}

	// 获取底层客户端
	kdniaoClient, ok := client.(*kdniaoClient)
	if !ok {
		kdniaoClient = nil
	}

	return &CallbackServer{
		client:   kdniaoClient,
		handlers: h,
	}
}

// HandlePush 处理快递鸟推送请求
func (s *CallbackServer) HandlePush(ctx context.Context, req PushRequest) (*PushResponse, error) {
	// 遍历所有推送数据项
	for _, item := range req.Data {
		// 从Callback字段解析target
		target := ParseCallback(item.Callback)

		// 查找对应的处理器
		handler, ok := s.handlers[target]
		if !ok {
			continue // 没有找到处理器，跳过
		}

		// 转换为抽象的TraceItem
		traceItem := types.TraceItem{
			ShipperCode:    item.ShipperCode,
			LogisticCode:   item.LogisticCode,
			Status:         convertKdniaoStatus(item.State),
			Success:        item.Success,
			Reason:         item.Reason,
			Callback:       item.Callback,
			Traces:         convertTraces(item.Traces),
			DeliveryMan:    item.DeliveryMan,
			DeliveryManTel: item.DeliveryManTel,
		}

		// 调用处理器
		err := handler.Handle(ctx, traceItem)
		if err != nil {
			// 记录错误但继续处理其他项
			continue
		}
	}

	return &PushResponse{
		Success: true,
	}, nil
}

// HandlePushWithVerify 处理快递鸟推送请求（带签名验证）
func (s *CallbackServer) HandlePushWithVerify(ctx context.Context, req PushRequest, appKey string) (*PushResponse, error) {
	// 快递鸟推送的签名验证需要根据实际接口实现
	// 这里暂时不实现签名验证，后续根据快递鸟文档补充
	return s.HandlePush(ctx, req)
}

// ParseCallback 从Callback字段解析target
// 格式: "target:xxx" 或直接返回原值
func ParseCallback(callback string) string {
	// 简单实现：如果包含":"则取冒号后的部分
	// 实际格式可能是URL参数或自定义格式
	// 这里支持格式: "target:crm" 或 "?target=crm"
	if callback == "" {
		return ""
	}

	// 尝试解析 "target:xxx" 格式
	if len(callback) > 7 && callback[:7] == "target:" {
		return callback[7:]
	}

	// 尝试解析URL参数格式
	// ?target=crm&other=xxx
	if len(callback) > 8 && callback[:8] == "?target=" {
		// 简单解析
		for i := 8; i < len(callback); i++ {
			if callback[i] == '&' {
				return callback[8:i]
			}
		}
		return callback[8:]
	}

	return callback
}

// AddHandler 添加处理器
func (s *CallbackServer) AddHandler(target string, handler types.LogisticsPushHandler) {
	s.handlers[target] = handler
}

// RemoveHandler 移除处理器
func (s *CallbackServer) RemoveHandler(target string) {
	delete(s.handlers, target)
}
