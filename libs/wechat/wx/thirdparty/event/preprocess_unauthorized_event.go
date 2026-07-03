package event

import (
	"encoding/xml"
)

type xmlUnauthorizedEvent struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

type unauthorizedEventProcessor struct {
}

func newUnauthorizedEventProcessor() PreProcessor {
	return &unauthorizedEventProcessor{}
}

func (h unauthorizedEventProcessor) Process(data []byte) (string, error) {
	var e xmlUnauthorizedEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}

	return e.AuthorizerAppid, nil

}
