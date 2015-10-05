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
		text = strings.TrimLeft(text, "what is")
		text = strings.Trim(text, " ")

		response, err := h.alf.brain.Get("whatis", text)
		if err == nil {
			h.alf.Send(text+" is "+response+".", msg.Channel)
		} else {
			h.alf.Send("I don't know what "+text+" is.", msg.Channel)
		}
	} else if strings.HasPrefix(text, "know that") {
		text = strings.TrimLeft(text, "know that")
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
