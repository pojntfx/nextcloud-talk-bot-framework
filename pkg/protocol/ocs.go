package protocol

type OCSRoom struct {
	OCS struct {
		Data []Room `json:"data"`
	} `json:"ocs"`
}

type OCSChat struct {
	OCS struct {
		Data []Chat `json:"chat"`
	} `json:"ocs"`
}
