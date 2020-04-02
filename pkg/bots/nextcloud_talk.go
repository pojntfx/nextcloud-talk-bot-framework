package bots

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/pojntfx/nextcloud-talk-jitsi-bot/pkg/protocol"
)

// NextcloudTalk is a Nextcloud Talk bot.
type NextcloudTalk struct {
	URL                *url.URL
	Username, Password string
}

// NewNextcloudTalk creates a new Nextcloud Talk bot.
func NewNextcloudTalk(url *url.URL, username, password string) *NextcloudTalk {
	return &NextcloudTalk{
		url,
		username,
		password,
	}
}

// GetRooms gets the rooms of the user.
func (n *NextcloudTalk) GetRooms() ([]protocol.Room, error) {
	client := resty.New()

	res, err := client.R().SetBasicAuth(n.Username, n.Password).SetHeaders(map[string]string{
		"OCS-APIRequest": "true",
		"Accept":         "application/json",
	}).Get(n.URL.String() + "/" + path.Join("ocs", "v2.php", "apps", "spreed", "api", "v1", "room"))
	if err != nil {
		return []protocol.Room{}, err
	}

	var resStruct protocol.OCSRoom
	if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
		return []protocol.Room{}, err
	}

	return resStruct.OCS.Data, nil
}
