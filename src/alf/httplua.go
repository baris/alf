package main

import (
	"io/ioutil"
	"net/http"

	"github.com/yuin/gopher-lua"
)

func luaHTTPLoader(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"GET": luaGET,
	})
	L.SetGlobal("http", mod)
}

func luaGET(L *lua.LState) int {
	url := L.ToString(1)

	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)

	L.Push(lua.LString(body))
	return 1
}
