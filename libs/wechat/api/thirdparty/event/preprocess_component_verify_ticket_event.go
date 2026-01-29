package event

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type xmlComponentVerifyTicketEvent struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

type componentVerifyTicketEventProcessor struct {
}

func newComponentVerifyTicketEventProcessor() PreProcessor {
	return &componentVerifyTicketEventProcessor{}
}

func (h componentVerifyTicketEventProcessor) Process(data []byte) (string, error) {
	var e xmlComponentVerifyTicketEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}

	if e.ComponentVerifyTicket == "" {
		return "", errors.New("empty component verify ticket")
	}

	return e.ComponentVerifyTicket, nil
}
