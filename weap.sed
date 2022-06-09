s;func weapon\(char core\.Character, c \*core\.Core, r int, param map\[string\]int\) string {;\
type Weap struct { \
Index int\
}\
func (w *Weapon) SetIndex(idx int) { w.Index = idx}\
func (w *Weapon) Init() error      { return nil }\
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {\
;