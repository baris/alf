package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type AlfHandler struct {
	alf *Alf
}

func (h *AlfHandler) Help() string {
	return `help -- print this help string
default channel is [CHANNEL] -- sets the default channel to CHANNEL
`
}

func (h *AlfHandler) ProcessCurrentEvent() {
}

func (h *AlfHandler) ProcessMessage(msg *slack.MessageEvent) {
	text := strings.ToLower(msg.Text)
	if strings.HasPrefix(text, h.alf.name) {
		text = strings.TrimPrefix(text, h.alf.name)
		text = strings.TrimLeft(text, ":@ ")
		text = strings.TrimRight(text, ".!?,:;")
	} else {
		return
	}

	if text == "help" {
		for _, handler := range h.alf.handlers {
			h.alf.Send(handler.Help(), msg.Channel)
		}
	} else if strings.HasPrefix(text, "default channel is ") {
		text = strings.TrimPrefix(text, "default channel is ")
		text = strings.Trim(text, " ")
		if h.alf.IsMemberOf(text, h.alf.name) {
			h.alf.Send("OK! Default channel is now #"+text, msg.Channel)
			h.alf.defaultChannel = text
		} else {
			h.alf.Send("Nope, can't do that. You need to invite me there first.", msg.Channel)
		}
	}
}

func (h *AlfHandler) ProcessIdleEvent() {
}
