package character

import (
	"fmt"
	"strings"
)

func (c *Character) Condition(fields []string) (any, error) {
	return false, fmt.Errorf("invalid character condition: .%v.%v", c.Base.Key.String(), strings.Join(fields, "."))
}
