package message

import "encoding/json"

// LocationNormalMessage 位置消息
type LocationNormalMessage struct {
	*Message
	Location *LocationContent `json:"location"`
}

// LocationContent 位置内容
type LocationContent struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
	Name      string  `json:"name"`      // 位置名
	Address   string  `json:"address"`   // 地址详情说明
}

var _ Messager = (*LocationNormalMessage)(nil)

// newLocationNormalMessage 创建位置消息
func newLocationNormalMessage(msg *Message, data []byte) (Messager, error) {
	var fullMsg struct {
		MsgType  string           `json:"msgtype"`
		Location *LocationContent `json:"location"`
	}

	if err := json.Unmarshal(data, &fullMsg); err != nil {
		return nil, err
	}

	return &LocationNormalMessage{
		Message:  msg,
		Location: fullMsg.Location,
	}, nil
}

// GetKind 获取消息类型
func (m *LocationNormalMessage) GetKind() Kind {
	return MessageKindLocation
}
