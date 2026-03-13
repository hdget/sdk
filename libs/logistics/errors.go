package logistics

import "errors"

var (
	// ErrInvalidShipperCode 无效的快递公司编码
	ErrInvalidShipperCode = errors.New("invalid shipper code")
	// ErrInvalidTrackingNo 无效的快递单号
	ErrInvalidTrackingNo = errors.New("invalid tracking number")
	// ErrQueryFailed 查询失败
	ErrQueryFailed = errors.New("query failed")
	// ErrSubscribeFailed 订阅失败
	ErrSubscribeFailed = errors.New("subscribe failed")
	// ErrRecognizeFailed 识别失败
	ErrRecognizeFailed = errors.New("recognize failed")
	// ErrInvalidSign 签名验证失败
	ErrInvalidSign = errors.New("invalid sign")
	// ErrUnknownVendor 未知的供应商
	ErrUnknownVendor = errors.New("unknown vendor")
	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig = errors.New("invalid config")
	// ErrEmptyConfig 空配置
	ErrEmptyConfig = errors.New("empty config")
	// ErrParseCallbackFailed 解析回调失败
	ErrParseCallbackFailed = errors.New("parse callback failed")
)