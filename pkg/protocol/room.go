package protocol

type Room struct {
	Token       string  `json:"token"`
	LastMessage Message `json:"lastMessage"`
}
