package chasca

const (
	a4Key    = "faruzan-a4"
	a4ICDKey = "faruzan-a4-icd"
)

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
}
