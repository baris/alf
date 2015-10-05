package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type WhatisHandler struct {
	alf *Alf
}

func (h *WhatisHandler) ProcessCurrentEvent() {
}

func (h *WhatisHandler) ProcessMessage(msg *slack.MessageEvent) {
	text := strings.ToLower(msg.Text)
	if strings.HasPrefix(text, "what is") {
		text = strings.TrimPrefix(text, "what is")
		text = strings.Trim(text, " ")

		response, err := h.alf.brain.Get("whatis", text)
		if err == nil && response != "" {
			h.alf.Send(text+" is "+response+".", msg.Channel)
		} else {
			h.alf.Send("I don't know what "+text+" is.", msg.Channel)
			if h.alf.hubotNick != "" {
				h.alf.Send(h.alf.hubotNick+": google me "+text, msg.Channel)
			}
		}
	} else if strings.HasPrefix(text, "know that") {
		text = strings.TrimPrefix(text, "know that")
		text = strings.Trim(text, " ")

		parts := strings.SplitN(text, " is ", 2)
		if len(parts) == 2 {
			h.alf.brain.Put("whatis", parts[0], parts[1])
			h.alf.Send("OK!", msg.Channel)
		}
	}
}

func (h *WhatisHandler) ProcessIdleEvent() {
}
