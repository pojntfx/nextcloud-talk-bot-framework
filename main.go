package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/client"
)

func main() {
	var (
		url        = os.Getenv("BOT_NEXTCLOUD_URL")
		username   = os.Getenv("BOT_NEXTCLOUD_USERNAME")
		password   = os.Getenv("BOT_NEXTCLOUD_PASSWORD")
		dbLocation = os.Getenv("BOT_DB_LOCATION")
		jitsiURL   = os.Getenv("BOT_JITSI_URL")
	)

	log.Printf(`Starting NextcloudTalk-Jitsi-Bot on "%v" (Jitsi: "%v")`, url, jitsiURL)

	chatChan, statusChan := make(chan client.Chat), make(chan string)
	bot := client.NewNextcloudTalk(url, username, password, dbLocation, chatChan, statusChan)

	defer bot.Close()
	if err := bot.Open(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			if err := bot.ReadRooms(); err != nil {
				log.Println(err)
			}
		}
	}()

	go func() {
		for {
			if err := bot.ReadChats(); err != nil {
				log.Println(err)
			}
		}
	}()

	go func() {
		for status := range statusChan {
			log.Println(status)
		}
	}()

	for chat := range chatChan {
		log.Printf(`Received message from "%v" ("%v") in room "%v" with ID "%v": "%v"`, chat.ActorDisplayName, chat.ActorID, chat.Token, chat.ID, chat.Message)

		reg := regexp.MustCompile("^#video(chat|call)")

		if reg.Match([]byte(chat.Message)) {
			log.Printf(`"%v" ("%v") has requested a video call in room "%v" with ID "%v"; creating video call.`, chat.ActorDisplayName, chat.ActorID, chat.Token, chat.ID)

			bot.CreateChat(chat.Token, fmt.Sprintf("@%v started a video call. Tap on %v to join!", chat.ActorID, jitsiURL+"/"+chat.Token))
		}
	}
}
