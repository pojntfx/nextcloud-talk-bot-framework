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
		room     = os.Getenv("NEXTCLOUD_ROOM")
	)

	chats, err := client.GetChats(url, username, password, room)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(chats[0])
}
