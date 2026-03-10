package server

type Message struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}
