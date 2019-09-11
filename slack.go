package main

import (
	"log"
	"strings"

	"github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.
	actionAccept = "appcept"
	actionRepeat = "repeat"
)

// SlackListener listen and response to slack event.
type SlackListener struct {
	client          *slack.Client
	botID           string
	lot             *Lot
	memberCollector *MemberCollector
}

// ListenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response mention to bot. Ignore else.
	if !strings.Contains(ev.Msg.Text, s.botID) || len(ev.Inviter) != 0 {
		return nil
	}
	if len(ev.BotID) != 0 {
		return nil
	}
	members, err := s.memberCollector.Collect(ev.Channel)
	if err != nil {
		return err
	}
	s.lot.DrawLots(ev.Channel, members)

	return nil
}
