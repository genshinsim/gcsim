package parse

import (
	"strconv"
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

	// ActionStmt represents a sim action; Does not produce a value
	ActionStmt struct {
		Pos
	}

	// AssignStmt represents assigning of a value to a previously declared variable
	AssignStmt struct {
		Pos
	}

	// LetStmt represents a variable assignment
	LetStmt struct {
		Pos
		Ident Token
		Val   Expr
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
		Ident Token
		Block *BlockStmt
	}
)

// stmtNode()
func (*BlockStmt) stmtNode() {}
func (*LetStmt) stmtNode()   {}
func (*IfStmt) stmtNode()    {}
func (*FnStmt) stmtNode()    {}

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
		sb.WriteString(";\n")
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

// LetStmt.

func (l *LetStmt) String() string {
	var sb strings.Builder
	l.writeTo(&sb)
	return sb.String()
}

func (l *LetStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("let ")
	sb.WriteString(l.Ident.String())
	sb.WriteString(" = ")
	l.Val.writeTo(sb)
}

func (l *LetStmt) CopyBlock() *LetStmt {
	if l == nil {
		return l
	}
	n := &LetStmt{
		Pos:   l.Pos,
		Ident: l.Ident,
	}
	n.Val = l.Val.CopyExpr()
	return n
}

func (l *LetStmt) Copy() Node {
	return l.CopyBlock()
}

// FnStmt.

func (f *FnStmt) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *FnStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("fn ")
	sb.WriteString(f.Ident.String())
	sb.WriteString(" {\n")
	f.Block.writeTo(sb)
	sb.WriteString(" }")
}

func (f *FnStmt) CopyFn() *FnStmt {
	if f == nil {
		return f
	}
	return &FnStmt{
		Pos:   f.Pos,
		Ident: f.Ident,
		Block: f.Block.CopyBlock(),
	}
}

func (f *FnStmt) Copy() Node {
	return f.CopyFn()
}

// IfStmt.

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
	sb.WriteString(" {\n")
	i.IfBlock.writeTo(sb)
	sb.WriteString("}")
	if i.ElseBlock != nil {
		sb.WriteString("else {\n")
		sb.WriteString(i.ElseBlock.String())
		sb.WriteString("}")
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
	NumberLit struct {
		Pos
		IntVal   int64
		FloatVal float64
		IsInt    bool
	}

	StringLit struct {
		Pos
		Value float64
	}
	BoolLit struct {
		Pos
		Value float64
	}

	Ident struct {
		Pos
		Value string
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
		Op    Token
		Right Expr // operand
	}

	//A BinaryExpr node represents a binary expression i.e. a > b, 1 + 1, etc..
	BinaryExpr struct {
		Pos
		Left  Expr
		Right Expr  // need to evalute to same type as lhs
		Op    Token //should be > itemCompareOP and < itemDot
	}
)

//exprNode()
func (*NumberLit) exprNode()  {}
func (*Ident) exprNode()      {}
func (*CallExpr) exprNode()   {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}

// BasicLit.

func (b *NumberLit) CopyExpr() Expr {
	return &NumberLit{Pos: b.Pos, IntVal: b.IntVal}
}

func (b *NumberLit) Copy() Node {
	return b.CopyExpr()
}

func (b *NumberLit) String() string {
	var sb strings.Builder
	b.writeTo(&sb)
	return sb.String()
}

func (b *NumberLit) writeTo(sb *strings.Builder) {
	if b.IsInt {
		sb.WriteString(strconv.FormatInt(b.IntVal, 10))
	} else {
		sb.WriteString(strconv.FormatFloat(b.FloatVal, 'f', -1, 64))
	}
}

// Ident.

func (i *Ident) CopyExpr() Expr {
	return &Ident{Pos: i.Pos, Value: i.Value}
}

func (i *Ident) Copy() Node {
	return i.CopyExpr()
}

func (b *Ident) String() string {
	var sb strings.Builder
	b.writeTo(&sb)
	return sb.String()
}

func (b *Ident) writeTo(sb *strings.Builder) {
	sb.WriteString(b.Value)
}

// UnaryExpr.

func (u *UnaryExpr) CopyUnaryExpr() *UnaryExpr {
	if u == nil {
		return u
	}
	n := &UnaryExpr{Pos: u.Pos}
	n.Right = u.Right.CopyExpr()
	n.Op = u.Op
	return n
}

func (u *UnaryExpr) CopyExpr() Expr {
	return u.CopyUnaryExpr()
}

func (u *UnaryExpr) Copy() Node {
	return u.CopyUnaryExpr()
}

func (u *UnaryExpr) String() string {
	var sb strings.Builder
	u.writeTo(&sb)
	return sb.String()
}

func (u *UnaryExpr) writeTo(sb *strings.Builder) {
	sb.WriteString("(")
	sb.WriteString(u.Op.String())
	u.Right.writeTo(sb)
	sb.WriteString(")")
}

// BinaryExpr.

func (b *BinaryExpr) CopyBinaryExpr() *BinaryExpr {
	if b == nil {
		return b
	}
	n := &BinaryExpr{Pos: b.Pos}
	n.Left = b.Left.CopyExpr()
	n.Right = b.Right.CopyExpr()
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
	sb.WriteString("(")
	b.Left.writeTo(sb)
	sb.WriteString(" ")
	sb.WriteString(b.Op.String())
	sb.WriteString(" ")
	b.Right.writeTo(sb)
	sb.WriteString(")")
}
