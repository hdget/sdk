package event

import (
	"encoding/xml"
)

type xmlAuthAuthorizedEvent struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
}

type authorizedEventProcessor struct {
}

func newAuthorizedEventProcessor() PreProcessor {
	return &authorizedEventProcessor{}
}

func (h authorizedEventProcessor) Process(data []byte) (string, error) {
	var e xmlAuthAuthorizedEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}
	return e.AuthorizationCode, nil
}
