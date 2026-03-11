package logistics_kdniao

import (
	"errors"
)

var (
	// ErrInvalidShipperCode 无效的快递公司编码
	ErrInvalidShipperCode = errors.New("invalid shipper code")
	// ErrInvalidLogisticCode 无效的快递单号
	ErrInvalidLogisticCode = errors.New("invalid logistic code")
	// ErrQueryFailed 查询失败
	ErrQueryFailed = errors.New("query failed")
	// ErrSubscribeFailed 订阅失败
	ErrSubscribeFailed = errors.New("subscribe failed")
	// ErrRecognizeFailed 识别失败
	ErrRecognizeFailed = errors.New("recognize failed")
	// ErrInvalidSign 签名验证失败
	ErrInvalidSign = errors.New("invalid sign")
	// ErrNoHandler 无处理器
	ErrNoHandler = errors.New("no handler for target")
)