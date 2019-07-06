package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	key, keypresent := os.LookupEnv("MIRROR_GO_SLACK_OAUTH")

	if !keypresent {
		log.Fatal("No MIRROR_GO_SLACK_OAUTH")
	}

	api := slack.New(
		key,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	fmt.Println("Got API and...")

	// Who am I?
	buildUserCache(api)
	fmt.Printf("I am %s\n", iam.ID)
	fmt.Printf("%+v\n", iam)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			processMessageEvent(ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

}

var cachedUsers []slack.User
var iam slack.User

func buildUserCache(api *slack.Client) {
	users, err := api.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	cachedUsers = users

	for _, user := range cachedUsers {
		if user.Name == "mirrorgo" {
			iam = user
			return
		}
	}
}
func processMessageEvent(ev *slack.MessageEvent) {
	switch ev.Type {

	}
}
