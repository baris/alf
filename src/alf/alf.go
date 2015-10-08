package main

import (
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type Alf struct {
	name           string
	hubotNick      string
	api            *slack.Client
	rtm            *slack.RTM
	brain          *Brain
	users          []slack.User
	channels       []slack.Channel
	handlers       []Handler
	updateInterval int
	scriptsDir     string
	defaultChannel string
	currentEvent   slack.RTMEvent
}

func NewAlf(name, hubotNick, token, defaultChannel, databaseFile, scriptsDir string, updateInterval int) *Alf {
	alf := new(Alf)
	alf.name = name
	alf.hubotNick = hubotNick
	alf.api = slack.New(token)
	alf.rtm = alf.api.NewRTM()
	alf.brain = NewBrain(databaseFile)
	alf.defaultChannel = defaultChannel
	alf.handlers = make([]Handler, 0)
	alf.users = make([]slack.User, 0)
	alf.channels = make([]slack.Channel, 0)
	alf.scriptsDir = scriptsDir
	alf.updateInterval = updateInterval
	alf.api.SetDebug(false)
	return alf
}

func (alf *Alf) start() {
	go alf.rtm.ManageConnection()
	go alf.updateChannelsLoop()
	go alf.updateUsersLoop()
	go alf.idleLoop()

	for {
		select {
		case ev := <-alf.rtm.IncomingEvents:
			log.Debug(ev)
			alf.currentEvent = ev
			switch ev.Data.(type) {
			case *slack.MessageEvent:
				msg := ev.Data.(*slack.MessageEvent)
				for _, h := range alf.handlers {
					h.ProcessMessage(msg)
				}
			}
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

func (alf *Alf) IsMemberOf(channelName, userName string) bool {
	channel := alf.getChannel(channelName)
	userId := alf.getUserID(userName)
	for _, member := range channel.Members {
		if member == userId {
			return true
		}
	}
	return false
}

func (alf *Alf) AddHandler(h Handler) {
	alf.handlers = append(alf.handlers, h)
}

func (alf *Alf) updateChannelsLoop() {
	for {
		channels, err := alf.api.GetChannels(true)
		alf.channels = channels
		if err != nil {
			log.Error("Cannot update channels: ", err)
		}
		time.Sleep(time.Duration(alf.updateInterval) * time.Second)
	}

}

func (alf *Alf) updateUsersLoop() {
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

func (alf *Alf) getChannel(channelName string) slack.Channel {
	for _, channel := range alf.channels {
		if channel.Name == channelName {
			return channel
		}
	}
	log.Debug("Cannot find channel: ", channelName)
	return slack.Channel{}
}

func (alf *Alf) getChannelID(channelName string) string {
	for _, channel := range alf.channels {
		if channel.Name == channelName {
			return strings.ToLower(channel.ID)
		}
	}
	log.Debug("Cannot find channel: ", channelName)
	return ""
}

func (alf *Alf) getUserID(userName string) string {
	for _, user := range alf.users {
		if user.Name == userName {
			return strings.ToLower(user.ID)
		}
	}
	log.Debug("Cannot find user: ", userName)
	return ""
}
