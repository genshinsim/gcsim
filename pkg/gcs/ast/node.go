package ast

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
	CopyStmt() Stmt
	stmtNode()
}

type (

	// BlockStmt represents a brace statement list
	BlockStmt struct {
		List []Node
		Pos
	}

	// AssignStmt represents assigning of a value to a previously declared variable
	AssignStmt struct {
		Pos
		Ident Token
		Val   Expr
	}

	// LetStmt represents a variable assignment. Number only
	LetStmt struct {
		Pos
		Ident Token
		Type  ExprType
		Val   Expr
	}

	// ReturnStmt represents return <expr>.
	ReturnStmt struct {
		Pos
		Val Expr
	}

	// CtrlStmt represents continue, break, and fallthrough
	CtrlStmt struct {
		Pos
		Typ CtrlTyp
	}

	// IfStmt represents an if block
	IfStmt struct {
		Pos
		Condition Expr       //TODO: this should be an expr?
		IfBlock   *BlockStmt // What to execute if true
		ElseBlock Stmt       // What to execute if false
	}

	// SwitchStmt represent a switch block
	SwitchStmt struct {
		Pos
		Condition Expr // the condition to switch on
		Cases     []*CaseStmt
		Default   *BlockStmt // default case
	}

	// CaseStmt represents a case in a switch block
	CaseStmt struct {
		Pos
		Condition Expr
		Body      *BlockStmt
	}

	// A FnStmt node represents a function declared with syntax fn ident(..args) { block }.
	// Functionally the same as a LetStmt
	FnStmt struct {
		Pos
		Ident Token
		Func  *FuncLit
	}

	// WhileStmt represents a while block
	WhileStmt struct {
		Pos
		Condition  Expr       //TODO: this should be an expr?
		WhileBlock *BlockStmt // What to execute if true
	}

	// ForStmt represents a for block
	ForStmt struct {
		Pos
		Init Stmt // initialization statement; or nil
		Cond Expr // condition; or nil
		Post Stmt // post iteration statement; or nil
		Body *BlockStmt
	}
)

type CtrlTyp int

const (
	InvalidCtrl CtrlTyp = iota
	CtrlBreak
	CtrlContinue
	CtrlFallthrough
)

// stmtNode()
func (*BlockStmt) stmtNode()  {}
func (*AssignStmt) stmtNode() {}
func (*LetStmt) stmtNode()    {}
func (*CtrlStmt) stmtNode()   {}
func (*ReturnStmt) stmtNode() {}
func (*IfStmt) stmtNode()     {}
func (*SwitchStmt) stmtNode() {}
func (*CaseStmt) stmtNode()   {}
func (*FnStmt) stmtNode()     {}
func (*WhileStmt) stmtNode()  {}
func (*ForStmt) stmtNode()    {}

// BlockStmt.
func newBlockStmt(pos Pos) *BlockStmt {
	return &BlockStmt{Pos: pos}
}
func (b *BlockStmt) append(n Node) {
	b.List = append(b.List, n)
}

func (b *BlockStmt) String() string {
	var sb strings.Builder
	b.writeTo(&sb)
	return sb.String()
}

func (b *BlockStmt) writeTo(sb *strings.Builder) {
	for _, n := range b.List {
		n.writeTo(sb)
		sb.WriteString("\n")
	}
}

func (b *BlockStmt) CopyBlock() *BlockStmt {
	if b == nil {
		return b
	}
	n := newBlockStmt(b.Pos)
	for _, elem := range b.List {
		n.append(elem.Copy())
	}
	return n
}

func (b *BlockStmt) CopyStmt() Stmt {
	return b.CopyBlock()
}

func (b *BlockStmt) Copy() Node {
	return b.CopyBlock()
}

// AssignStmt.

func (a *AssignStmt) String() string {
	var sb strings.Builder
	a.writeTo(&sb)
	return sb.String()
}

func (a *AssignStmt) writeTo(sb *strings.Builder) {
	sb.WriteString(a.Ident.String())
	sb.WriteString(" = ")
	a.Val.writeTo(sb)
}

func (a *AssignStmt) CopyAssign() *AssignStmt {
	if a == nil {
		return a
	}
	n := &AssignStmt{
		Pos:   a.Pos,
		Ident: a.Ident,
	}
	n.Val = a.Val.CopyExpr()
	return n
}

func (a *AssignStmt) CopyStmt() Stmt {
	return a.CopyAssign()
}

func (a *AssignStmt) Copy() Node {
	return a.CopyAssign()
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
	if l.Type != nil {
		sb.WriteString(" ")
		l.Type.writeTo(sb)
	}
	sb.WriteString(" = ")
	if l.Val != nil {
		l.Val.writeTo(sb)
	}
}

func (l *LetStmt) CopyLet() *LetStmt {
	if l == nil {
		return l
	}
	n := &LetStmt{
		Pos:   l.Pos,
		Ident: l.Ident,
		Type:  l.Type.CopyExprType(),
	}
	n.Val = l.Val.CopyExpr()
	return n
}

func (l *LetStmt) CopyStmt() Stmt {
	return l.CopyLet()
}

func (l *LetStmt) Copy() Node {
	return l.CopyLet()
}

// ReturnStmt.

func (r *ReturnStmt) String() string {
	var sb strings.Builder
	r.writeTo(&sb)
	return sb.String()
}

func (r *ReturnStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("return ")
	r.Val.writeTo(sb)
}

func (r *ReturnStmt) CopyReturn() *ReturnStmt {
	if r == nil {
		return r
	}
	n := &ReturnStmt{
		Pos: r.Pos,
	}
	n.Val = r.Val.CopyExpr()
	return n
}

func (r *ReturnStmt) CopyStmt() Stmt {
	return r.CopyReturn()
}

func (r *ReturnStmt) Copy() Node {
	return r.CopyReturn()
}

// CtrlStmt.

func (c *CtrlStmt) String() string {
	var sb strings.Builder
	c.writeTo(&sb)
	return sb.String()
}

func (c *CtrlStmt) writeTo(sb *strings.Builder) {
	switch c.Typ {
	case CtrlContinue:
		sb.WriteString("continue")
	case CtrlBreak:
		sb.WriteString("break")
	case CtrlFallthrough:
		sb.WriteString("fallthrough")
	}
}

func (c *CtrlStmt) CopyControl() *CtrlStmt {
	if c == nil {
		return c
	}
	n := &CtrlStmt{
		Pos: c.Pos,
		Typ: c.Typ,
	}
	return n
}

func (c *CtrlStmt) CopyStmt() Stmt {
	return c.CopyControl()
}

func (c *CtrlStmt) Copy() Node {
	return c.CopyControl()
}

// IfStmt.

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

func (i *IfStmt) CopyIfStmt() *IfStmt {
	if i == nil {
		return nil
	}
	n := &IfStmt{
		Pos:       i.Pos,
		Condition: i.Condition.CopyExpr(),
		IfBlock:   i.IfBlock.CopyBlock(),
	}
	if i.ElseBlock != nil {
		n.ElseBlock = i.ElseBlock.CopyStmt()
	}
	return n
}

func (i *IfStmt) CopyStmt() Stmt {
	return i.CopyIfStmt()
}

func (i *IfStmt) Copy() Node {
	return i.CopyIfStmt()
}

// SwitchStmt.

func (s *SwitchStmt) String() string {
	var sb strings.Builder
	s.writeTo(&sb)
	return sb.String()
}

func (s *SwitchStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("switch ")
	s.Condition.writeTo(sb)
	sb.WriteString(" {\n")
	for _, v := range s.Cases {
		v.writeTo(sb)
	}
	if s.Default != nil {
		sb.WriteString("default: {\n")
		s.Default.writeTo(sb)
		sb.WriteString("}")
	}
	sb.WriteString("}")
}

func (s *SwitchStmt) CopySwitch() *SwitchStmt {
	if s == nil {
		return nil
	}
	n := &SwitchStmt{
		Pos:     s.Pos,
		Cases:   make([]*CaseStmt, 0, len(s.Cases)),
		Default: s.Default.CopyBlock(),
	}
	if s.Condition != nil {
		n.Condition = s.Condition.CopyExpr()
	}
	for i := range s.Cases {
		n.Cases = append(n.Cases, s.Cases[i].CopyCase())
	}
	return n
}

func (s *SwitchStmt) CopyStmt() Stmt {
	return s.CopySwitch()
}

func (s *SwitchStmt) Copy() Node {
	return s.CopySwitch()
}

// CaseStmt.

func (c *CaseStmt) String() string {
	var sb strings.Builder
	c.writeTo(&sb)
	return sb.String()
}

func (c *CaseStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("case ")
	c.Condition.writeTo(sb)
	sb.WriteString(" {\n")
	c.Body.writeTo(sb)
	sb.WriteString("}")
}

func (c *CaseStmt) CopyCase() *CaseStmt {
	if c == nil {
		return nil
	}
	return &CaseStmt{
		Pos:       c.Pos,
		Condition: c.Condition.CopyExpr(),
		Body:      c.Body.CopyBlock(),
	}
}

func (c *CaseStmt) CopyStmt() Stmt {
	return c.CopyCase()
}

func (c *CaseStmt) Copy() Node {
	return c.CopyCase()
}

// FnStmt.

func (f *FnStmt) CopyFn() Stmt {
	if f == nil {
		return nil
	}
	n := &FnStmt{
		Pos:   f.Pos,
		Ident: f.Ident,
		Func:  f.Func.copyFuncLit(),
	}

	return n
}

func (f *FnStmt) CopyStmt() Stmt {
	return f.CopyFn()
}

func (f *FnStmt) Copy() Node {
	return f.CopyStmt()
}

func (f *FnStmt) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *FnStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("fn ")
	sb.WriteString(f.Ident.String())
	f.Func.writeTo(sb)
}

// WhileStmt.

func (w *WhileStmt) String() string {
	var sb strings.Builder
	w.writeTo(&sb)
	return sb.String()
}

func (w *WhileStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("while ")
	w.Condition.writeTo(sb)
	sb.WriteString(" {\n")
	w.WhileBlock.writeTo(sb)
	sb.WriteString("}")
}

func (w *WhileStmt) CopyWhileStmt() *WhileStmt {
	if w == nil {
		return nil
	}
	return &WhileStmt{
		Pos:        w.Pos,
		Condition:  w.Condition.CopyExpr(),
		WhileBlock: w.WhileBlock.CopyBlock(),
	}
}

func (w *WhileStmt) CopyStmt() Stmt {
	return w.CopyWhileStmt()
}

func (w *WhileStmt) Copy() Node {
	return w.CopyWhileStmt()
}

// ForStmt.

func (f *ForStmt) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *ForStmt) writeTo(sb *strings.Builder) {
	sb.WriteString("for ")
	if f.Init != nil {
		f.Init.writeTo(sb)
	}
	sb.WriteString("; ")
	if f.Cond != nil {
		f.Cond.writeTo(sb)
	}
	sb.WriteString("; ")
	if f.Post != nil {
		f.Post.writeTo(sb)
	}
	sb.WriteString(" {\n")
	f.Body.writeTo(sb)
	sb.WriteString("}")
}

func (f *ForStmt) CopyForStmt() *ForStmt {
	if f == nil {
		return nil
	}
	n := &ForStmt{
		Pos:  f.Pos,
		Body: f.Body.CopyBlock(),
	}
	if f.Init != nil {
		n.Init = f.Init.CopyStmt()
	}
	if f.Cond != nil {
		n.Cond = f.Cond.CopyExpr()
	}
	if f.Post != nil {
		n.Post = f.Post.CopyStmt()
	}
	return n
}

func (f *ForStmt) CopyStmt() Stmt {
	return f.CopyForStmt()
}

func (f *ForStmt) Copy() Node {
	return f.CopyForStmt()
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
		IsFloat  bool
	}

	StringLit struct {
		Pos
		Value string
	}

	// FuncExpr is just a wrapper around FuncLit representing an anonymous function declaration
	// This node should only exists following a let statement
	// The FuncExpr itself should have null type. The actual FuncLit will have it's own return type
	FuncExpr struct {
		Pos  // position of the fn keyword
		Func *FuncLit
	}

	// A FuncLit node represents a function literal.
	FuncLit struct {
		Pos       // position of the starting (
		Signature *FuncType
		Args      []*Ident
		Body      *BlockStmt
	}

	Ident struct {
		Pos
		Value string
	}

	Field struct {
		Pos
		Value []string
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Pos
		Fun  Expr   // function expression
		Args []Expr // function arguments; or nil
	}

	// A UnaryExpr node represents a unary expression.
	UnaryExpr struct {
		Pos
		Op    Token
		Right Expr // operand
	}

	// A BinaryExpr node represents a binary expression i.e. a > b, 1 + 1, etc..
	BinaryExpr struct {
		Pos
		Left  Expr
		Right Expr  // need to evalute to same type as lhs
		Op    Token // should be > itemCompareOP and < itemDot
	}

	MapExpr struct {
		Pos
		Fields map[string]Expr
	}
)

// exprNode()
func (*NumberLit) exprNode()  {}
func (*StringLit) exprNode()  {}
func (*FuncExpr) exprNode()   {}
func (*FuncLit) exprNode()    {}
func (*Ident) exprNode()      {}
func (*Field) exprNode()      {}
func (*CallExpr) exprNode()   {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*MapExpr) exprNode()    {}

// NumberLit.

func (n *NumberLit) CopyExpr() Expr {
	if n == nil {
		return nil
	}
	return &NumberLit{
		Pos:      n.Pos,
		IntVal:   n.IntVal,
		FloatVal: n.FloatVal,
		IsFloat:  n.IsFloat,
	}
}

func (n *NumberLit) Copy() Node {
	return n.CopyExpr()
}

func (n *NumberLit) String() string {
	var sb strings.Builder
	n.writeTo(&sb)
	return sb.String()
}

func (n *NumberLit) writeTo(sb *strings.Builder) {
	if n.IsFloat {
		sb.WriteString(strconv.FormatFloat(n.FloatVal, 'f', -1, 64))
	} else {
		sb.WriteString(strconv.FormatInt(n.IntVal, 10))
	}
}

// StringLit.

func (n *StringLit) CopyExpr() Expr {
	if n == nil {
		return nil
	}
	return &StringLit{
		Pos:   n.Pos,
		Value: n.Value,
	}
}

func (n *StringLit) Copy() Node {
	return n.CopyExpr()
}

func (n *StringLit) String() string {
	return n.Value
}

func (n *StringLit) writeTo(sb *strings.Builder) {
	sb.WriteString(n.Value)
}

// FuncExpr.

func (f *FuncExpr) CopyExpr() Expr {
	if f == nil {
		return nil
	}
	n := &FuncExpr{
		Pos:  f.Pos,
		Func: f.Func.copyFuncLit(),
	}
	return n
}

func (f *FuncExpr) Copy() Node {
	return f.CopyExpr()
}

func (f *FuncExpr) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *FuncExpr) writeTo(sb *strings.Builder) {
	sb.WriteString("fn")
	f.Func.writeTo(sb)
}

// FuncLit.

func (f *FuncLit) CopyExpr() Expr {
	if f == nil {
		return nil
	}
	return f.copyFuncLit()
}

func (f *FuncLit) copyFuncLit() *FuncLit {
	n := &FuncLit{
		Pos:       f.Pos,
		Signature: f.Signature.copyFuncType(),
		Body:      f.Body.CopyBlock(),
	}
	n.Args = make([]*Ident, 0, len(f.Args))
	for _, v := range f.Args {
		n.Args = append(n.Args, v.CopyIdent())
	}
	return n
}

func (f *FuncLit) Copy() Node {
	return f.CopyExpr()
}

func (f *FuncLit) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *FuncLit) writeTo(sb *strings.Builder) {
	sb.WriteString("(")
	for i, v := range f.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
		if f.Signature.ArgsType[i] != nil {
			sb.WriteString(" ")
			f.Signature.ArgsType[i].writeTo(sb)
		}
	}
	sb.WriteString(")")
	if f.Signature.ResultType != nil {
		sb.WriteString(" ")
		f.Signature.ResultType.writeTo(sb)
	}
	sb.WriteString(" {\n")
	f.Body.writeTo(sb)
	sb.WriteString("}")
}

// Ident.

func (i *Ident) CopyIdent() *Ident {
	if i == nil {
		return nil
	}
	return &Ident{Pos: i.Pos, Value: i.Value}
}

func (i *Ident) CopyExpr() Expr {
	return i.CopyIdent()
}

func (i *Ident) Copy() Node {
	return i.CopyIdent()
}

func (i *Ident) String() string {
	var sb strings.Builder
	i.writeTo(&sb)
	return sb.String()
}

func (i *Ident) writeTo(sb *strings.Builder) {
	sb.WriteString(i.Value)
}

// Field.

func (i *Field) CopyField() *Field {
	if i == nil {
		return nil
	}
	dst := make([]string, len(i.Value))
	copy(dst, i.Value)
	return &Field{Pos: i.Pos, Value: dst}
}

func (i *Field) CopyExpr() Expr {
	return i.CopyField()
}

func (i *Field) Copy() Node {
	return i.CopyField()
}

func (i *Field) String() string {
	var sb strings.Builder
	i.writeTo(&sb)
	return sb.String()
}

func (i *Field) writeTo(sb *strings.Builder) {
	for _, v := range i.Value {
		sb.WriteString(v)
	}
}

// CallExpr.

func (f *CallExpr) CopyFn() Expr {
	if f == nil {
		return nil
	}
	n := &CallExpr{
		Pos:  f.Pos,
		Fun:  f.Fun.CopyExpr(),
		Args: make([]Expr, 0, len(f.Args)),
	}
	for i := range f.Args {
		n.Args = append(n.Args, f.Args[i].CopyExpr())
	}

	return n
}

func (f *CallExpr) CopyExpr() Expr {
	return f.CopyFn()
}

func (f *CallExpr) Copy() Node {
	return f.CopyExpr()
}

func (f *CallExpr) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *CallExpr) writeTo(sb *strings.Builder) {
	f.Fun.writeTo(sb)
	sb.WriteString("(")
	for i, v := range f.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		v.writeTo(sb)
	}
	sb.WriteString(")")
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

// MapExpr.

func (m *MapExpr) CopyExpr() Expr {
	if m == nil {
		return m
	}
	n := &MapExpr{
		Pos:    m.Pos,
		Fields: make(map[string]Expr),
	}
	for k, v := range m.Fields {
		n.Fields[k] = v.CopyExpr()
	}
	return n
}

func (m *MapExpr) Copy() Node {
	return m.CopyExpr()
}

func (m *MapExpr) String() string {
	var sb strings.Builder
	m.writeTo(&sb)
	return sb.String()
}

func (m *MapExpr) writeTo(sb *strings.Builder) {
	sb.WriteString("[")
	done := false
	for k, v := range m.Fields {
		if done {
			sb.WriteString(", ")
		}
		done = true

		sb.WriteString(k)
		sb.WriteString(" = ")
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
}

type ExprType interface {
	Node
	exprTypeNode()
	CopyExprType() ExprType
}

type (
	NumberType struct {
		Pos // position of :, or of the ident if defaulting to NumberType
	}
	StringType struct {
		Pos // position of :
	}

	MapType struct {
		Pos // position of keyword map
	}
	FuncType struct {
		Pos        // position of opening (
		ArgsType   []ExprType
		ResultType ExprType
	}
)

// exprTypeNode()
func (*NumberType) exprTypeNode() {}
func (*StringType) exprTypeNode() {}
func (*MapType) exprTypeNode()    {}
func (*FuncType) exprTypeNode()   {}

// NumberType.
func (n *NumberType) CopyExprType() ExprType {
	if n == nil {
		return nil
	}
	return &NumberType{
		Pos: n.Pos,
	}
}

func (n *NumberType) Copy() Node {
	return n.CopyExprType()
}

func (n *NumberType) String() string {
	var sb strings.Builder
	n.writeTo(&sb)
	return sb.String()
}

func (n *NumberType) writeTo(sb *strings.Builder) {
	sb.WriteString("number")
}

// StringType.
func (s *StringType) CopyExprType() ExprType {
	if s == nil {
		return nil
	}
	return &StringType{
		Pos: s.Pos,
	}
}

func (s *StringType) Copy() Node {
	return s.CopyExprType()
}

func (s *StringType) String() string {
	var sb strings.Builder
	s.writeTo(&sb)
	return sb.String()
}

func (s *StringType) writeTo(sb *strings.Builder) {
	sb.WriteString("string")
}

// MapType.
func (m *MapType) CopyExprType() ExprType {
	if m == nil {
		return nil
	}
	return &MapType{
		Pos: m.Pos,
	}
}

func (m *MapType) Copy() Node {
	return m.CopyExprType()
}

func (m *MapType) String() string {
	var sb strings.Builder
	m.writeTo(&sb)
	return sb.String()
}

func (m *MapType) writeTo(sb *strings.Builder) {
	sb.WriteString("map")
}

// FuncType.
func (f *FuncType) copyFuncType() *FuncType {
	next := &FuncType{
		Pos:        f.Pos,
		ArgsType:   make([]ExprType, 0, len(f.ArgsType)),
		ResultType: f.ResultType.CopyExprType(),
	}

	for _, v := range f.ArgsType {
		next.ArgsType = append(next.ArgsType, v.CopyExprType())
	}

	return next
}

func (f *FuncType) CopyExprType() ExprType {
	if f == nil {
		return nil
	}
	return f.copyFuncType()
}

func (f *FuncType) Copy() Node {
	return f.CopyExprType()
}

func (f *FuncType) String() string {
	var sb strings.Builder
	f.writeTo(&sb)
	return sb.String()
}

func (f *FuncType) writeTo(sb *strings.Builder) {
	sb.WriteString("fn(")
	for i, v := range f.ArgsType {
		if i > 0 {
			sb.WriteString(", ")
		}
		if v != nil {
			v.writeTo(sb)
		}
	}
	sb.WriteString(")")
	if f.ResultType != nil {
		sb.WriteString(" ")
		f.ResultType.writeTo(sb)
	}
}
