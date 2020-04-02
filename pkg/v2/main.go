package main

import (
	"log"
	"os"
	"time"

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

	chatChan := make(chan client.Chat)
	for i := range rooms {
		go func(token string) {
			for {
				chats, err := client.GetChats(url, username, password, token)
				if err != nil {
					log.Fatal(err)
				}

				chatChan <- chats[0]

				time.Sleep(time.Second * 5)
			}
		}(rooms[i].Token)
	}

	for chat := range chatChan {
		log.Println(chat)
	}
}
