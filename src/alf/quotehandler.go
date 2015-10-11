package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type QuoteHandler struct {
}

var quotes []string = []string{
	"Call the police! Call the police!",
	"There's hair in this tuna fish...I like it!",
	"Fine, I'll make a peanut butter sandwich...where's the blender?",
	"Now tell me you love me!",
	"What, are you talking to me?",
	"Help, help, I'm stuck in the outhouse!",
	"Don't make me use this!",
	"He's quick, I'll give 'em that!",
	"Lets have a snack now, we'll get friendly later. You got a cat?",
	"Yo, Lucky my man!",
	"Brilliant! This and the letter 'I' in one day.",
	"Hey, don't worry about the old ALFer...Channel 9 is running Psycho!",
	"No problem, just leave me the keys to the liquor cabinet!",
	"I learned one thing about eating jigsaw puzzles...an hour later, you're hungry again.",
	"Yo Kate, where do you keep the casserole dishes? (Why?) The cat won't fit in the toaster.",
	"Did you say I should get hair in the peanut butter, or I shouldn't?",
	"Grease fire! Grease fire!",
	"Nevermind the curtains, put me out!",
	"Oreos!?! My kinda people!",
	"Fine, don't believe me! They didn't believe the boy who cried wolf!",
	"You wanted me to use a flash?!?",
	"Orphins have to eat gruel, and tap-dance with mops!",
	"You want me to press my lips up against your forehead?",
	"I tried to puree a rock...it didn't work.",
	"Hmmm, immediate gratification versus long term security...I'M THINKING, I'M THINKING!",
	"I've decided to reveal myself to the world. This way I can meet new people, travel, see a Grateful Dead concert.",
	"Look at this, they've got me wired for cable! Let's see, which was the button for a cheeseburger?",
	"I wanna be alone. Come on Brian, keep me company!",
	"Why would he even TRY making banana coffee?",
	"Hey, you crawl under people's houses, you hear things.",
	"I like the sauce that Kate opens!",
	"Great, a new baby! We'll raise him as our own.",
	"Shoot bullets through me, I felt like a snack!",
	"Have 'em throw the book at this guy, preferably something by James Mitchman.",
	"This is the way we diaper our kid, diaper our kid, diaper our kid, this is the way we diaper our kid and (drop!)...this is how we drop it.",
	"I still think I should have brought her something, you know? Some candy, some flowers...a rambo doll.",
	"The only good cat is a stir-fried cat.",
	"Did you know that if you eat fast you can eat more?",
	"Oh no! Rain drops are falling on my pig!",
	"Trust me on this one, I've been wrong so many times before.",
	"Quick, quick, hang up, hang up, dial 9-1-1, nine uno uno!",
	"(I was) looking for tomato paste, I broke a tomato.",
	"Naaaaa, that's stupid...I'll do anyway.",
}

func randomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}

func (h *QuoteHandler) Help() string {
	return ""
}

func (h *QuoteHandler) ProcessCurrentEvent() {
}

func (h *QuoteHandler) ProcessMessage(msg *slack.MessageEvent) {
	name := alf.name
	userId := alf.getUserID(alf.name)
	if !IsToUser(msg.Text, name, userId) {
		return
	}
	text := strings.ToLower(RemoveMention(msg.Text, name, userId))
	if strings.HasPrefix(text, "say something") || strings.HasPrefix(text, "talk") {
		alf.Send(randomQuote(), msg.Channel)
	}
}

func (h *QuoteHandler) ProcessIdleEvent() {
	if rand.Intn(86400/alf.updateInterval) == 0 {
		if rand.Intn(2) == 0 && alf.hubotNick != "" {
			alf.Send(alf.hubotNick+": do you feel love?", alf.defaultChannel)
			time.Sleep(3 * time.Second)
			alf.Send(alf.hubotNick+": sarcastic clap", alf.defaultChannel)
		} else {
			alf.Send(randomQuote(), alf.defaultChannel)
		}
	}
}
