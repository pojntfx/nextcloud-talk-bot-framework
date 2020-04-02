package protocol

type OCSRoom struct {
	OCS struct {
		Data []Room `json:"data"`
	} `json:"ocs"`
}

type OCSMessage struct {
	OCS struct {
		Data []Message `json:"data"`
	} `json:"ocs"`
}
