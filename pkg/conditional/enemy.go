package conditional

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func evalTarget(c *core.Core, trg *enemy.Enemy, fields []string) (any, error) {
	typ := fields[1]
	switch typ {
	case "mods", "status":
		if err := fieldsCheck(fields, 3, "target "+typ); err != nil {
			return 0, err
		}
		return trg.StatusDuration(fields[2]), nil
	case "element":
		if err := fieldsCheck(fields, 3, "target "+typ); err != nil {
			return 0, err
		}
		ele := fields[2]
		elekey := attributes.StringToEle(ele)
		if elekey == attributes.UnknownElement {
			return 0, fmt.Errorf("bad element condition: invalid element %v", ele)
		}
		result := combat.Durability(0)
		for i := reactable.ModifierInvalid; i < reactable.EndReactableModifier; i++ {
			if i.Element() == elekey && trg.Durability[i] > reactable.ZeroDur && trg.Durability[i] > result {
				result = trg.Durability[i]
			}
		}
		return float64(result), nil
	}
	return nil, fmt.Errorf("bad target condition: invalid type %v", typ)
}

func evalDebuff(c *core.Core, fields []string) (any, error) {
	//.debuff.res.t1.name
	if err := fieldsCheck(fields, 4, "debuff"); err != nil {
		return false, err
	}
	typ := fields[1]
	trg := fields[2]
	mod := fields[3]

	e, err := parseTarget(c, trg)
	if err != nil {
		return false, fmt.Errorf("bad debuff condition: %v", err)
	}

	switch typ {
	case "def", "res":
		log.Printf("WARN: .debuff.%v.t0.%v is deprecated, use .target0.mods.%v", typ, mod, mod)
		return e.StatusIsActive(mod), nil
	default:
		return false, fmt.Errorf("bad debuff condition: invalid type %v", typ)
	}
}

func evalElement(c *core.Core, fields []string) (any, error) {
	//.element.t1.pyro
	if err := fieldsCheck(fields, 3, "element"); err != nil {
		return 0, err
	}
	trg := fields[1]
	ele := fields[2]

	e, err := parseTarget(c, trg)
	if err != nil {
		return 0, fmt.Errorf("bad element condition: %v", err)
	}

	log.Printf("WARN: .element.t0.%v is deprecated, use .target0.element.%v", ele, ele)
	return evalTarget(c, e, []string{"", "element", ele})
}

func parseTarget(c *core.Core, trg string) (*enemy.Enemy, error) {
	pat := regexp.MustCompile(`t(?:arget)?(\d+)`)
	matches := pat.FindStringSubmatch(trg)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid target %v", trg)
	}
	tid, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid target %v", trg)
	}

	t := c.Combat.Enemy(int(tid))
	if t == nil {
		return nil, fmt.Errorf("invalid target %v", tid)
	}

	e, ok := t.(*enemy.Enemy)
	if !ok {
		return nil, fmt.Errorf("target %v is not an enemy", tid)
	}
	return e, nil
}
