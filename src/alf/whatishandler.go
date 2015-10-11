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

func (h *WhatisHandler) Help() string {
	return `what is [QUERY STRING]  -- lookup QUERY STRING.
know that [QUERY STRING] is [VALUE STRING] -- set QUERY STRING as VALUE STRING.
forget [QUERY STRING] -- delete QUERY STRING.
what do you know? -- brain dump.
`
}

func (h *WhatisHandler) ProcessMessage(msg *slack.MessageEvent) {
	name := h.alf.name
	userId := h.alf.getUserID(h.alf.name)
	if !IsToUser(msg.Text, name, userId) {
		return
	}
	text := strings.ToLower(RemoveMention(msg.Text, name, userId))
	text = strings.TrimRight(text, ".!?,:;")
	if strings.HasPrefix(text, "what is") {
		text = strings.TrimPrefix(text, "what is")
		text = strings.Trim(text, " ")

		value, err := h.alf.brain.Get("whatis", text)
		if err == nil && value != "" {
			h.alf.Send(text+" is "+value+".", msg.Channel)
		} else {
			h.alf.Send("I don't know what "+text+" is.", msg.Channel)
			if h.alf.hubotNick != "" {
				h.alf.Send(h.alf.hubotNick+": google me "+text, msg.Channel)
			}
		}

	} else if strings.HasPrefix(text, "know that") {
		text = strings.Trim(strings.TrimPrefix(text, "know that"), " ")

		parts := strings.SplitN(text, " is ", 2)
		if len(parts) == 2 {
			h.alf.brain.Put("whatis", parts[0], parts[1])
			h.alf.Send("OK!", msg.Channel)
		}

	} else if strings.HasPrefix(text, "forget") {
		text = strings.Trim(strings.TrimPrefix(text, "forget"), " ")
		if _, err := h.alf.brain.Get("whatis", text); err == nil {
			h.alf.brain.Delete("whatis", text)
			h.alf.Send("OK!", msg.Channel)
		}

	} else if strings.HasPrefix(text, "what do you know") {
		all, err := h.alf.brain.GetAll("whatis")
		if err == nil {
			h.alf.Send("I know that...", msg.Channel)
			things := make([]string, len(all))
			for k, v := range all {
				things = append(things, k+" is "+v)
			}
			h.alf.Send(strings.Join(things, "\n"), msg.Channel)
		}
		time.Sleep(3 * time.Second)
		h.alf.Send(h.alf.hubotNick+": image me You know nothing, Jon Snow.", msg.Channel)

	}
}

func (h *WhatisHandler) ProcessIdleEvent() {
}
