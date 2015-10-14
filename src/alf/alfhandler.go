package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type AlfHandler struct {
}

func (h *AlfHandler) Help() string {
	help := `NICK: help -- print this help string
NICK: default channel is [CHANNEL] -- sets the default channel to CHANNEL
`
	return strings.Replace(help, "NICK:", alf.name+":", -1)
}

func (h *AlfHandler) ProcessCurrentEvent() {
	ev := alf.currentEvent
	switch ev.Data.(type) {
	case *slack.HelloEvent:
		alf.Send("Now tell me you love me!", alf.defaultChannel)
	}
}

func (h *AlfHandler) ProcessMessage(msg *slack.MessageEvent) {
	name := alf.name
	userId := alf.getUserID(alf.name)
	if !IsToUser(msg.Text, name, userId) {
		return
	}
	text := strings.ToLower(RemoveMention(msg.Text, name, userId))
	text = strings.TrimRight(text, ".!?,:;")
	if text == "help" {
		var help []string
		for _, handler := range alf.handlers {
			if handler.Help() != "" {
				help = append(help, handler.Help())
			}
		}
		alf.Send(strings.Join(help, "-----\n"), msg.Channel)
	} else if strings.HasPrefix(text, "default channel is ") {
		text = strings.TrimPrefix(text, "default channel is ")
		text = strings.Trim(text, " ")
		if alf.IsMemberOf(text, alf.name) {
			alf.Send("OK! Default channel is now #"+text, msg.Channel)
			alf.defaultChannel = text
		} else {
			alf.Send("Nope, can't do that. You need to invite me there first.", msg.Channel)
		}
	}
}

func (h *AlfHandler) ProcessIdleEvent() {
}
