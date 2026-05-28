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
		MsgID   int      `json:"msg_id"`
		NodeID  string   `json:"node_id"`
		NodeIDs []string `json:"node_ids"`
	} `json:"body"`
}

type InitOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		InReplyTo int    `json:"in_reply_to"`
		MsgID     int    `json:"msg_id"`
	} `json:"body"`
}
