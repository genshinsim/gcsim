import os
import re

replace = [
    [r"core\.AttackInfo", "combat.AttackInfo"],
    [r"core\.AttackTag(\w+)", "attacks.AttackTag\\1"],
    [
        r"func weapon\(char core\.Character, c \*core\.Core, r int, param map\[string\]int\) string {",
        'type Weapon struct {\n'
        '   Index int\n'
        '}\n'
        'func (w * Weapon) SetIndex(idx int) {w.Index = idx}\n'
        'func (w * Weapon) Init() error      {return nil}\n\n'
        'func NewWeapon(c * core.Core, char * character.CharWrapper, p weapon.WeaponProfile)(weapon.Weapon, error) {\n'
        '   w := &Weapon{}\n'
        '   r := p.Refine\n'
    ],
    [r"char\.CharIndex\(\)", "char.Index"],
    [r"c\.ActiveChar", "c.Player.Active()"],
    [r"char\.Name\(\)", "char.Base.Key.String()"],
    [r"char\.CurrentEnergy\(\)", "char.Energy"],
    [r"char\.MaxEnergy\(\)", "char.EnergyMax"],
    [r"\.Status\.AddStatus\(", ".Status.Add("],
    [r"core\.RegisterWeaponFunc\(\"(.+)\", weapon\)", ""],
    [r"return \"(.+)\"", "return w, nil"],
    [r"c\.Combat\.Queue", "c.Queue"]
]


files = [f for f in os.listdir('.') if os.path.isfile(f)]
for f in files:
    print(f"Processing {f}")
    with open(f, "r+") as file:
        data = file.read()
        file.seek(0)
        # replace using regex
        for x in replace:
            data = re.sub(x[0], x[1], data)
        # print(data)
        file.write(data)
        file.truncate()
