package main

import (
	"log"
	"net/url"
	"os"

	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/bots"
	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/protocol"
)

func main() {
	testUsername := os.Getenv("NEXTCLOUD_USERNAME")
	testPassword := os.Getenv("NEXTCLOUD_PASSWORD")
	testURL := os.Getenv("NEXTCLOUD_URL")
	addr, err := url.Parse(testURL)
	if err != nil {
		log.Fatal(err)
	}

	msgChan := make(chan protocol.Message)
	bot := bots.NewNextcloudTalk(addr, testUsername, testPassword, msgChan)

	go func() {
		for {
			if err := bot.ReadRooms(); err != nil {
				log.Println(err)
			}
		}
	}()

	go func() {
		for {
			if err := bot.ReadMessages(); err != nil {
				log.Println(err)
			}
		}
	}()

	for {
		msg := <-msgChan

		log.Println(msg)
	}
}
