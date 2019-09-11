package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nlopes/slack"
)

// Lot decide member randomly and make slackAttachment.
type Lot struct {
	client *slack.Client
}

// DrawLots decide member randomly and send message to slack.
func (l *Lot) DrawLots(channelID string, members []Member) error {
	if len(members) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	winner := members[rand.Intn(len(members))]
	attachment := slack.Attachment{
		Text:       "",
		Color:      "#42f46e",
		CallbackID: "taskuji",
		Actions: []slack.AttachmentAction{
			{
				Name:  actionAccept,
				Text:  "OK!",
				Type:  "button",
				Style: "primary",
				Value: winner.ID,
			},
			{
				Name:  actionRepeat,
				Text:  "NG:cry:",
				Type:  "button",
				Style: "danger",
			},
		},
		Title:  fmt.Sprintf("I choose you <@%s>! ", winner.ID),
		Footer: "Push the Button",
	}
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := l.client.PostMessage(channelID, fmt.Sprintf("<@%s>", winner.ID), params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
