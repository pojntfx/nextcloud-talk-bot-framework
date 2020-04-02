package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/akrylysov/pogreb"
	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/v2/pkg/client"
)

func main() {
	var (
		url        = os.Getenv("NEXTCLOUD_URL")
		username   = os.Getenv("NEXTCLOUD_USERNAME")
		password   = os.Getenv("NEXTCLOUD_PASSWORD")
		dbLocation = os.Getenv("DB_LOCATION")
		jitsiURL   = os.Getenv("JITSI_URL")
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

	knownIDs, err := pogreb.Open(dbLocation, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer knownIDs.Close()

	chatChan := make(chan client.Chat)
	go func() {
		for room := range roomChan {
			log.Printf(`Joined room "%v" ("%v") with ID "%v" and token "%v"`, room.DisplayName, room.Name, room.ID, room.Token)

			go func(currentRoom client.Room) {
				for {
					lastID := []byte{}
					has, err := knownIDs.Has([]byte(currentRoom.Token))
					if err != nil {
						log.Fatal(err)
					}

					if has {
						lastID, err = knownIDs.Get([]byte(currentRoom.Token))
						if err != nil {
							log.Fatal(err)
						}
					}

					chats, err := client.GetChats(url, username, password, currentRoom.Token)
					if err != nil {
						if err.Error() == "invalid character '<' looking for beginning of value" {
							log.Printf(`Left room "%v" ("%v") with ID "%v" and token "%v"`, currentRoom.DisplayName, currentRoom.Name, currentRoom.ID, currentRoom.Token)

							return
						}

						log.Fatal(err)
					}

					chat := chats[0]
					if strconv.Itoa(chat.ID) != string(lastID) {
						chatChan <- chats[0]

						if err := knownIDs.Put([]byte(currentRoom.Token), []byte(strconv.Itoa(chat.ID))); err != nil {
							log.Fatal(err)
						}
					}

					time.Sleep(time.Second * 5)
				}
			}(room)
		}
	}()

	for chat := range chatChan {
		log.Printf(`Received message from "%v" ("%v") in room "%v" with ID "%v": "%v"`, chat.ActorDisplayName, chat.ActorID, chat.Token, chat.ID, chat.Message)

		reg := regexp.MustCompile("^#videochat")

		if reg.Match([]byte(chat.Message)) {
			client.CreateChat(url, username, password, chat.Token, fmt.Sprintf("@%v started a video call. Tap on %v to join!", chat.ActorID, jitsiURL+"/"+chat.Token))
		}
	}
}
