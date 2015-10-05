package main

import (
	"time"

	"github.com/nlopes/slack"
)

type Alf struct {
	name           string
	api            *slack.Client
	rtm            *slack.RTM
	brain          *Brain
	users          []slack.User
	channels       []slack.Channel
	handlers       []Handler
	updateInterval int
	defaultChannel string
	currentEvent   slack.RTMEvent
}

func NewAlf(name, token, defaultChannel, databaseFile string, updateInterval int) *Alf {
	alf := new(Alf)
	alf.name = name
	alf.api = slack.New(token)
	alf.rtm = alf.api.NewRTM()
	alf.brain = NewBrain(databaseFile)
	alf.defaultChannel = defaultChannel
	alf.handlers = make([]Handler, 0)
	alf.users = make([]slack.User, 0)
	alf.channels = make([]slack.Channel, 0)
	alf.updateInterval = updateInterval
	alf.api.SetDebug(false)
	return alf
}

func (alf *Alf) start() {
	go alf.rtm.ManageConnection()
	go alf.updateChannels()
	go alf.updateUsers()
	go alf.idleLoop()

	for {
		select {
		case ev := <-alf.rtm.IncomingEvents:
			log.Debug(ev)
			alf.currentEvent = ev
			for _, h := range alf.handlers {
				h.ProcessCurrentEvent()
			}
		}
	}

}

func (alf *Alf) Send(msg, channelNameOrID string) {
	channelID := channelNameOrID
	if channel, err := alf.api.GetChannelInfo(channelID); err != nil || channel == nil {
		channelID = alf.getChannelID(channelNameOrID)
	}
	alf.rtm.SendMessage(alf.rtm.NewOutgoingMessage(msg, channelID))
}

func (alf *Alf) AddHandler(h Handler) {
	alf.handlers = append(alf.handlers, h)
}

func (alf *Alf) updateChannels() {
	for {
		channels, err := alf.api.GetChannels(true)
		alf.channels = channels
		if err != nil {
			log.Error("Cannot update channels: ", err)
		}
		time.Sleep(time.Duration(alf.updateInterval) * time.Second)
	}

}

func (alf *Alf) updateUsers() {
	for {
		users, err := alf.api.GetUsers()
		alf.users = users
		if err != nil {
			log.Error("Cannot update users: ", err)
		}
		time.Sleep(time.Duration(alf.updateInterval) * time.Second)
	}

}

func (alf *Alf) idleLoop() {
	for {
		for _, h := range alf.handlers {
			go h.ProcessIdleEvent()
		}
		time.Sleep(time.Duration(alf.updateInterval) * time.Second)
	}
}

func (alf *Alf) getChannelID(channelName string) string {
	for _, channel := range alf.channels {
		if channel.Name == channelName {
			return channel.ID
		}
	}
	log.Debug("Cannot find channel: ", channelName)
	return ""
}
