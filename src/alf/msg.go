package main

import (
	"strings"
)

func IsToUser(msgText, name, userId string) bool {
	text := strings.ToLower(msgText)
	return strings.HasPrefix(text, name) || strings.HasPrefix(text, "<@"+userId+">")
}

func RemoveMention(msgText, name, userId string) string {
	removeBefore := 0
	text := strings.ToLower(msgText)
	if strings.HasPrefix(text, name) {
		removeBefore = len(name)
	} else if strings.HasPrefix(text, "<@"+userId+">") {
		removeBefore = len("<@" + userId + ">")
	}
	return strings.TrimLeft(msgText[removeBefore:], ": ")
}
