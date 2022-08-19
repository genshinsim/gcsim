package conditional

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func evalElement[V core.Number](c *core.Core, fields []string) (V, error) {
	//.element.t1.pyro
	if err := fieldsCheck(fields, 3, "element"); err != nil {
		return 0, err
	}
	trg := strings.TrimPrefix(fields[1], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("bad element condition: invalid target %v", trg)
	}
	ele := strings.TrimPrefix(fields[2], ".")
	elekey := attributes.StringToEle(ele)
	if elekey == attributes.UnknownElement {
		return 0, fmt.Errorf("bad element condition: invalid element %v", ele)
	}

	t := c.Combat.Target(int(tid))
	if t == nil {
		return 0, fmt.Errorf("bad element condition: invalid target %v", tid)
	}

	enemy, ok := t.(*enemy.Enemy)
	if !ok {
		return 0, fmt.Errorf("bad element condition: target %v is not an enemy", tid)
	}

	if enemy.AuraContains(elekey) {
		return 1, nil
	}

	return 0, nil
}
