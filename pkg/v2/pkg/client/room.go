package client

import (
	"encoding/json"
	"path"

	"github.com/go-resty/resty/v2"
)

// Room is a chat room
type Room struct {
	ID          int    `json:"id"`
	Token       string `json:"token"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// RoomResponse is the API response for rooms
type RoomResponse struct {
	OCS struct {
		Data []Room `json:"data"`
	} `json:"ocs"`
}

// GetRooms gets the chat rooms
func GetRooms(url, username, password string) ([]Room, error) {
	client := resty.New()

	res, err := client.R().
		SetHeaders(map[string]string{
			"OCS-APIRequest": "true",
			"Accept":         "application/json",
		}).
		SetBasicAuth(username, password).
		Get(url + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "room"))
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
