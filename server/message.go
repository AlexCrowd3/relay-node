package server

type Message struct {
	Type string                 `json:"type"`
	From string                 `json:"from"`
	To   string                 `json:"to"`
	Data map[string]interface{} `json:"data"`
}
