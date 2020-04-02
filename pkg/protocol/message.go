package protocol

type Message struct {
	Id               int    `json:"id"`
	Timestamp        int    `json:"timestamp"`
	ActorId          string `json:"actorId"`
	ActorDisplayName string `json:"actorDisplayName"`
	Content          string `json:"message"`
}
