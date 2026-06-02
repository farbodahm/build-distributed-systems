package core

import "reflect"

// Node holds per-process protocol state.
type Node struct {
	ID    string
	Peers []string

	nextMsgID int
}

func NewNode() *Node {
	return &Node{}
}

// Init records the node's identity and peers (typically from an init body).
func (n *Node) Init(id string, peers []string) {
	n.ID = id
	n.Peers = peers
}

// NextMsgID returns the next outgoing msg_id and bumps the counter.
func (n *Node) NextMsgID() int {
	id := n.nextMsgID
	n.nextMsgID++
	return id
}

// Reply sends body as a JSON reply to req. Src/Dest are swapped. The body's
// MsgID is set to a fresh counter value and InReplyTo is copied from the
// request's Body.MsgID. Field lookups go through reflection, so body just
// needs to embed BodyCommon (or otherwise expose MsgID/InReplyTo as exported
// fields) and req's body needs to expose MsgID.
//
// body must be a pointer so reflection can write to its fields.
func (n *Node) Reply(req Replyable, body interface{}) {
	msgID := n.NextMsgID()
	reqMsgID := requestMsgID(req)

	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if f := v.FieldByName("MsgID"); f.IsValid() && f.CanSet() {
		f.SetInt(int64(msgID))
	}
	if f := v.FieldByName("InReplyTo"); f.IsValid() && f.CanSet() {
		f.SetInt(int64(reqMsgID))
	}
	if t, ok := body.(Typed); ok {
		if f := v.FieldByName("Type"); f.IsValid() && f.CanSet() {
			f.SetString(t.ReplyType())
		}
	}

	Log.PrintJSON(struct {
		Src  string      `json:"src"`
		Dest string      `json:"dest"`
		Body interface{} `json:"body"`
	}{req.Destination(), req.Source(), body})
}

// requestMsgID reads req.Body.MsgID via reflection. Returns 0 if either field
// is missing.
func requestMsgID(req interface{}) int {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	body := v.FieldByName("Body")
	if !body.IsValid() {
		return 0
	}
	f := body.FieldByName("MsgID")
	if !f.IsValid() {
		return 0
	}
	return int(f.Int())
}
