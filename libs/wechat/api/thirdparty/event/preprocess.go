package event

// PreProcessor 预处理器，对消息进行预处理
type PreProcessor interface {
	Process(data []byte) (string, error)
}

var (
	// 授权事件预处理器
	_preProcessors = map[AuthEventKind]PreProcessor{
		AuthEventKindComponentVerifyTicket: newComponentVerifyTicketEventProcessor(),
		AuthEventKindAuthorized:            newAuthorizedEventProcessor(),
		AuthEventKindUnauthorized:          newUnauthorizedEventProcessor(),
		AuthEventKindUpdateAuthorized:      newAuthorizedEventProcessor(),
	}
)
