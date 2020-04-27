package clients

// Room is a chat room
type Room struct {
	ID               int    `json:"id"`
	Token            string `json:"token"`
	Type             int    `json:"type"`
	Name             string `json:"name"`
	DisplayName      string `json:"displayName"`
	ParticipantType  int    `json:"participantType"`
	ParticipantFlags int    `json:"participantFlags"`
	ReadOnly         int    `json:"readOnly"`
	HasPassword      bool   `json:"hasPassword"`
	LastActivity     int    `json:"lastActivity"`
	NotificatonLevel int    `json:"notificatonLevel"` // Participant::NOTIFY_* (1-3)
	LastReadMessage  int    `json:"lastReadMessage"`
	LastMessage      int    `json:"lastMessage"`
}

// RoomResponse is the API response for rooms
type RoomResponse struct {
	OCS struct {
		Data []Room `json:"data"`
	} `json:"ocs"`
}
