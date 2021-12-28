package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (p *Parser) parseCondition() (core.Condition, error) {
	var c core.Condition
	var n item
LOOP:
	for {
		//look for a field
		n = p.next()
		if n.typ != itemField {
			return c, fmt.Errorf("<if - field> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		c.Fields = append(c.Fields, n.val)
		//see if any more fields
		n = p.peek()
		if n.typ != itemField {
			break LOOP
		}
	}

	//scan for comparison op
	n = p.next()
	if n.typ <= itemCompareOp || n.typ >= itemKeyword {
		return c, fmt.Errorf("<if - comp> bad token at line %v - %v: %v", n.line, n.pos, n)
	}
	c.Op = n.val
	//scan for value
	n, err := p.consume(itemNumber)
	if err != nil {
		return c, err
	}
	c.Value, err = itemNumberToInt(n)

	return c, err
}
