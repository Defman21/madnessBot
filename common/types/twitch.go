package types

type TwitchHub struct {
	Callback     string `json:"hub.callback"`
	Mode         string `json:"hub.mode"`
	LeaseSeconds int    `json:"hub.lease_seconds"`
	Topic        string `json:"hub.topic"`
}
