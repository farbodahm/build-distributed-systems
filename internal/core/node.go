package core

import (
	"reflect"
	"sync"
)

// Node holds per-process protocol state.
type Node struct {
	ID    string
	Peers []string

	mu        sync.Mutex
	nextMsgID int
	handlers  map[IncomingMessageType]func(Incoming) error
}

func NewNode() *Node {
	return &Node{
		handlers: make(map[IncomingMessageType]func(Incoming) error),
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
func (n *Node) RegisterHandler(t IncomingMessageType, handler func(Incoming) error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.handlers[t]; ok {
		panic("handler already registered for " + t)
	}
	n.handlers[t] = handler
}

// Run reads messages from stdin and dispatches each to its registered handler,
// running every handler in its own goroutine. It returns after stdin closes
// and all in-flight handlers finish.
func (n *Node) Run() error {
	Log.Info("starting node")
	var wg sync.WaitGroup

	var msg Incoming
	for ScanTyped(&msg) {
		t := msg.Type()

		n.mu.Lock()
		handler, ok := n.handlers[t]
		n.mu.Unlock()
		if !ok {
			Log.Warn("no handler registered for message type %q", t)
			continue
		}

		wg.Add(1)
		go func(m Incoming, t IncomingMessageType, h func(Incoming) error) {
			defer wg.Done()
			if err := h(m); err != nil {
				Log.Error("handler %q: %v", t, err)
			}
		}(msg, t, handler)
	}

	wg.Wait()
	return nil
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
