package main

import (
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type WhatisHandler struct {
}

func (h *WhatisHandler) ProcessCurrentEvent() {
}

func (h *WhatisHandler) Help() string {
	help := `NICK: what is [QUERY STRING]  -- lookup QUERY STRING.
NICK: know that [QUERY STRING] is [VALUE STRING] -- set QUERY STRING as VALUE STRING.
NICK: forget [QUERY STRING] -- delete QUERY STRING.
NICK: what do you know? -- brain dump.
`
	return strings.Replace(help, "NICK:", alf.name+":", -1)
}

func (h *WhatisHandler) ProcessMessage(msg *slack.MessageEvent) {
	name := alf.name
	userId := alf.getUserID(alf.name)
	if !IsToUser(msg.Text, name, userId) {
		return
	}
	text := strings.ToLower(RemoveMention(msg.Text, name, userId))
	text = strings.TrimRight(text, ".!?,:;")
	if strings.HasPrefix(text, "what is") {
		text = strings.TrimPrefix(text, "what is")
		text = strings.Trim(text, " ")

		value, err := brain.Get("whatis", text)
		if err == nil && value != "" {
			alf.Send(text+" is "+value+".", msg.Channel)
		} else {
			alf.Send("I don't know what "+text+" is.", msg.Channel)
			if alf.hubotNick != "" {
				alf.Send(alf.hubotNick+": google me "+text, msg.Channel)
			}
		}

	} else if strings.HasPrefix(text, "know that") {
		text = strings.Trim(strings.TrimPrefix(text, "know that"), " ")

		parts := strings.SplitN(text, " is ", 2)
		if len(parts) == 2 {
			brain.Put("whatis", parts[0], parts[1])
			alf.Send("OK!", msg.Channel)
		}

	} else if strings.HasPrefix(text, "forget") {
		text = strings.Trim(strings.TrimPrefix(text, "forget"), " ")
		if _, err := brain.Get("whatis", text); err == nil {
			brain.Delete("whatis", text)
			alf.Send("OK!", msg.Channel)
		}

	} else if strings.HasPrefix(text, "what do you know") {
		all, err := brain.GetAll("whatis")
		if err == nil {
			alf.Send("I know that...", msg.Channel)
			things := make([]string, len(all))
			for k, v := range all {
				things = append(things, k+" is "+v)
			}
			alf.Send(strings.Join(things, "\n"), msg.Channel)
		}
		time.Sleep(3 * time.Second)
		alf.Send(alf.hubotNick+": image me You know nothing, Jon Snow.", msg.Channel)

	}
}

func (h *WhatisHandler) ProcessIdleEvent() {
}
