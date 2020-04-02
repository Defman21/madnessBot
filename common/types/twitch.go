package types

type TwitchHub struct {
	Callback     string `json:"hub.callback"`
	Mode         string `json:"hub.mode"`
	LeaseSeconds int    `json:"hub.lease_seconds"`
	Topic        string `json:"hub.topic"`
}

type TwitchStreamRequest struct {
	UserLogin string `json:"user_login"`
}

type TwitchStream struct {
	Title   string `json:"title"`
	Viewers int    `json:"viewer_count"`
	Game    string `json:"game_id"`
}

type TwitchStreamResponse struct {
	Data []TwitchStream `json:"data"`
}

type TwitchUserRequest struct {
	Login string `json:"login"`
}

type TwitchUser struct {
	ID string `json:"id"`
}

type TwitchUserResponse struct {
	Data []TwitchUser `json:"data"`
}

type TwitchGameRequest struct {
	ID string `json:"id"`
}

type TwitchGame struct {
	Name string `json:"name"`
}

type TwitchGameResponse struct {
	Data []TwitchGame `json:"data"`
}

type TwitchWebHookNotification struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Game    string `json:"game_id"`
	Viewers int    `json:"viewer_count"`
}

type TwitchWebHookNotificationRequest struct {
	Data []TwitchWebHookNotification `json:"data"`
}
