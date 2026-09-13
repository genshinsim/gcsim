package attributes

//go:generate go tool github.com/dmarkham/enumer -text -json -linecomment -type=Element $GOFILE
type Element int

// ElementType should be Pyro, Hydro, Cryo, Electro, Geo, Anemo and maybe Dendro
const (
	NoElement Element = iota // none
	Pyro                     // pyro
	Hydro                    // hydro
	Dendro                   // dendro
	Electro                  // electro
	Cryo                     // cryo
	Anemo                    // anemo
	Geo                      // geo
	Physical                 // physical
	Frozen                   // frozen
	Quicken                  // quicken
)

func EleToDmgP(e Element) Stat {
	switch e {
	case Anemo:
		return AnemoP
	case Cryo:
		return CryoP
	case Electro:
		return ElectroP
	case Geo:
		return GeoP
	case Hydro:
		return HydroP
	case Pyro:
		return PyroP
	case Dendro:
		return DendroP
	case Physical:
		return PhyP
	}
	return NoStat
}
