package ast

import (
	"fmt"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast/infix"
)

// Node represents a node within the AST
//
// All node types implement the Node interface.
type Node interface {
	node() // marks this as a node
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

////////////////////////////////////////////////////////////////////////////////

// A Program node is the root node within the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}

func (p *Program) String() string {
	stmts := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		if s != nil {
			stmts[i] = s.String()
		}
	}
	return strings.Join(stmts, "\n")
}

var (
	_ Node = &Program{}
)

// A ReturnStatement represents a return node which yields another Expression.
type ReturnStatement struct {
	ReturnValue Expression
}

func (rs *ReturnStatement) node()          {}
func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

var (
	_ Node      = &ReturnStatement{}
	_ Statement = &ReturnStatement{}
)

// An ExpressionStatement is a Statement wrapping an Expression
type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) node()          {}
func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

var (
	_ Node      = &ExpressionStatement{}
	_ Statement = &ExpressionStatement{}
)

// BlockStatement represents a list of statements
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) node()          {}
func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) String() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		if s != nil {
			out.WriteString(s.String())
		}
	}
	return out.String()
}

var (
	_ Node      = &BlockStatement{}
	_ Statement = &BlockStatement{}
)

// A BreakStatement represents a break statement
type BreakStatement struct {
	Condition Expression
	Unless    bool
}

func (bs *BreakStatement) node()          {}
func (bs *BreakStatement) statementNode() {}

func (bs *BreakStatement) String() string {
	var out strings.Builder
	out.WriteString("break")
	out.WriteString(" ")
	out.WriteString(bs.Condition.String())
	return out.String()
}

var (
	_ Node      = &BreakStatement{}
	_ Statement = &BreakStatement{}
)

// Assignment represents a generic assignment
type Assignment struct {
	Left  Expression
	Right Expression
}

func (a *Assignment) node()           {}
func (a *Assignment) expressionNode() {}

func (a *Assignment) String() string {
	var out strings.Builder
	out.WriteString(encloseInParensIfNeeded(a.Left))
	out.WriteString(" = ")
	out.WriteString(encloseInParensIfNeeded(a.Right))
	return out.String()
}

var (
	_ Node       = &Assignment{}
	_ Expression = &Assignment{}
)

// MultiAssignment represents multiple variables on the left-hand side
type MultiAssignment struct {
	Variables []*Identifier
	Values    []Expression
}

func (m *MultiAssignment) node()           {}
func (m *MultiAssignment) expressionNode() {}

func (m *MultiAssignment) String() string {
	var out strings.Builder
	vars := make([]string, len(m.Variables))
	for i, v := range m.Variables {
		vars[i] = v.Value
	}
	out.WriteString(strings.Join(vars, ", "))
	out.WriteString(" = ")
	values := make([]string, len(m.Values))
	for i, v := range m.Values {
		values[i] = v.String()
	}
	out.WriteString(strings.Join(values, ", "))
	return out.String()
}

var (
	_ Node       = &MultiAssignment{}
	_ Expression = &MultiAssignment{}
)

// Keyword__FILE__ represents __FILE__ in the AST
type Keyword__FILE__ struct {
	Filename string
}

func (f *Keyword__FILE__) node()           {}
func (f *Keyword__FILE__) expressionNode() {}
func (f *Keyword__FILE__) String() string  { return "__FILE__" }

var (
	_ Node       = &Keyword__FILE__{}
	_ Expression = &Keyword__FILE__{}
)

// An Identifier represents an identifier in the program
type Identifier struct {
	Constant bool // true if the identifier is a constant
	Value    string
}

func (i *Identifier) node()           {}
func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string { return i.Value }

var (
	_ Node       = &Identifier{}
	_ Expression = &Identifier{}
)

// IsConstant returns true if the Identifier represents a Constant, false otherwise
// func (i *Identifier) IsConstant() bool { return i.Token.Type == token.CONST }

// Global represents a global in the AST
type Global struct {
	Value string
}

func (g *Global) node()           {}
func (g *Global) expressionNode() {}

func (g *Global) String() string { return g.Value }

var (
	_ Node       = &Global{}
	_ Expression = &Global{}
)

// IntegerLiteral represents an integer in the AST
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) node()           {}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }

var (
	_ Node       = &IntegerLiteral{}
	_ Expression = &IntegerLiteral{}
)

// FloatLiteral represents a float in the AST
type FloatLiteral struct {
	Value float64
}

func (fl *FloatLiteral) node()           {}
func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) String() string { return fmt.Sprintf("%f", fl.Value) }

var (
	_ Node       = &FloatLiteral{}
	_ Expression = &FloatLiteral{}
)

// Nil represents the 'nil' keyword
type Nil struct{}

func (n *Nil) node()           {}
func (n *Nil) expressionNode() {}

func (n *Nil) String() string { return "nil" }

var (
	_ Node       = &Nil{}
	_ Expression = &Nil{}
)

// Boolean represents a boolean in the AST
type Boolean struct {
	Value bool
}

func (b *Boolean) node()           {}
func (b *Boolean) expressionNode() {}

func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

var (
	_ Node       = &Boolean{}
	_ Expression = &Boolean{}
)

// StringLiteral represents a double quoted string in the AST
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) node()           {}
func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string { return sl.Value }

var (
	_ Node       = &StringLiteral{}
	_ Expression = &StringLiteral{}
)

// Comment represents a double quoted string in the AST
type Comment struct {
	Value string
}

func (c *Comment) node()          {}
func (c *Comment) statementNode() {}

func (c *Comment) String() string { return "#" + c.Value }

var (
	_ Node      = &Comment{}
	_ Statement = &Comment{}
)

// SymbolLiteral represents a symbol within the AST
type SymbolLiteral struct {
	Value Expression
}

func (s *SymbolLiteral) node()           {}
func (s *SymbolLiteral) expressionNode() {}

func (s *SymbolLiteral) String() string { return ":" + s.Value.String() }

var (
	_ Node       = &SymbolLiteral{}
	_ Expression = &SymbolLiteral{}
)

// ConditionalExpression represents an if expression within the AST
type ConditionalExpression struct {
	Unless      bool // true = unless, false = if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ce *ConditionalExpression) node()           {}
func (ce *ConditionalExpression) expressionNode() {}

func (ce *ConditionalExpression) String() string {
	var out strings.Builder
	if ce.Unless {
		out.WriteString("unless ")
	} else {
		out.WriteString("if ")
	}
	out.WriteString(ce.Condition.String())
	out.WriteString(" ")
	out.WriteString(ce.Consequence.String())
	if ce.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ce.Alternative.String())
	}
	out.WriteString(" end")
	return out.String()
}

var (
	_ Node       = &ConditionalExpression{}
	_ Expression = &ConditionalExpression{}
)

// A LoopExpression represents a loop
type LoopExpression struct {
	Condition Expression
	Block     *BlockStatement
}

func (ce *LoopExpression) node()           {}
func (ce *LoopExpression) expressionNode() {}

func (ce *LoopExpression) String() string {
	var out strings.Builder
	out.WriteString("while ")
	out.WriteString(ce.Condition.String())
	out.WriteString(" do ")
	out.WriteString(ce.Block.String())
	out.WriteString(" end")
	return out.String()
}

var (
	_ Node       = &LoopExpression{}
	_ Expression = &LoopExpression{}
)

// ExpressionList represents a list of expressions within the AST divided by commas
type ExpressionList []Expression

func (el ExpressionList) node()           {}
func (el ExpressionList) expressionNode() {}

func (el ExpressionList) String() string {
	var out strings.Builder
	elements := []string{}
	for _, e := range el {
		elements = append(elements, e.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	return out.String()
}

var (
	_ Node       = ExpressionList{}
	_ Expression = ExpressionList{}
)

// ArrayLiteral represents an Array literal within the AST
type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) node()           {}
func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

var (
	_ Node       = &ArrayLiteral{}
	_ Expression = &ArrayLiteral{}
)

// HashLiteral represents an Hash literal within the AST
type HashLiteral struct {
	Map map[Expression]Expression
}

func (hl *HashLiteral) node()           {}
func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for key, val := range hl.Map {
		elements = append(elements, fmt.Sprintf("%q => %q", key.String(), val.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}

var (
	_ Node       = &HashLiteral{}
	_ Expression = &HashLiteral{}
)

// RangeLiteral represents a range literal within the AST
type RangeLiteral struct {
	Left      Expression
	Right     Expression
	Inclusive bool
}

func (rl *RangeLiteral) node()           {}
func (rl *RangeLiteral) expressionNode() {}

// String returns the string representation of the range
func (rl *RangeLiteral) String() string {
	var out strings.Builder
	out.WriteString(rl.Left.String())
	out.WriteString(" ")
	if rl.Inclusive {
		out.WriteString("..")
	} else {
		out.WriteString("...")
	}
	out.WriteString(" ")
	out.WriteString(rl.Right.String())
	return out.String()
}

var (
	_ Node       = &RangeLiteral{}
	_ Expression = &RangeLiteral{}
)

// A BlockCapture represents a function scoped variable capturing a block
type BlockCapture struct {
	Name *Identifier
}

func (b *BlockCapture) node()           {}
func (b *BlockCapture) expressionNode() {}

func (b *BlockCapture) String() string {
	return "&" + b.Name.Value
}

var (
	_ Node       = &BlockCapture{}
	_ Expression = &BlockCapture{}
)

// A FunctionLiteral represents a function definition in the AST
type FunctionLiteral struct {
	Name       *Identifier
	Parameters []*FunctionParameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) node()           {}
func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) String() string {
	var out strings.Builder
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("def ")
	out.WriteString(fl.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if fl.Body != nil {
		out.WriteString(fl.Body.String())
	}
	out.WriteString(" end")
	return out.String()
}

var (
	_ Node       = &FunctionLiteral{}
	_ Expression = &FunctionLiteral{}
)

// A FunctionParameter represents a parameter in a function literal
type FunctionParameter struct {
	Name    *Identifier
	Default Expression
	Splat   bool
}

func (f *FunctionParameter) node()           {}
func (f *FunctionParameter) expressionNode() {}

func (f *FunctionParameter) String() string {
	var out strings.Builder
	if f.Splat {
		out.WriteString("*")
	}
	out.WriteString(f.Name.String())
	if f.Default != nil {
		out.WriteString(" = ")
		out.WriteString(encloseInParensIfNeeded(f.Default))
	}
	return out.String()
}

var (
	_ Node       = &FunctionParameter{}
	_ Expression = &FunctionParameter{}
)

// A Splat represents a splat operator in the AST
type Splat struct {
	Value Expression
}

func (s *Splat) node()           {}
func (s *Splat) expressionNode() {}

func (s *Splat) String() string {
	var out strings.Builder
	out.WriteString("*")
	out.WriteString(s.Value.String())
	return out.String()
}

var (
	_ Node       = &Splat{}
	_ Expression = &Splat{}
)

// An IndexExpression represents an array or hash access in the AST
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) node()           {}
func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

var (
	_ Node       = &IndexExpression{}
	_ Expression = &IndexExpression{}
)

// A ContextCallExpression represents a method call on a given Context
type ContextCallExpression struct {
	Context   Expression       // The lefthandside expression
	Function  *Identifier      // The function to call
	Arguments []Expression     // The function arguments
	Block     *BlockExpression // The function block
}

func (ce *ContextCallExpression) node()           {}
func (ce *ContextCallExpression) expressionNode() {}

func (ce *ContextCallExpression) String() string {
	var out strings.Builder
	if ce.Context != nil {
		out.WriteString(ce.Context.String())
		out.WriteString(".")
	}
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	if ce.Block != nil {
		out.WriteString("\n")
		out.WriteString(ce.Block.String())
	}
	return out.String()
}

var (
	_ Node       = &ContextCallExpression{}
	_ Expression = &ContextCallExpression{}
)

// A BlockExpression represents a Ruby block
type BlockExpression struct {
	Parameters []*FunctionParameter // the block parameters
	Body       *BlockStatement      // the block body
}

func (b *BlockExpression) node()           {}
func (b *BlockExpression) expressionNode() {}

// String returns a string representation of the block statement
func (b *BlockExpression) String() string {
	var out strings.Builder
	out.WriteString("{")
	if len(b.Parameters) != 0 {
		args := []string{}
		for _, a := range b.Parameters {
			args = append(args, a.String())
		}
		out.WriteString("|")
		out.WriteString(strings.Join(args, ", "))
		out.WriteString("|")
		out.WriteString("\n")
	}
	out.WriteString(b.Body.String())
	out.WriteString("\n")
	out.WriteString("}")
	return out.String()
}

var (
	_ Node       = &BlockExpression{}
	_ Expression = &BlockExpression{}
)

// PrefixExpression represents a prefix operator
type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) node()           {}
func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// An InfixExpression represents an infix operator in the AST
type InfixExpression struct {
	Left     Expression
	Operator infix.Infix
	Right    Expression
}

func (oe *InfixExpression) node()           {}
func (oe *InfixExpression) expressionNode() {}

func (oe *InfixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator.String() + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

func encloseInParensIfNeeded(expr Expression) string {
	val := expr.String()
	return val
	// hasParens := strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")")
	// if !hasParens {
	// 	val = "(" + val + ")"
	// }
	// return val
}
