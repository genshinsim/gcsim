package conditional

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func evalDebuff[V core.Number](c *core.Core, fields []string) (V, error) {
	//.debuff.res.t1.name
	if err := fieldsCheck(fields, 4, "debuff"); err != nil {
		return 0, err
	}
	typ := strings.TrimPrefix(fields[1], ".")
	trg := strings.TrimPrefix(fields[2], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("bad debuff condition: invalid target %v", trg)
	}

	d := strings.TrimPrefix(fields[3], ".")
	t := c.Combat.Target(int(tid))
	if t == nil {
		return 0, fmt.Errorf("bad debuff condition: invalid target %v", tid)
	}

	enemy, ok := t.(*enemy.Enemy)
	if !ok {
		return 0, fmt.Errorf("bad debuff condition: target %v is not an enemy", tid)
	}

	switch typ {
	case "res":
		if enemy.ResistModIsActive(d) {
			return 1, nil
		}
	case "def":
		if enemy.DefModIsActive(d) {
			return 1, nil
		}
	}
	return 0, nil
}
