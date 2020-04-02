package main

import (
	"log"
	"os"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/v2/pkg/client"
)

func main() {
	var (
		url      = os.Getenv("NEXTCLOUD_URL")
		username = os.Getenv("NEXTCLOUD_USERNAME")
		password = os.Getenv("NEXTCLOUD_PASSWORD")
	)

	roomChan := make(chan client.Room)
	go func() {
		var lastRooms []client.Room

		for {
			rooms, err := client.GetRooms(url, username, password)
			if err != nil {
				log.Fatal(err)
			}

			for _, room := range rooms {
				exists := false

				for _, lastRoom := range lastRooms {
					if room.ID == lastRoom.ID {
						exists = true

						break
					}
				}

				if !exists {
					roomChan <- room
				}
			}

			lastRooms = rooms

			time.Sleep(time.Second * 5)
		}
	}()

	knownIDs := cmap.New()
	chatChan := make(chan client.Chat)
	go func() {
		for room := range roomChan {
			go func(token string) {
				for {
					lastID, _ := knownIDs.Get(token)

					chats, err := client.GetChats(url, username, password, token)
					if err != nil {
						log.Fatal(err)
					}

					chat := chats[0]
					if chat.ID != lastID {
						chatChan <- chats[0]

						log.Println(token, lastID, chat.ID)

						knownIDs.Set(token, chat.ID)
					}

					time.Sleep(time.Second * 5)
				}
			}(room.Token)
		}
	}()

	for chat := range chatChan {
		log.Println(chat)
	}
}
