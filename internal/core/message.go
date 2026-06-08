package core

// Replyable is what Node.Reply needs from an incoming message's envelope.
// The request's msg_id is found separately via reflection on the Body field.
type Replyable interface {
	Source() string
	Destination() string
}

// Typed lets a reply body declare its protocol type string.
type Typed interface {
	ReplyType() string
}

// Incoming is the sealed set of message types ScanTyped can return. The
// unexported marker method means only types in this package can join the set,
// so a type switch over an Incoming has a known, closed list of cases.
type Incoming interface {
	isIncoming()
}

func (InitMessage) isIncoming() {}
func (EchoMessage) isIncoming() {}

// Envelope is the outer shell shared by every message type.
type Envelope struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

// Source and Destination implement the Replyable interface for any message
// that embeds Envelope.
func (e Envelope) Source() string      { return e.Src }
func (e Envelope) Destination() string { return e.Dest }

// BodyCommon holds the protocol fields every message body has.
type BodyCommon struct {
	Type      string `json:"type"`
	MsgID     int    `json:"msg_id"`
	InReplyTo int    `json:"in_reply_to,omitempty"`
}

// Init
type InitMessage struct {
	Envelope
	Body struct {
		BodyCommon
		NodeID  string   `json:"node_id"`
		NodeIDs []string `json:"node_ids"`
	} `json:"body"`
}

type InitOkBody struct {
	BodyCommon
}

func (InitOkBody) ReplyType() string { return "init_ok" }

// Echo
type EchoMessage struct {
	Envelope
	Body struct {
		BodyCommon
		Echo string `json:"echo"`
	} `json:"body"`
}

type EchoOkBody struct {
	BodyCommon
	Echo string `json:"echo"`
}

func (EchoOkBody) ReplyType() string { return "echo_ok" }
