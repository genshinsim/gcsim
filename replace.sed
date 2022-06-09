#!/bin/sed -f

s;"github\.com/genshinsim/gcsim/internal/tmpl/character";tmpl "github.com/genshinsim/gcsim/internal/template/character" \
"github.com/genshinsim/gcsim/pkg/core/player/character";
s/\*character\.Tmpl/*tmpl.Character/

s/Register\(.*\)Func(core\./Register\1Func(keys./

s/p core\.CharacterProfile/w *character.CharWrapper, p character.CharacterProfile/
s/(core\.Character, error)/error/
s;t, err := character\.NewTemplateChar(s, p);t := tmpl.New(s) \
t.CharWrapper = w;
s;c\.Tmpl = t;c.Character = t \
;
s;return &c, nil;w.Character = \&c \
 \
return nil;
s/func (c \*char) Init() {/func (c *char) Init() error {/
s/c\.Tmpl\.Init()//

s/(int, int)/action.ActionInfo/

s/Core\.Events\.Subscribe(core\./Core.Events.Subscribe(event./

s/core\.Action\(\w*\)/action.Action\1/g
s/core\.AttackTag\(\w*\)/combat.AttackTag\1/g
s/core\.ICD\(\w*\)/combat.ICD\1/g
s/core\.Log\(\w*\)Event/glog.Log\1Event/g
s/core\.StrikeType\(\w*\)/combat.StrikeType\1/g
s/core\.WeaponClass\(\w*\)/weapon.WeaponClass\1/g
s/core\.Zone\(\w*\)/character.Zone\1/g
s/core\.New\(\w*\)Hit/combat.New\1Hit/g
s/core\.NewDefSingleTarget/combat.NewDefSingleTarget/g

s/AddTask/Core.Tasks.Add/g
s/Core\.Chars/Core.Player.Chars()/g
s/Core\.Combat\.QueueAttack/Core.QueueAttack/g
s/Core\.Health\.Heal/Core.Player.Heal/g
s/Core\.Status\.AddStatus/Core.Status.Add/g
s/Core\.Status\.DeleteStatus/Core.Status.Delete/g
s/Tmpl\.ActionReady/Character.ActionReady/g
s/Tmpl\.ActionStam/Character.ActionStam/g
s/Tmpl\.Snapshot/Character.Snapshot/g
s/\b\(\w\)\.QueueParticle/\1.Core.QueueParticle/g
s/core\.AttackCB/combat.AttackCB/g
s/core\.AttackEvent/combat.AttackEvent/g
s/core\.AttackInfo/combat.AttackInfo/g
s/core\.HealInfo/player.HealInfo/g
s/core\.Snapshot/combat.Snapshot/g
s/core\.Target/combat.Target/g

s/MaxEnergy()/EnergyMax/g

s/core\.EleType\b/attributes.Element/g
s/core\.Electro\b/attributes.Electro/g
s/core\.Pyro\b/attributes.Pyro/g
s/core\.Cryo\b/attributes.Cryo/g
s/core\.Hydro\b/attributes.Hydro/g
s/core\.Frozen\b/attributes.Frozen/g
s/core\.Anemo\b/attributes.Anemo/g
s/core\.Dendro\b/attributes.Dendro/g
s/core\.Geo\b/attributes.Geo/g
s/core\.NoElement\b/attributes.NoElement/g
s/core\.Physical\b/attributes.Physical/g

s/core\.DEFP\b/attributes.DEFP/g
s/core\.DEF\b/attributes.DEF/g
s/core\.HP\b/attributes.HP/g
s/core\.HPP\b/attributes.HPP/g
s/core\.ATK\b/attributes.ATK/g
s/core\.ATKP\b/attributes.ATKP/g
s/core\.ER\b/attributes.ER/g
s/core\.EM\b/attributes.EM/g
s/core\.CR\b/attributes.CR/g
s/core\.CD\b/attributes.CD/g
s/core\.Heal\b/attributes.Heal/g
s/core\.PyroP\b/attributes.PyroP/g
s/core\.HydroP\b/attributes.HydroP/g
s/core\.CryoP\b/attributes.CryoP/g
s/core\.ElectroP\b/attributes.ElectroP/g
s/core\.AnemoP\b/attributes.AnemoP/g
s/core\.GeoP\b/attributes.GeoP/g
s/core\.PhyP\b/attributes.PhyP/g
s/core\.DendroP\b/attributes.DendroP/g
s/core\.AtkSpd\b/attributes.AtkSpd/g
s/core\.DmgP\b/attributes.DmgP/g
s/core\.EndStatType\b/attributes.EndStatType/g

s/action\.ActionType/action.Action/g

# need to convert manually
s/AddMod(core\./AddStatMod(core./
s/AddPreDamageMod(core\./AddAttackMod(core./
s/\b\(\w\)\.AddWeaponInfuse(core\.WeaponInfusion{/\1.Core.Player.AddWeaponInfuse(core.WeaponInfusion{/
s/\b\(\w\)\.Core\.ActiveChar/\1.Core.Player.Active()/g
s/Core\.Player\.Chars()\[\(\w\)\.Core\.Player\.Active()\]/\1.Core.Player.ActiveChar()/g
s/Core\.Player\.Chars()\[\(\w\)\.\(\w\)\.Core\.Player\.Active()\]/\1.\1.Core.Player.ActiveChar()/g

s;Amount:;attributes.NoStat, \
;
s/Expiry://
s/Key://