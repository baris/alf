package main

import (
	"path/filepath"

	"github.com/nlopes/slack"
	"github.com/yuin/gopher-lua"
)

var brain *Brain // TODO: This is too hacky. Get rid of this global!!!!

type ScriptsHandler struct {
	alf *Alf
}

func (h *ScriptsHandler) Help() string {
	return ""
}

func (h *ScriptsHandler) ProcessCurrentEvent() {
}

func (h *ScriptsHandler) ProcessMessage(msg *slack.MessageEvent) {
	if brain == nil {
		brain = h.alf.brain // TODO: Get rid of this global!
	}

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

	ls.SetGlobal("AlfBrainGet", ls.NewFunction(BrainGet))
	ls.SetGlobal("AlfBrainPut", ls.NewFunction(BrainPut))

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

func BrainGet(L *lua.LState) int {
	bucket := L.ToString(1) // get the first argument
	key := L.ToString(2)    // get the second argument

	value, _ := brain.Get(bucket, key)

	L.Push(lua.LString(value)) // push result
	return 1                   // number of results
}

func BrainPut(L *lua.LState) int {
	bucket := L.ToString(1) // get the first argument
	key := L.ToString(2)    // get the second argument
	value := L.ToString(3)  // get the third argument

	err := brain.Put(bucket, key, value)
	if err != nil {
		L.Push(lua.LBool(false))
	} else {
		L.Push(lua.LBool(true))
	}

	return 1 // number of results
}
