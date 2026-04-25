package conditional

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func evalDebuff(c *core.Core, fields []string) (bool, error) {
	// .debuff.res.t1.name
	if err := fieldsCheck(fields, 4, "debuff"); err != nil {
		return false, err
	}
	typ := fields[1]
	trg := fields[2]
	mod := fields[3]

	e, err := parseTarget(c, trg)
	if err != nil {
		return false, fmt.Errorf("bad debuff condition: %w", err)
	}

	switch typ {
	case "def":
		return e.DefModIsActive(mod), nil
	case "res":
		return e.ResistModIsActive(mod), nil
	default:
		return false, fmt.Errorf("bad debuff condition: invalid type %s", typ)
	}
}

func evalElement(c *core.Core, fields []string) (float64, error) {
	// .element.t1.pyro
	if err := fieldsCheck(fields, 3, "element"); err != nil {
		return 0, err
	}
	trg := fields[1]
	ele := fields[2]

	e, err := parseTarget(c, trg)
	if err != nil {
		return 0, fmt.Errorf("bad element condition: %w", err)
	}

	var eleKey info.ReactionModKey
	switch ele {
	case "burning":
		result := info.Durability(0)
		if e.IsBurning() {
			// TODO: this really should be a min of Burning and BurningFuel, but leaving it as is to avoid
			// breaking existing sims that may be relying on this behavior in the current set of changes
			// this should be fixed
			result = e.GetAuraDurability(info.ReactionModKeyBurning)
		}
		return float64(result), nil
	case "electro":
		eleKey = info.ReactionModKeyElectro
	case "pyro":
		eleKey = info.ReactionModKeyPyro
	case "cryo":
		eleKey = info.ReactionModKeyCryo
	case "hydro":
		eleKey = info.ReactionModKeyHydro
	case "dendro":
		eleKey = info.ReactionModKeyDendro
	case "quicken":
		eleKey = info.ReactionModKeyQuicken
	case "frozen":
		eleKey = info.ReactionModKeyFrozen
	case "geo":
		eleKey = info.ReactionModKeyGeo
	default:
		return 0, fmt.Errorf("bad element condition: invalid element %s", ele)
	}

	dur := e.GetAuraDurability(eleKey)
	if dur < info.ZeroDur {
		return 0, nil
	}
	return float64(dur), nil
}

func parseTarget(c *core.Core, trg string) (*enemy.Enemy, error) {
	trg = strings.TrimPrefix(trg, "t")
	tid, err := strconv.ParseInt(trg, 10, 64)
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
