package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`
}

// Member to assign task
type Member struct {
	ID   string
	Name string
}

// MemberList to assign task
type MemberList struct {
	members []Member
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)

	lot := &Lot{client: client}
	memberCollector := &MemberCollector{client: client}
	slackListener := &SlackListener{
		client:          client,
		botID:           env.BotID,
		lot:             lot,
		memberCollector: memberCollector,
	}
	go slackListener.ListenAndResponse()

	http.Handle("/taskuji", interactionHandler{
		slackClient:       client,
		verificationToken: env.VerificationToken,
		lot:               lot,
		memberCollector:   memberCollector,
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
