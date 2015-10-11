package main

import (
	"path/filepath"

	"github.com/nlopes/slack"
)

type ScriptsHandler struct {
}

func (h *ScriptsHandler) Help() string {
	help := ""
	for _, script := range scripts() {
		ret := luaCallScript(script, "help")
		if ret != "" {
			help += ret + "\n"
		}
	}
	return help

}

func (h *ScriptsHandler) ProcessCurrentEvent() {
}

func (h *ScriptsHandler) ProcessMessage(msg *slack.MessageEvent) {
	for _, script := range scripts() {
		ret := luaCallScript(script, "processMessage")
		if ret != "" {
			alf.Send(ret, msg.Channel)
		}
	}
}

func (h *ScriptsHandler) ProcessIdleEvent() {
}

func scripts() []string {
	lst, err := filepath.Glob(alf.scriptsDir + "/*.lua")
	if err != nil {
		log.Error("Cannot find scripts file")
		return nil
	}
	return lst
}
