package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
)

type MediumHandler struct {
}

type TopStoriesResponse struct {
	Payload TopStoriesPayload `json:payload`
}

type TopStoriesPayload struct {
	Value TopStoriesValue `json:value`
}

type TopStoriesValue struct {
	Posts []TopStoriesPost `json:posts`
}

type TopStoriesPost struct {
	Id      string            `json:id`
	Title   string            `json:title`
	Creator TopStoriesCreator `json:creator`
}

type TopStoriesCreator struct {
	Name     string `json:name`
	Username string `json:username`
}

func (h *MediumHandler) ProcessCurrentEvent() {
}

func (h *MediumHandler) Help() string {
	help := `NICK: medium top -- list top 5 stories from medium.com
NICK: medium top all -- list all top stories on the home page
`
	return strings.Replace(help, "NICK:", alf.name+":", -1)
}

func (h *MediumHandler) ProcessMessage(msg *slack.MessageEvent) {
	name := alf.name
	userId := alf.getUserID(alf.name)
	if !IsToUser(msg.Text, name, userId) {
		return
	}
	var stories []string
	text := strings.ToLower(RemoveMention(msg.Text, name, userId))
	if strings.HasPrefix(text, "medium top") {
		limit := 5
		if strings.HasPrefix(text, "medium top all") {
			limit = -1
		}
		for index, post := range getTopStories() {
			story := fmt.Sprintf(
				"⚫ %s\n\t⤷by %s (%s) 〜 http://medium.com/p/%s",
				post.Title,
				post.Creator.Name,
				post.Creator.Username,
				post.Id,
			)
			stories = append(stories, story)

			if index == limit {
				break
			}
		}
		alf.Send(strings.Join(stories, "\n"), msg.Channel)
	}
}

func (h *MediumHandler) ProcessIdleEvent() {
}

func getTopStories() []TopStoriesPost {
	url := "https://api.medium.com/top-stories?format=json"
	res, err := http.Get(url)
	if err != nil {
		log.Error("Failed to fetch medium's top stories: ", err)
		return make([]TopStoriesPost, 0)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to parse medium's top stories: ", err)
		return make([]TopStoriesPost, 0)
	}

	// remove XSSI protection
	body = body[16:]

	var result TopStoriesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error("Failed to unmarshall results: ", err)
		return make([]TopStoriesPost, 0)
	}

	return result.Payload.Value.Posts
}
