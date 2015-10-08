package main

import (
	"path/filepath"

	"github.com/nlopes/slack"
	"github.com/yuin/gopher-lua"
)

type ScriptsHandler struct {
	alf *Alf
}

func (h *ScriptsHandler) Help() string {
	return ""
}

func (h *ScriptsHandler) ProcessCurrentEvent() {
}

func (h *ScriptsHandler) ProcessMessage(msg *slack.MessageEvent) {
	scripts, err := filepath.Glob(h.alf.scriptsDir + "/*.lua")
	if err != nil {
		log.Error("Cannot find scripts file")
		return
	}
	for _, script := range scripts {
		ret := callScript(script, "processMessage", msg.Text)
		if ret != "" {
			h.alf.Send(ret, msg.Channel)
		}
	}
}

func (h *ScriptsHandler) ProcessIdleEvent() {
}

func callScript(scriptPath, method, input string) string {
	ls := lua.NewState()
	defer ls.Close()

	if err := ls.DoFile(scriptPath); err != nil {
		log.Error("Failed to load script file", err)
	}

	param := lua.P{
		Fn:      ls.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}
	if err := ls.CallByParam(param, lua.LString(input)); err != nil {
		log.Error("Failed to call method", err)
	}

	ret := ls.Get(-1)
	ls.Pop(1)

	if ret.Type() != lua.LTNil {
		return ret.String()
	}
	return ""
}
