package client

// Chat is a chat message
type Chat struct {
	ID               int    `json:"id"`
	Token            string `json:"token"`
	ActorID          string `json:"actorId"`
	ActorDisplayName string `json:"actorDisplayName"`
	Message          string `json:"message"`
}

// ChatResponse is the API response for chats
type ChatResponse struct {
	OCS struct {
		Data []Chat `json:"data"`
	} `json:"ocs"`
}
