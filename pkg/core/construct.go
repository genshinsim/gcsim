package core

type GeoConstructType int

const (
	GeoConstructNingSkill GeoConstructType = iota
	GeoConstructZhongliSkill
	GeoConstructTravellerSkill
	GeoConstructTravellerBurst
	GeoConstructAlbedoSkill
	EndGeoConstructType
)

var ConstructString = [...]string{
	"NingSkill",
	"ZhongliSkill",
	"TravellerSkill",
	"TravellerBurst",
	"AlbedoSkill",
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
}
