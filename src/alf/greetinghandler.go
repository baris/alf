package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type GreetingHandler struct {
	alf *Alf
}

func (h *GreetingHandler) ProcessCurrentEvent() {
	ev := h.alf.currentEvent
	switch ev.Data.(type) {
	case *slack.HelloEvent:
		h.alf.Send("Now tell me you love me!", h.alf.defaultChannel)
	case *slack.MessageEvent:
		msg := ev.Data.(*slack.MessageEvent)
		h.handleMessage(msg)
	}

}

func (h *GreetingHandler) handleMessage(msg *slack.MessageEvent) {
	log.Debug(msg)
	log.Debug(msg.Text)
	text := strings.ToLower(msg.Text)
	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		user, _ := h.alf.api.GetUserInfo(msg.User)
		h.alf.Send("Hey, "+user.Name+"! Hello!", msg.Channel)
	}
}

func (h *GreetingHandler) ProcessIdleEvent() {
}
