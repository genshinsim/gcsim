package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func evalGadgets(c *core.Core, fields []string) (int, error) {
	if err := fieldsCheck(fields, 3, "gadgets"); err != nil {
		return 0, err
	}
	switch fields[1] {
	case "dendrocore":
		return evalDendroCore(c, fields[2])
	case "sourcewaterdroplet":
		return evalSourcewaterDroplet(c, fields[2])
	case "crystallizeshard":
		return evalCrystallizeShard(c, fields[2])
	default:
		return 0, fmt.Errorf("bad gadgets condition: invalid criteria %v", fields[1])
	}
}

func evalDendroCore(c *core.Core, key string) (int, error) {
	switch key {
	case countField:
		count := 0
		for i := 0; i < c.Combat.GadgetCount(); i++ {
			if _, ok := c.Combat.Gadget(i).(*reactable.DendroCore); ok {
				count++
			}
		}
		return count, nil
	default:
		return 0, fmt.Errorf("bad gadgets (dendrocore) condition: invalid criteria %v", key)
	}
}

func evalSourcewaterDroplet(c *core.Core, key string) (int, error) {
	switch key {
	case countField:
		count := 0
		for i := 0; i < c.Combat.GadgetCount(); i++ {
			if _, ok := c.Combat.Gadget(i).(*common.SourcewaterDroplet); ok {
				count++
			}
		}
		return count, nil
	default:
		return 0, fmt.Errorf("bad gadgets (sourcewaterdroplet) condition: invalid criteria %v", key)
	}
}

func evalCrystallizeShard(c *core.Core, key string) (int, error) {
	switch key {
	case "all":
		count := 0
		for i := 0; i < c.Combat.GadgetCount(); i++ {
			cs, ok := c.Combat.Gadget(i).(*reactable.CrystallizeShard)
			if !ok {
				continue
			}
			if c.F < cs.EarliestPickup {
				continue
			}
			count++
		}
		return count, nil
	case attributes.Pyro.String(), attributes.Hydro.String(), attributes.Electro.String(), attributes.Cryo.String():
		return countElementalCrystallizeShards(c, key), nil
	default:
		return 0, fmt.Errorf("bad gadgets (crystallizeshard) condition: invalid criteria %v", key)
	}
}

func countElementalCrystallizeShards(c *core.Core, ele string) int {
	count := 0
	for i := 0; i < c.Combat.GadgetCount(); i++ {
		cs, ok := c.Combat.Gadget(i).(*reactable.CrystallizeShard)
		if !ok {
			continue
		}
		if c.F < cs.EarliestPickup {
			continue
		}
		if cs.Shield.Ele.String() == ele {
			count++
		}
	}
	return count
}
