import os
import re

replace = [
    [r"core\.AttackInfo", "combat.AttackInfo"],
    [r"core\.AttackTag(\w+)", "combat.AttackTag\\1"],
    [
        r"func weapon\(char core\.Character, c \*core\.Core, r int, param map\[string\]int\) string {",
        'type Weapon struct {\n'
        'Index int\n'
        '}\n'
        'func (w * Weapon) SetIndex(idx int) {w.Index = idx}\n'
        'func(w * Weapon) Init() error      {return nil}\n'
        'func NewWeapon(c * core.Core, char * character.CharWrapper, p weapon.WeaponProfile)(weapon.Weapon, error) {\n'
        'w := &Weapon{}\n'
        'r := p.Refine\n'
    ],
    [r"core\.EleType\b", "attributes.Element"],
    [r"core\.Electro\b", "attributes.Electro"],
    [r"core\.Pyro\b", "attributes.Pyro"],
    [r"core\.Cryo\b", "attributes.Cryo"],
    [r"core\.Hydro\b", "attributes.Hydro"],
    [r"core\.Frozen\b", "attributes.Frozen"],
    [r"core\.Anemo\b", "attributes.Anemo"],
    [r"core\.Dendro\b", "attributes.Dendro"],
    [r"core\.Geo\b", "attributes.Geo"],
    [r"core\.NoElement\b", "attributes.NoElement"],
    [r"core\.Physical\b", "attributes.Physical"],
    [r"core\.DEFP\b", "attributes.DEFP"],
    [r"core\.DEF\b", "attributes.DEF"],
    [r"core\.HP\b", "attributes.HP"],
    [r"core\.HPP\b", "attributes.HPP"],
    [r"core\.ATK\b", "attributes.ATK"],
    [r"core\.ATKP\b", "attributes.ATKP"],
    [r"core\.ER\b", "attributes.ER"],
    [r"core\.EM\b", "attributes.EM"],
    [r"core\.CR\b", "attributes.CR"],
    [r"core\.CD\b", "attributes.CD"],
    [r"core\.Heal\b", "attributes.Heal"],
    [r"core\.PyroP\b", "attributes.PyroP"],
    [r"core\.HydroP\b", "attributes.HydroP"],
    [r"core\.CryoP\b", "attributes.CryoP"],
    [r"core\.ElectroP\b", "attributes.ElectroP"],
    [r"core\.AnemoP\b", "attributes.AnemoP"],
    [r"core\.GeoP\b", "attributes.GeoP"],
    [r"core\.PhyP\b", "attributes.PhyP"],
    [r"core\.DendroP\b", "attributes.DendroP"],
    [r"core\.AtkSpd\b", "attributes.AtkSpd"],
    [r"core\.DmgP\b", "attributes.DmgP"],
    [r"core\.EndStatType\b", "attributes.EndStatType"],
    [r"core\.OnAttackWillLand\b", "event.OnAttackWillLand"],
    [r"core\.OnDamage\b", "event.OnDamage"],
    [r"core\.OnAuraDurabilityAdded\b", "event.OnAuraDurabilityAdded"],
    [r"core\.OnAuraDurabilityDepleted\b", "event.OnAuraDurabilityDepleted"],
    [r"core\.ReactionEventStartDeli\b", "event.ReactionEventStartDeli"],
    [r"core\.OnOverload\b", "event.OnOverload"],
    [r"core\.OnSuperconduct\b", "event.OnSuperconduct"],
    [r"core\.OnMelt\b", "event.OnMelt"],
    [r"core\.OnVaporize\b", "event.OnVaporize"],
    [r"core\.OnFrozen\b", "event.OnFrozen"],
    [r"core\.OnElectroCharged\b", "event.OnElectroCharged"],
    [r"core\.OnSwirlHydr\b", "event.OnSwirlHydr"],
    [r"core\.OnSwirlCry\b", "event.OnSwirlCry"],
    [r"core\.OnSwirlElectr\b", "event.OnSwirlElectr"],
    [r"core\.OnSwirlPyr\b", "event.OnSwirlPyr"],
    [r"core\.OnCrystallizeHydr\b", "event.OnCrystallizeHydr"],
    [r"core\.OnCrystallizeCry\b", "event.OnCrystallizeCry"],
    [r"core\.OnCrystallizeElectr\b", "event.OnCrystallizeElectr"],
    [r"core\.OnCrystallizePyr\b", "event.OnCrystallizePyr"],
    [r"core\.ReactionEventEndDeli\b", "event.ReactionEventEndDeli"],
    [r"core\.OnStamUse\b", "event.OnStamUse"],
    [r"core\.OnShielded\b", "event.OnShielded"],
    [r"core\.OnCharacterSwap\b", "event.OnCharacterSwap"],
    [r"core\.OnParticleReceived\b", "event.OnParticleReceived"],
    [r"core\.OnEnergyChange\b", "event.OnEnergyChange"],
    [r"core\.OnTargetDied\b", "event.OnTargetDied"],
    [r"core\.OnCharacterHurt\b", "event.OnCharacterHurt"],
    [r"core\.OnHeal\b", "event.OnHeal"],
    [r"core\.OnActionExec\b", "event.OnActionExec"],
    [r"core\.PreSkill\b", "event.PreSkill"],
    [r"core\.PostSkill\b", "event.PostSkill"],
    [r"core\.PreBurst\b", "event.PreBurst"],
    [r"core\.PostBurst\b", "event.PostBurst"],
    [r"core\.PreAttack\b", "event.PreAttack"],
    [r"core\.PostAttack\b", "event.PostAttack"],
    [r"core\.PreChargeAttack\b", "event.PreChargeAttack"],
    [r"core\.PostChargeAttack\b", "event.PostChargeAttack"],
    [r"core\.PrePlunge\b", "event.PrePlunge"],
    [r"core\.PostPlunge\b", "event.PostPlunge"],
    [r"core\.PreAimShoot\b", "event.PreAimShoot"],
    [r"core\.PostAimShoot\b", "event.PostAimShoot"],
    [r"core\.PreDas\b", "event.PreDas"],
    [r"core\.PostDas\b", "event.PostDas"],
    [r"core\.OnInitialize\b", "event.OnInitialize"],
    [r"core\.OnStateChange\b", "event.OnStateChange"],
    [r"core\.OnTargetAdded\b", "event.OnTargetAdded"],
    [r"core\.EndEventTypes\b", "event.EndEventTypes"],
    [r"core\.AttackEvent\b", "combat.AttackEvent"],
    [r"core\.Target\b", "combat.Target"],
    [r"AddPreDamageMod\(core\.PreDamageMod{", "AddAttackMod("],
    [r"\.AddMod\(core\.CharStatMod{", ".AddStatMod("],
    [r"char\.CharIndex\(\)", "char.Index"],
    [r"c\.ActiveChar", "c.Player.Active()"],
    [r"char\.Name\(\)", "char.Base.Name"],
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
