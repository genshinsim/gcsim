package queue

import "strings"

type ExprTreeNode struct {
	Left   *ExprTreeNode
	Right  *ExprTreeNode
	IsLeaf bool
	Op     string //&& || ( )
	Expr   Condition
}

func (e *ExprTreeNode) Clone() *ExprTreeNode {
	//recursively clone left and right
	var next ExprTreeNode
	next.IsLeaf = e.IsLeaf

	//if is leaf then no more conditions or operation
	if next.IsLeaf {
		next.Expr = e.Expr.Clone()
		return &next
	}

	//operation is only on nodes
	next.Op = e.Op

	//other wise grab left and right
	//both shouldn't be nil?
	if e.Left != nil {
		next.Left = e.Left.Clone()
	}
	if e.Right != nil {
		next.Right = e.Right.Clone()
	}

	return &next
}

type Condition struct {
	Fields []string
	Op     string
	Value  int
}

func (c *Condition) String() {
	var sb strings.Builder
	for _, v := range c.Fields {
		sb.WriteString(v)
	}
	sb.WriteString(c.Op)
}

func (c *Condition) Clone() Condition {
	next := *c

	if c.Fields != nil {
		next.Fields = make([]string, len(c.Fields))
		copy(next.Fields, c.Fields)
	}

	return next
}
