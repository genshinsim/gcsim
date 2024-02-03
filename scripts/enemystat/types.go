package main

//nolint:tagliatelle // need to match datamine
type HpDrop struct {
	DropId    int     `json:"dropId"`
	HpPercent float64 `json:"hpPercent"`
}

//nolint:tagliatelle // need to match datamine
type propGrowCurve struct {
	GrowCurve string `json:"growCurve"`
	// ...
}

//nolint:tagliatelle // need to match datamine
type monsterExcelConfig struct {
	MonsterName     string          `json:"monsterName"`
	Typ             string          `json:"type"`
	HpDrops         []HpDrop        `json:"hpDrops"`
	DescribeId      int             `json:"describeId"`
	KillDropId      int             `json:"killDropId"`
	HpBase          float64         `json:"hpBase"`
	PropGrowCurves  []propGrowCurve `json:"propGrowCurves"`
	FireSubHurt     float64         `json:"fireSubHurt"`
	GrassSubHurt    float64         `json:"grassSubHurt"`
	WaterSubHurt    float64         `json:"waterSubHurt"`
	ElecSubHurt     float64         `json:"elecSubHurt"`
	WindSubHurt     float64         `json:"windSubHurt"`
	IceSubHurt      float64         `json:"iceSubHurt"`
	RockSubHurt     float64         `json:"rockSubHurt"`
	PhysicalSubHurt float64         `json:"physicalSubHurt"`
	Id              int             `json:"id"`
	// ...
}

//nolint:tagliatelle // need to match datamine
type monsterDescribeExcelConfig struct {
	Id              int    `json:"id"`
	NameTextMapHash int    `json:"nameTextMapHash"`
	Icon            string `json:"icon"`
	// ...
}
