package main

import (
	"log"
	"os"

	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/v2/pkg/client"
)

func main() {
	var (
		url      = os.Getenv("NEXTCLOUD_URL")
		username = os.Getenv("NEXTCLOUD_USERNAME")
		password = os.Getenv("NEXTCLOUD_PASSWORD")
	)

	rooms, err := client.GetRooms(url, username, password)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		chats, err := client.GetChats(url, username, password, room.Token)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(chats[0])
	}
}
