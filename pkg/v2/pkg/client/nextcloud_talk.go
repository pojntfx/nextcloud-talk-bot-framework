package client

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/akrylysov/pogreb"
	"github.com/go-resty/resty/v2"
)

// NextcloudTalk is a Nextcloud Talk Chatbot
type NextcloudTalk struct {
	url, username, password, dbLocation string
	chatChan                            chan Chat
	statusChan                          chan string
	knownIDs                            *pogreb.DB
}

// NewNextcloudTalk creates a new Nextcloud Talk Chatbot
func NewNextcloudTalk(url, username, password, dbLocation string, chatChan chan Chat, statusChan chan string) *NextcloudTalk {
	return &NextcloudTalk{
		url, username, password, dbLocation, chatChan, statusChan, nil,
	}
}

func (n *NextcloudTalk) getRooms() ([]Room, error) {
	client := resty.New()

	res, err := client.R().
		SetHeaders(map[string]string{
			"OCS-APIRequest": "true",
			"Accept":         "application/json",
		}).
		SetBasicAuth(n.username, n.password).
		Get(n.url + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "room"))
	if err != nil {
		return nil, err
	}

	var resStruct RoomResponse
	if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
		if err != nil {
			return nil, err
		}
	}

	return resStruct.OCS.Data, nil
}

func (n *NextcloudTalk) getChats(room string) ([]Chat, error) {
	client := resty.New()

	res, err := client.R().
		SetHeaders(map[string]string{
			"OCS-APIRequest": "true",
			"Accept":         "application/json",
		}).
		SetQueryParams(map[string]string{
			"setReadMarker":  "true",
			"lookIntoFuture": "0",
		}).
		SetBasicAuth(n.username, n.password).
		Get(n.url + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "chat", room))
	if err != nil {
		return nil, err
	}

	var resStruct ChatResponse
	if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
		if err != nil {
			return nil, err
		}
	}

	return resStruct.OCS.Data, nil
}

// Open opens the bot
func (n *NextcloudTalk) Open() error {
	knownIDs, err := pogreb.Open(n.dbLocation, nil)
	if err != nil {
		return err
	}
	n.knownIDs = knownIDs

	return nil
}

// Close closes the chat bot
func (n *NextcloudTalk) Close() error {
	return n.knownIDs.Close()
}

// Read starts reading rooms and chats
func (n *NextcloudTalk) Read() error {
	go func() {
		roomChan := make(chan Room)
		go func() {
			var lastRooms []Room

			for {
				rooms, err := n.getRooms()
				if err != nil {
					n.statusChan <- err.Error()

					continue
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

		for room := range roomChan {
			n.statusChan <- fmt.Sprintf(`Joined room "%v" ("%v") with ID "%v" and token "%v"`, room.DisplayName, room.Name, room.ID, room.Token)

			go func(currentRoom Room) {
				for {
					lastID := []byte{}
					has, err := n.knownIDs.Has([]byte(currentRoom.Token))
					if err != nil {
						n.statusChan <- err.Error()

						continue
					}

					if has {
						lastID, err = n.knownIDs.Get([]byte(currentRoom.Token))
						if err != nil {
							n.statusChan <- err.Error()

							continue
						}
					}

					chats, err := n.getChats(currentRoom.Token)
					if err != nil {
						if err.Error() == "invalid character '<' looking for beginning of value" {
							n.statusChan <- fmt.Sprintf(`Left room "%v" ("%v") with ID "%v" and token "%v"`, currentRoom.DisplayName, currentRoom.Name, currentRoom.ID, currentRoom.Token)

							return
						}

						n.statusChan <- err.Error()

						continue
					}

					chat := chats[0]
					if strconv.Itoa(chat.ID) != string(lastID) {
						n.chatChan <- chats[0]

						if err := n.knownIDs.Put([]byte(currentRoom.Token), []byte(strconv.Itoa(chat.ID))); err != nil {
							n.statusChan <- err.Error()
						}
					}

					time.Sleep(time.Second * 5)
				}
			}(room)
		}
	}()

	return nil
}

// CreateChat creates a chat message in a room
func (n *NextcloudTalk) CreateChat(room string, message string) error {
	client := resty.New()

	_, err := client.R().
		SetHeaders(map[string]string{
			"OCS-APIRequest": "true",
			"Accept":         "application/json",
		}).
		SetQueryParams(map[string]string{
			"message": message,
		}).
		SetBasicAuth(n.username, n.password).
		Post(n.url + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "chat", room))

	return err
}
