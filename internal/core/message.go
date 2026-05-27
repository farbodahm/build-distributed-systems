package core

type Message struct {
	Src  string         `json:"src"`
	Dest string         `json:"dest"`
	Body map[string]interface{} `json:"body"`
}
