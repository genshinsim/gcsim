package construct

import "github.com/genshinsim/gcsim/pkg/core/info"

type GeoConstructType int

const (
	GeoConstructInvalid GeoConstructType = iota
	GeoConstructNingSkill
	GeoConstructZhongliSkill
	GeoConstructTravellerSkill
	GeoConstructTravellerBurst
	GeoConstructAlbedoSkill
	GeoConstructIttoSkill
	GeoConstructLunarCrystallize
	EndGeoConstructType
)

var ConstructString = [...]string{
	"Invalid",
	"NingSkill",
	"ZhongliSkill",
	"TravellerSkill",
	"TravellerBurst",
	"AlbedoSkill",
	"IttoSkill",
	"LunarCrystallize",
}

var ConstructNameToKey = map[string]GeoConstructType{
	"ningguang":         GeoConstructNingSkill,
	"zhongli":           GeoConstructZhongliSkill,
	"traveler_skill":    GeoConstructTravellerSkill,
	"traveler_burst":    GeoConstructTravellerBurst,
	"albedo":            GeoConstructAlbedoSkill,
	"itto":              GeoConstructIttoSkill,
	"lunar_crystallize": GeoConstructLunarCrystallize,
}

func (c GeoConstructType) String() string {
	return ConstructString[c]
}

type Construct interface {
	OnDestruct()
	Key() int
	Type() GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
	Direction() info.Point
	Pos() info.Point
}
