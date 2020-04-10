package clients

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
