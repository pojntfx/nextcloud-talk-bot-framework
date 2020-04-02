package client

import (
	"encoding/json"
	"path"

	"github.com/go-resty/resty/v2"
)

// Chat is a chat message
type Chat struct {
	ID               int    `json:"id"`
	Token            string `json:"token"`
	ActorID          string `json:"actorId"`
	ActorDisplayName string `json:"actorDisplayName"`
	Message          string `json:"message"`
}

// ChatResponse is the API response
type ChatResponse struct {
	OCS struct {
		Data []Chat `json:"data"`
	} `json:"ocs"`
}

// GetChats gets the chat messages
func GetChats(url, username, password, room string) ([]Chat, error) {
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
		SetBasicAuth(username, password).
		Get(url + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "chat", room))
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
