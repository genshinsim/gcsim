package parse

import (
	"strings"
)

type Node interface {
	String() string
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// writeTo writes the String output to the builder.
	writeTo(*strings.Builder)
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

// Stmt.

type Stmt interface {
	Node
	stmtNode()
}

type (

	// BlockStmt represents a brace statement list
	BlockStmt struct {
		List []Node
		Pos
	}

	// CommentStmt holds a comment.
	CommentStmt struct {
		Pos
		Text string // Comment text.
	}

	// ActionStmt represents a sim action; Does not produce a value
	ActionStmt struct {
		Pos
	}

	// IfStmt represents an if block
	IfStmt struct {
		Pos
		Condition Expr       //TODO: this should be an expr?
		IfBlock   *BlockStmt // What to execute if true
		ElseBlock *BlockStmt // What to execute if false
	}

	// FnStmt represents a fn block
	FnStmt struct {
		Pos
		Block *BlockStmt
	}
)

// stmtNode()
func (*BlockStmt) stmtNode()
func (*CommentStmt) stmtNode()

// BlockStmt.
func newBlockStmt(pos Pos) *BlockStmt {
	return &BlockStmt{Pos: pos}
}
func (l *BlockStmt) append(n Node) {
	l.List = append(l.List, n)
}

func (l *BlockStmt) String() string {
	var sb strings.Builder
	l.writeTo(&sb)
	return sb.String()
}

func (l *BlockStmt) writeTo(sb *strings.Builder) {
	for _, n := range l.List {
		n.writeTo(sb)
	}
}

func (l *BlockStmt) CopyBlock() *BlockStmt {
	if l == nil {
		return l
	}
	n := newBlockStmt(l.Pos)
	for _, elem := range l.List {
		n.append(elem.Copy())
	}
	return n
}

func (l *BlockStmt) Copy() Node {
	return l.CopyBlock()
}

func newComment(pos Pos, text string) *CommentStmt {
	return &CommentStmt{Pos: pos, Text: text}
}

// CommentStmt.

func (c *CommentStmt) String() string {
	var sb strings.Builder
	c.writeTo(&sb)
	return sb.String()
}

func (c *CommentStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("//")
	sb.WriteString(c.Text)
}

func (c *CommentStmt) Copy() Node {
	return &CommentStmt{Pos: c.Pos, Text: c.Text}
}

// IfStmt.

func newIfStmt(pos Pos) *IfStmt {
	return &IfStmt{Pos: pos}
}

func (i *IfStmt) SetCondition(e Expr) {
	i.Condition = e
}

func (i *IfStmt) SetIfBlock(b *BlockStmt) {
	i.IfBlock = b
}

func (i *IfStmt) SetElseBlock(b *BlockStmt) {
	i.ElseBlock = b
}

func (i *IfStmt) String() string {
	var sb strings.Builder
	i.writeTo(&sb)
	return sb.String()
}

func (i *IfStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("if ")
	i.Condition.writeTo(sb)
	sb.WriteString(" {")
	i.IfBlock.writeTo(sb)
	sb.WriteString(" }")
	if i.ElseBlock != nil {
		sb.WriteString("else ")
		sb.WriteString(" {")
		sb.WriteString(i.ElseBlock.String())
		sb.WriteString(" }")
	}
}

func (i *IfStmt) Copy() Node {
	return &IfStmt{
		Pos:       i.Pos,
		Condition: i.Condition.CopyExpr(),
		IfBlock:   i.IfBlock.CopyBlock(),
		ElseBlock: i.ElseBlock.CopyBlock(),
	}
}

// Expr.

type Expr interface {
	Node
	exprNode()
	CopyExpr() Expr
}

// An expression is represented by a tree consisting of one or
// more of the following concrete expression nodes
type (
	BasicLit struct {
		Pos
		Kind  tokenType
		Value string // literal string; eg. 42, 3.14
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Pos
		FunVal string // function name
		Fun    Expr   // function expression
		Args   []Expr // function arguments; or nil
	}

	// A UnaryExpr node represents a unary expression.
	UnaryExpr struct {
		Pos
		Op Token
		X  Expr // operand
	}

	//A BinaryExpr node represents a binary expression i.e. a > b, 1 + 1, etc..
	BinaryExpr struct {
		Pos
		LHS Expr
		RHS Expr  // need to evalute to same type as lhs
		Op  Token //should be > itemCompareOP and < itemDot
	}
)

//exprNode()
func (*BinaryExpr) exprNode()

// BinaryExpr.

func newBinaryExpr(pos Pos) *BinaryExpr {
	return &BinaryExpr{Pos: pos}
}

func (b *BinaryExpr) SetLHS(e Expr) {
	b.LHS = e
}

func (b *BinaryExpr) SetRHS(e Expr) {
	b.RHS = e
}

func (b *BinaryExpr) SetOP(op Token) {
	b.Op = op
}

func (b *BinaryExpr) CopyBinaryExpr() *BinaryExpr {
	if b == nil {
		return b
	}
	n := newBinaryExpr(b.Pos)
	n.LHS = b.LHS.CopyExpr()
	n.RHS = b.RHS.CopyExpr()
	n.Op = b.Op
	return n
}

func (b *BinaryExpr) CopyExpr() Expr {
	return b.CopyBinaryExpr()
}

func (b *BinaryExpr) Copy() Node {
	return b.CopyBinaryExpr()
}

func (b *BinaryExpr) String() string {
	var sb strings.Builder
	b.writeTo(&sb)
	return sb.String()
}

func (b *BinaryExpr) writeTo(sb *strings.Builder) {
	b.LHS.writeTo(sb)
	sb.WriteString(" ")
	b.RHS.writeTo(sb)
	sb.WriteString(" ")
	sb.WriteString(b.Op.String())
}
