package models

// MessageParams represents generic message parameters.
// In real-world application it is better to avoid such types for better
// performance.
type MessageParams map[string]interface{}

type WsMessage struct {
	Params     MessageParams `json:"params"`
	ActionType string        `json:"method"`
	Data       string        `json:"data"`
	Body       string        `json:"body"`
	UserId     int           `json:"id"`
}

type WsResponse struct {
	UserId int           `json:"id"`
	Result MessageParams `json:"result"`
}

type WsError struct {
	UserId int           `json:"id"`
	Error  MessageParams `json:"error"`
}
