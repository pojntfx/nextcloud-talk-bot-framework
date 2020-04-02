package main

import (
	"log"
	"net/url"
	"os"

	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/bots"
)

func main() {
	testUsername := os.Getenv("NEXTCLOUD_USERNAME")
	testPassword := os.Getenv("NEXTCLOUD_PASSWORD")
	testURL := os.Getenv("NEXTCLOUD_URL")
	addr, err := url.Parse(testURL)
	if err != nil {
		log.Fatal(err)
	}

	bot := bots.NewNextcloudTalk(addr, testUsername, testPassword)

	rooms, err := bot.GetRooms()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Available rooms: %v\n", rooms)
}
