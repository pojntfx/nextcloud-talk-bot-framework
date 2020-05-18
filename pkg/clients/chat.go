package clients

// Chat is a chat message.
type Chat struct {
	ID                int    `json:"id"`
	Token             string `json:"token"`
	ActorType         string `json:"actorType"`          // guest, user
	ActorID           string `json:"actorId"`
	ActorDisplayName  string `json:"actorDisplayName"`
	IsReplyable       bool   `json:"isReplyable"`
	Message           string `json:"message"`
	MessageParameters string `json:"messageParamertes"` //RichObjectString
}

// ChatResponse is the API response for chats.
type ChatResponse struct {
	OCS struct {
		Data []Chat `json:"data"`
	} `json:"ocs"`
}
