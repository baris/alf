package main

import (
	"github.com/nlopes/slack"
	"github.com/yuin/gopher-lua"
)

func luaCallScript(scriptPath, method string) string {
	L := lua.NewState()
	defer L.Close()

	luaAlfLoader(L)

	if err := L.DoFile(scriptPath); err != nil {
		log.Error("Failed to load script file", err)
		return ""
	}

	param := lua.P{
		Fn:      L.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}
	if err := L.CallByParam(param); err != nil {
		log.Error("Failed to call method", err)
		return ""
	}

	ret := L.Get(-1)
	L.Pop(1)

	if ret.Type() != lua.LTNil {
		return ret.String()
	}
	return ""
}

func luaAlfLoader(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"brainGet": luaBrainGet,
		"brainPut": luaBrainPut,
		"msg":      luaMessage,
		"user":     luaMessageUser,
	})
	L.SetField(mod, "name", lua.LString(alf.name))
	L.SetGlobal("alf", mod)
}

func luaBrainGet(L *lua.LState) int {
	bucket := L.ToString(1) // get the first argument
	key := L.ToString(2)    // get the second argument

	value, _ := brain.Get(bucket, key)

	L.Push(lua.LString(value)) // push result
	return 1                   // number of results
}

func luaBrainPut(L *lua.LState) int {
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

func luaMessage(L *lua.LState) int {
	switch alf.currentEvent.Data.(type) {
	case *slack.MessageEvent:
		msg := alf.currentEvent.Data.(*slack.MessageEvent)
		L.Push(lua.LString(msg.Text))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func luaMessageUser(L *lua.LState) int {
	switch alf.currentEvent.Data.(type) {
	case *slack.MessageEvent:
		msg := alf.currentEvent.Data.(*slack.MessageEvent)
		L.Push(lua.LString(alf.getUserName(msg.User)))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}
