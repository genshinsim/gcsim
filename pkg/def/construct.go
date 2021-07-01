package def

type GeoConstructType int

const (
	GeoConstructNingSkill GeoConstructType = iota
	GeoConstructZhongliSkill
	GeoConstructTravellerSkill
	GeoConstructTravellerBurst
	GeoConstructAlbedoSkill
)

type Construct interface {
	OnDestruct()
	Key() int
	Type() GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
}
