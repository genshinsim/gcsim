package bennett

// Decreases Passion Overload's CD by 20%.
func (c *char) a1(cd int) int {
	if c.Base.Ascension < 1 {
		return cd
	}
	return int(float64(cd) * 0.8)
}

// Within the area created by Fantastic Voyage, Passion Overload takes on the following effects:
//
// - CD is reduced by 50%.
func (c *char) a4CD(cd int) int {
	if c.Base.Ascension < 4 || !c.StatModIsActive(burstFieldKey) {
		return cd
	}
	return int(float64(cd) * 0.5)
}

// Within the area created by Fantastic Voyage, Passion Overload takes on the following effects:
//
// - Bennett will not be launched by the effects of Charge Level 2.
func (c *char) a4NoLaunch() bool {
	return c.Base.Ascension >= 4 && c.StatModIsActive(burstFieldKey)
}
