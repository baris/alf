package main

type Handler interface {
	ProcessCurrentEvent()
	ProcessIdleEvent()
}
