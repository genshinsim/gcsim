package conditional

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func Eval[V core.Number](c *core.Core, fields []string) (V, error) {
	switch fields[0] {
	case ".debuff":
		return evalDebuff[V](c, fields)
	case ".element":
		return evalElement[V](c, fields)
	case ".status":
		if err := fieldsCheck(fields, 2, "status"); err != nil {
			return 0, err
		}
		return V(c.Status.Duration(strings.TrimPrefix(fields[1], "."))), nil
	case ".stam":
		return V(c.Player.Stam), nil
	case ".construct":
		return evalConstruct[V](c, fields)
	case ".keys":
		return evalKeys[V](c, fields)
	default:
		//check if it's a char name; if so check char custom eval func
		name := strings.TrimPrefix(fields[0], ".")
		if key, ok := shortcut.CharNameToKey[name]; ok {
			return evalCharacter[V](c, key, fields)
		}
		return 0, fmt.Errorf("invalid character %v in character condition", name)
	}
}
