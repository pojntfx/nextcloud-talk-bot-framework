package bots

import (
	"encoding/json"
	"log"
	"net/url"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/protocol"
)

// NextcloudTalk is a Nextcloud Talk bot.
type NextcloudTalk struct {
	URL                        *url.URL
	Username, Password         string
	rooms                      []protocol.Room
	newRoomChan                chan protocol.Room
	msgChan                    chan protocol.Message
	lastKnownMessages          cmap.ConcurrentMap
	lastKnownMessageTimestamps cmap.ConcurrentMap
}

// NewNextcloudTalk creates a new Nextcloud Talk bot.
func NewNextcloudTalk(url *url.URL, username, password string, msgChan chan protocol.Message) *NextcloudTalk {
	return &NextcloudTalk{
		url,
		username,
		password,
		make([]protocol.Room, 0),
		make(chan protocol.Room),
		msgChan,
		cmap.New(),
		cmap.New(),
	}
}

func (n *NextcloudTalk) ReadRooms() error {
	client := resty.New()

	for {
		res, err := client.R().SetBasicAuth(n.Username, n.Password).SetHeaders(map[string]string{
			"OCS-APIRequest": "true",
			"Accept":         "application/json",
		}).Get(n.URL.String() + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "room"))
		if err != nil {
			return err
		}

		var resStruct protocol.OCSRoom
		if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
			return err
		}

		for _, nroom := range resStruct.OCS.Data {
			exists := false
			for _, oroom := range n.rooms {
				if nroom.Token == oroom.Token {
					exists = true
					break
				}
			}

			if !exists {
				n.newRoomChan <- nroom
			}
		}

		n.rooms = resStruct.OCS.Data

		time.Sleep(time.Second * 10)
	}
}

func (n *NextcloudTalk) ReadMessages() error {
	for {
		room := <-n.newRoomChan

		go func(room2 protocol.Room) {
			lastId := room2.LastMessage.Id

			for {
				client := resty.New()
				res, err := client.R().SetBasicAuth(n.Username, n.Password).SetHeaders(map[string]string{
					"OCS-APIRequest": "true",
					"Accept":         "application/json",
				}).SetQueryParams(map[string]string{
					// "lastKnownMessageId": fmt.Sprintf("%v", lastId),
					"setReadMarker":  "false",
					"lookIntoFuture": "0",
				}).Get(n.URL.String() + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "chat", room2.Token))
				if err != nil {
					continue
				}

				var resStruct protocol.OCSMessage
				if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
					continue
				}

				newId := lastId
				var msg protocol.Message
				for _, smsg := range resStruct.OCS.Data {
					if smsg.Timestamp > lastId {
						newId = smsg.Id
						msg = smsg
					}
				}

				log.Println(lastId, "->", newId)

				if newId != lastId {
					n.msgChan <- msg
				}

				lastId = newId

				time.Sleep(time.Second * 5)
			}
		}(room)
	}
}
