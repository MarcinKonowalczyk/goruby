package trace

import (
	"fmt"
	"runtime"
	"strings"
)

type walkable interface {
	// WalkEnter(fn enter_func) error
	Walk(fn func(node Node) error) error
}

type Tracer interface {
	// add the function to the tracing stack
	Trace(where FunctionName) Exit
	Un(Exit)
	// add a message to the trace
	Message(args ...any)
	// get all the messages
	Messages() []Message
	// mark the end of the trace
	Done()
	//
	ToWalkable() (walkable, error)
}

func NewTracer() Tracer {
	chain := make([]Node, 1)
	chain[0] = &Enter{
		name: FunctionName(START_NODE),
	}
	return &tracer{
		stack: chain,
		where: chain[0].(*Enter),
	}

}

// local debug flag

const dEBUG = false

// const dEBUG = true

func debug(args ...any) {
	if dEBUG {
		fmt.Println(args...)
	}
}

// func debugf(format string, args ...any) {
// 	if dEBUG {
// 		fmt.Printf(format, args...)
// 	}
// }

//////////

type FunctionName string

const (
	START_NODE FunctionName = "<START>"
	END_NODE   FunctionName = "<END>"
)

func (f FunctionName) IsUntraced() bool {
	return strings.HasPrefix(string(f), "*")
}

type linked interface {
	Next() Node
	Prev() Node
}

type Node interface {
	linked
	Name() FunctionName
	// Message(args ...any)
	// Messages() []string
}

type Enter struct {
	name FunctionName
	// next/prev nodes in the chain
	next Node
	prev Node
	// parent node -- what we are entering from
	parent *Enter
}

func (n Enter) Name() FunctionName {
	return n.name
}

func (n Enter) String() string {
	var out strings.Builder
	if n.parent != nil {
		out.WriteString(string(n.parent.Name()))
		out.WriteString(" -> ")
	}
	out.WriteString(string(n.name))
	return out.String()
}

func (n *Enter) Next() Node {
	return n.next
}

func (n *Enter) Prev() Node {
	return n.prev
}

var (
	_ Node   = (*Enter)(nil)
	_ linked = (*Enter)(nil)
)

type Exit struct {
	name FunctionName
	// next/prev nodes in the chain
	next Node
	prev Node
	// parent node -- what we are exiting from
	parent *Enter
}

func (n *Exit) Name() FunctionName {
	return n.name
}

func (n *Exit) Message(args ...any) {
	panic("exit node should not receive messages")
}

func (n *Exit) Messages() []string {
	return nil
}

func (n *Exit) String() string {
	return string(n.name)
}

func (n *Exit) Next() Node {
	return n.next
}

func (n *Exit) Prev() Node {
	return n.prev
}

var (
	_ Node   = (*Exit)(nil)
	_ linked = (*Exit)(nil)
)

type Message struct {
	Message string
	// next/prev nodes in the chain
	next Node
	prev Node
	// parent node -- the enter node at which we are messaging
	parent *Enter
}

func (m Message) String() string {
	var out strings.Builder
	stack := m.Stack()
	for j := len(stack) - 1; j >= 0; j-- {
		out.WriteString(fmt.Sprintf("%s:", stack[j]))
	}
	out.WriteString(fmt.Sprintf(" %s", m.Message))
	return out.String()
}

func (m *Message) Stack() []FunctionName {
	stack := make([]FunctionName, 0)
	for n := m.parent; n != nil; n = n.parent {
		if n.Name() == START_NODE {
			break
		}
		stack = append(stack, n.Name())
	}
	// reverse the stack
	// for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
	// 	stack[i], stack[j] = stack[j], stack[i]
	// }
	return stack
}

func (m *Message) Name() FunctionName {
	return m.parent.Name() + " message"
}

func (m *Message) Next() Node {
	return m.next
}
func (m *Message) Prev() Node {
	return m.prev
}

var (
	_ Node   = (*Message)(nil)
	_ linked = (*Message)(nil)
)

type tracer struct {
	stack []Node
	// pointer to the enter node of the current function
	where *Enter
}

// Append any number of nodes to the chain
func (t *tracer) append(node ...Node) {
	for _, n := range node {
		// link new node with the top of the stack
		top := t.stack[len(t.stack)-1]
		top.(settable_next).SetNext(n)
		n.(settable_prev).SetPrev(top)
		// if we've added an new enter node, set the parent
		switch n := n.(type) {
		case *Enter:
			n.parent = t.where
			t.where = n // and update the where pointer
		case *Exit:
			t.where = n.parent.parent
		case *Message:
			// nothing to do for message nodes
		default:
			panic("unknown node type")
		}
		// actually set the next node
		t.stack = append(t.stack, n)
	}
}

func callerName(N int) string {
	parent, _, _, _ := runtime.Caller(N + 1)
	info := runtime.FuncForPC(parent)
	name := info.Name()
	return name
}

func Here() FunctionName {
	name := callerName(1)
	// strip everything before the last . to get just the function name
	name = name[strings.LastIndex(name, ".")+1:]
	return FunctionName(name)
}

// passes any information to the corresponding Un
// type trace_out struct {
// 	callers_name FunctionName
// 	parent       *Enter
// }

// func random_node_id() node_id {
// 	return node_id(rand.Int63())
// }

func (t *tracer) Trace(where FunctionName) Exit {
	n := &Enter{name: where}
	t.append(n)
	debug("> entering", t.where.Name())
	return Exit{
		name:   where,
		parent: n,
	}
}

// Usage pattern: defer t.Un(t.Trace(p, "..."))
func (t *tracer) Un(exit Exit) {
	t.append(&exit)
	debug("< exiting", t.where.Name())
}

func (t *tracer) Message(args ...any) {
	t.append(&Message{
		Message: fmt.Sprint(args...),
		parent:  t.where,
	})
}

type settable_next interface {
	SetNext(Node)
}

type settable_prev interface {
	SetPrev(Node)
}

func (t *Enter) SetNext(n Node) {
	t.next = n
}

func (t *Enter) SetPrev(n Node) {
	t.prev = n
}

func (t *Exit) SetNext(n Node) {
	t.next = n
}

func (t *Exit) SetPrev(n Node) {
	t.prev = n
}

func (t *Message) SetNext(n Node) {
	t.next = n
}
func (t *Message) SetPrev(n Node) {
	t.prev = n
}

var (
	_ settable_next = (*Enter)(nil)
	_ settable_next = (*Exit)(nil)
	_ settable_prev = (*Enter)(nil)
	_ settable_prev = (*Exit)(nil)
	_ settable_next = (*Message)(nil)
	_ settable_prev = (*Message)(nil)
)

func (t *tracer) Done() {
	t.stack = append(t.stack, &Exit{
		name: END_NODE,
	})

	// link the nodes together
	for i := 0; i < len(t.stack)-1; i++ {
		t.stack[i].(settable_next).SetNext(t.stack[i+1])
	}

	for i := len(t.stack) - 1; i > 0; i-- {
		t.stack[i].(settable_prev).SetPrev(t.stack[i-1])
	}
}

func (t *tracer) ToWalkable() (walkable, error) {
	if t.stack == nil {
		panic("call stack is empty")
	}
	// check we're done with the tracing
	switch n := t.stack[len(t.stack)-1].(type) {
	case *Enter:
		return nil, fmt.Errorf("not walkable. tracer is not done. last node was an enter node: %s", n.Name())
	case *Exit:
		if n.Name() != END_NODE {
			return nil, fmt.Errorf("not walkable. tracer is not done. last node was not an exit node of the root node: %s", n.Name())
		}
	}
	return t, nil
}

func (t *tracer) Messages() []Message {
	messages := make([]Message, 0)
	// stack := make([]*Enter, 0)
	walkable, err := t.ToWalkable()
	if err != nil {
		return nil
	}

	walkable.Walk(func(n Node) error {
		switch n := n.(type) {
		case *Message:
			messages = append(messages, *n)
		default:
			// do nothing
		}
		return nil
	})
	return messages
}

func (t *tracer) Walk(fn func(node Node) error) error {
	var err error
	for i := 0; i < len(t.stack); i++ {
		node := t.stack[i]
		err = fn(node)
		if err != nil {
			return fmt.Errorf("error in walk function: %w", err)
		}
	}
	return nil
}

var (
	_ Tracer   = (*tracer)(nil)
	_ walkable = (*tracer)(nil)
)
