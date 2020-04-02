package main

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
)

type ChatResponse struct {
	OCS struct {
		Data []struct {
			ID               int    `json:"id"`
			Token            string `json:"token"`
			ActorID          string `json:"actorId"`
			ActorDisplayName string `json:"actorDisplayName"`
			Message          string `json:"message"`
		} `json:"data"`
	} `json:"ocs"`
}

func main() {
	var (
		url      = os.Getenv("NEXTCLOUD_URL")
		username = os.Getenv("NEXTCLOUD_USERNAME")
		password = os.Getenv("NEXTCLOUD_PASSWORD")
		room     = os.Getenv("NEXTCLOUD_ROOM")
	)

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
		log.Fatal(err)
	}

	var resStruct ChatResponse
	if err := json.Unmarshal(res.Body(), &resStruct); err != nil {
		log.Fatal(err)
	}

	log.Println(resStruct.OCS.Data[0])
}
