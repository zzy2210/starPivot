package model

type Request struct {
	Messages string `json:"Messages"`
	ChatID   string `json:"ChatID"`
}

type Response struct {
	Messages string `json:"Messages"`
	ChatID   string `json:"ChatID"`
}
