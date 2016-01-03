package message

type Message struct {
	Address *Address `json:"address",omitempty`
	Body    string   `json:"body"`
}

type Address struct {
	Gateway string `json:"gateway"`
	Channel string `json:"channel"`
}
