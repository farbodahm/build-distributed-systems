package core

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

type InitMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type    string   `json:"type"`
		NodeID  string   `json:"node_id"`
		NodeIDs []string `json:"node_ids"`
	} `json:"body"`
}
