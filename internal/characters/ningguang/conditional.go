package ningguang

func (c *char) Condition(field string) int64 {
	switch field {
	case "jadeCount":
		return int64(c.jadeCount)
	case "prevAttack":
		return int64(c.prevAttack)
	default:
		return 0
	}
}
