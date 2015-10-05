package main

import (
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type WhatisHandler struct {
	alf *Alf
}

func (h *WhatisHandler) ProcessCurrentEvent() {
}

func (h *WhatisHandler) ProcessMessage(msg *slack.MessageEvent) {
	text := strings.ToLower(msg.Text)
	if strings.HasPrefix(text, h.alf.name) {
		text = strings.TrimPrefix(text, h.alf.name)
		text = strings.TrimLeft(text, ":@ ")
		text = strings.TrimRight(text, ".!?,:;")
	} else {
		return
	}

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

	} else if strings.HasPrefix(text, "what do you know") {
		all, err := h.alf.brain.GetAll("whatis")
		if err == nil {
			h.alf.Send("I know that...", msg.Channel)
			for k, v := range all {
				h.alf.Send(k+" is "+v, msg.Channel)
			}
		}
		time.Sleep(2 * time.Second)
		h.alf.Send(h.alf.hubotNick+": image me You know nothing, Jon Snow.", msg.Channel)

	}
}

func (h *WhatisHandler) ProcessIdleEvent() {
}
