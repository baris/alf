package main

import (
	"github.com/nlopes/slack"
)

type Handler interface {
	Help() string
	ProcessCurrentEvent()
	ProcessMessage(*slack.MessageEvent)
	ProcessIdleEvent()
}
