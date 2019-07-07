package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	key, keypresent := os.LookupEnv("MIRROR_GO_SLACK_OAUTH")

	if !keypresent {
		log.Fatal("No MIRROR_GO_SLACK_OAUTH")
	}

	api := slack.New(
		key,
		//		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	fmt.Println("Got API and...")

	// // Who am I?
	// buildUserCache(api)
	// fmt.Printf("I am %s\n", iam.ID)
	// fmt.Printf("%+v\n", iam)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			processInfo(ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %+v\n", ev)
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

var cachedInfo *slack.Info
var cachedUsers []slack.User
var iam *slack.UserDetails
var imatchid string

func processInfo(info *slack.Info) {

	fmt.Printf("%+v\n", info)
	cachedInfo = info

	iam = info.User
	imatchid = "<@" + iam.ID + ">"

	fmt.Println(info.Channels)
	fmt.Println(imatchid)
}

var ignorableTypes = []string{"error", "group_join", "member_joined_channel", "file_created", "hello", "desktop_notification", "user_typing", "file_public", "bot_added", "bot_changed", "apps_changed", "apps_installed", "user_change"}

func processMessageEvent(ev *slack.MessageEvent) {

	ignorable := false

	for _, v := range ignorableTypes {
		if strings.Compare(v, ev.Type) == 0 {
			ignorable = true
			break
		}
	}

	if ignorable {
		fmt.Printf("Ignoring %s\n", ev.Type)
		return
	}

	fmt.Printf("%s %s\n", imatchid, ev.Text)
	if strings.Contains(ev.Text, imatchid) {
		fmt.Printf("Mentions me %+v\n", ev)
	} else {
		fmt.Println("Something else")
	}

}
