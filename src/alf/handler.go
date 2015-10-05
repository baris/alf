package main

import (
	"github.com/nlopes/slack"
)

type Handler interface {
	ProcessCurrentEvent()
	ProcessMessage(*slack.MessageEvent)
	ProcessIdleEvent()
}
