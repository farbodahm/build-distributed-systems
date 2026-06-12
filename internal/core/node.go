package core

import (
	"reflect"
	"sync"
)

// Handler processes one incoming message and may reply via Node.Reply.
type Handler func(msg Incoming) error

// Node holds per-process protocol state.
type Node struct {
	ID    string
	Peers []string

	mu        sync.Mutex
	nextMsgID int
	handlers  map[IncomingMessageType]Handler
}

func NewNode() *Node {
	return &Node{
		handlers: make(map[IncomingMessageType]Handler),
	}
}

// Init records the node's identity and peers (typically from an init body).
func (n *Node) Init(id string, peers []string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ID = id
	n.Peers = peers
}

// NextMsgID returns the next outgoing msg_id and bumps the counter.
func (n *Node) NextMsgID() int {
	n.mu.Lock()
	defer n.mu.Unlock()

	id := n.nextMsgID
	n.nextMsgID++
	return id
}

// RegisterHandler registers a handler for incoming messages of type t.
// The handler can reply to the message using Node.Reply.
func (n *Node) RegisterHandler(t IncomingMessageType, handler Handler) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.handlers[t]; ok {
		panic("handler already registered for " + t)
	}
	n.handlers[t] = handler
}

// OnInit registers a typed handler for init messages.
func (n *Node) OnInit(h func(InitMessage) error) {
	n.RegisterHandler(MsgTypeInit, func(msg Incoming) error { return h(msg.(InitMessage)) })
}

// OnEcho registers a typed handler for echo messages.
func (n *Node) OnEcho(h func(EchoMessage) error) {
	n.RegisterHandler(MsgTypeEcho, func(msg Incoming) error { return h(msg.(EchoMessage)) })
}

// Run reads messages from stdin and dispatches each to its registered handler,
// one at a time, until stdin closes.
func (n *Node) Run() error {
	Log.Info("starting node")

	var msg Incoming
	for ScanTyped(&msg) {
		t := msg.Type()
		handler, ok := n.handlers[t]
		if !ok {
			Log.Warn("no handler registered for message type %q", t)
			continue
		}
		if err := handler(msg); err != nil {
			Log.Error("handler %q: %v", t, err)
		}
	}

	return nil
}

// Reply sends body as a JSON reply to req. Src/Dest are swapped. The body's
// MsgID is set to a fresh counter value and InReplyTo is copied from the
// request's Body.MsgID. Field lookups go through reflection, so body just needs
// to embed BodyCommon (or otherwise expose MsgID/InReplyTo as exported fields)
// and req's body needs to expose MsgID.
//
// body must be a pointer so reflection can write to its fields.
func (n *Node) Reply(req Replyable, body interface{}) {
	msgID := n.NextMsgID()
	reqMsgID := req.RequestMsgID()

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
