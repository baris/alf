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
	alf *Alf
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
	return `medium top -- list top 5 stories from medium.com
medium top all -- list all top stories on the home page
`
}

func (h *MediumHandler) ProcessMessage(msg *slack.MessageEvent) {
	text := strings.ToLower(msg.Text)
	if strings.HasPrefix(text, h.alf.name) {
		text = strings.TrimPrefix(text, h.alf.name)
		text = strings.TrimLeft(text, ":@ ")
		text = strings.TrimRight(text, ".!?,:;")
	} else {
		return
	}

	stories := make([]string, 0)
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
		h.alf.Send(strings.Join(stories, "\n"), msg.Channel)
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

	return result.Payload.Value.Posts
}
