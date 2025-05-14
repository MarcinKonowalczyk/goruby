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
	Trace(args ...string) Exit
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

func debugf(format string, args ...any) {
	if dEBUG {
		fmt.Printf(format, args...)
	}
}

//////////

type FunctionName string

const (
	START_NODE FunctionName = "<START>"
	END_NODE   FunctionName = "<END>"
)

func (f FunctionName) IsUntraced() bool {
	return strings.HasPrefix(string(f), "*")
}

type Message struct {
	Message string
	Stack   []FunctionName
}

func (m Message) String() string {
	var out strings.Builder
	for j := len(m.Stack) - 1; j >= 0; j-- {
		out.WriteString(fmt.Sprintf("%s:", m.Stack[j]))
	}
	out.WriteString(fmt.Sprintf(" %s", m.Message))
	return out.String()
}

type linked interface {
	Next() Node
	Prev() Node
}
type Node interface {
	linked
	Name() FunctionName
	Message(args ...any)
	Messages() []string
}

type Enter struct {
	name     FunctionName
	messages []string
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
	if len(n.messages) > 0 {
		out.WriteString(fmt.Sprintf(" (#%d)", len(n.messages)))
	}
	return out.String()
}

func (n *Enter) Message(args ...any) {
	if n.messages == nil {
		n.messages = make([]string, 0)
	}
	n.messages = append(n.messages, fmt.Sprint(args...))
}

func (n Enter) Messages() []string {
	return n.messages
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

type tracer struct {
	stack []Node
	// pointer to the enter node of the current function
	where Node
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
			n.parent = t.where.(*Enter)
			t.where = n // and update the where pointer
		case *Exit:
			t.where = n.parent.parent
		default:
			panic("unknown node type")
		}
		// actually set the next node
		t.stack = append(t.stack, n)
	}
}

func callerName(N int) FunctionName {
	parent, _, _, _ := runtime.Caller(N + 1)
	info := runtime.FuncForPC(parent)
	name := info.Name()
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

func (t *tracer) Trace(args ...string) Exit {
	var callers_name FunctionName
	switch len(args) {
	case 0:
		callers_name = callerName(1)
	case 1:
		callers_name = FunctionName(args[0])
	default:
		panic("too many arguments to Trace")
	}
	n := &Enter{name: callers_name}
	t.append(n)
	debug("> entering", t.where.Name())
	return Exit{
		name:   callers_name,
		parent: n,
	}
}

// Usage pattern: defer t.Un(t.Trace(p, "..."))
func (t *tracer) Un(exit Exit) {
	t.append(&exit)
	debug("< exiting", t.where.Name())
}

func (t *tracer) Message(args ...any) {
	debug("  messaging", t.where.Name())
	callers_name := callerName(1)
	if len(t.stack) == 0 {
		panic("call stack is empty")
	}
	// var top_node tracer_node = t.chain[len(t.chain)-1]
	if t.where.Name() == callers_name {
		// we are tracing the current function. Pass it the message
		t.where.Message(args...)
	} else {
		debug("!", t.where.Name(), callers_name)

		callers_callers_name := callerName(2)
		if t.where.Name() == callers_callers_name {
			// we are tracing the parent function. we can recover one-deep mishap
			n0 := &Enter{name: "*" + callers_name}
			n1 := &Exit{name: "*" + callers_name, parent: n0}
			n0.Message(args...) // pass the message
			t.append(n0, n1)
		} else {
			// we are not tracing the current function, nor the parent function
			n0 := &Enter{name: "..."}
			n1 := &Enter{name: "*" + callers_callers_name}
			n2 := &Enter{name: "*" + callers_name}
			n3 := &Exit{name: "*" + callers_name, parent: n2}
			n4 := &Exit{name: "*" + callers_callers_name, parent: n1}
			n5 := &Exit{name: "...", parent: n0}
			n2.Message(args...) // pass the message
			t.stack = append(t.stack, n0, n1, n2, n3, n4, n5)
		}
	}
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

var (
	_ settable_next = (*Enter)(nil)
	_ settable_next = (*Exit)(nil)
	_ settable_prev = (*Enter)(nil)
	_ settable_prev = (*Exit)(nil)
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
	stack := make([]*Enter, 0)
	walkable, err := t.ToWalkable()
	if err != nil {
		return nil
	}

	walkable.Walk(func(n Node) error {
		switch n := n.(type) {
		case *Enter:
			if n.Name() == START_NODE {
				// skip the first node
			} else {
				stack = append(stack, n)
			}
		case *Exit:
			if len(stack) == 0 {
				return fmt.Errorf("unmatched exit node: %s", n.Name())
			}
			if stack[len(stack)-1].Name() != n.Name() {
				return fmt.Errorf("unmatched exit node: %s", n.Name())
			}
			stack = stack[:len(stack)-1]
		default:
			return fmt.Errorf("unknown node type: %T", n)
		}
		if n.Name() == END_NODE {
			// skip the last node
			return nil
		}

		node_messages := n.Messages()
		if node_messages == nil {
			return nil
		}

		stack_copy := make([]FunctionName, len(stack))
		for i := 0; i < len(stack); i++ {
			stack_copy[i] = stack[i].Name()
		}
		for _, m := range node_messages {
			messages = append(messages, Message{
				Message: m,
				Stack:   stack_copy,
			})
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
