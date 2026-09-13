package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/model"
)

func ConvertProtoStat(in model.FightPropType) attributes.Stat {
	if out, ok := map[model.FightPropType]attributes.Stat{
		model.FightPropType_FIGHT_PROP_BASE_HP:           attributes.BaseHP,
		model.FightPropType_FIGHT_PROP_HP:                attributes.HP,
		model.FightPropType_FIGHT_PROP_HP_PERCENT:        attributes.HPP,
		model.FightPropType_FIGHT_PROP_BASE_ATTACK:       attributes.BaseATK,
		model.FightPropType_FIGHT_PROP_ATTACK:            attributes.ATK,
		model.FightPropType_FIGHT_PROP_ATTACK_PERCENT:    attributes.ATKP,
		model.FightPropType_FIGHT_PROP_BASE_DEFENSE:      attributes.BaseDEF,
		model.FightPropType_FIGHT_PROP_DEFENSE:           attributes.DEF,
		model.FightPropType_FIGHT_PROP_DEFENSE_PERCENT:   attributes.DEFP,
		model.FightPropType_FIGHT_PROP_CRITICAL:          attributes.CR,
		model.FightPropType_FIGHT_PROP_CRITICAL_HURT:     attributes.CD,
		model.FightPropType_FIGHT_PROP_CHARGE_EFFICIENCY: attributes.ER,
		model.FightPropType_FIGHT_PROP_HEAL_ADD:          attributes.Heal,
		model.FightPropType_FIGHT_PROP_ELEMENT_MASTERY:   attributes.EM,
		model.FightPropType_FIGHT_PROP_PHYSICAL_ADD_HURT: attributes.PhyP,
		model.FightPropType_FIGHT_PROP_FIRE_ADD_HURT:     attributes.PyroP,
		model.FightPropType_FIGHT_PROP_ELEC_ADD_HURT:     attributes.ElectroP,
		model.FightPropType_FIGHT_PROP_WATER_ADD_HURT:    attributes.HydroP,
		model.FightPropType_FIGHT_PROP_GRASS_ADD_HURT:    attributes.DendroP,
		model.FightPropType_FIGHT_PROP_WIND_ADD_HURT:     attributes.AnemoP,
		model.FightPropType_FIGHT_PROP_ROCK_ADD_HURT:     attributes.GeoP,
		model.FightPropType_FIGHT_PROP_ICE_ADD_HURT:      attributes.CryoP,
	}[in]; ok {
		return out
	}
	return attributes.NoStat
}

func ConvertProtoElement(in model.ElementType) attributes.Element {
	if out, ok := map[model.ElementType]attributes.Element{
		model.ElementType_Fire:     attributes.Pyro,
		model.ElementType_Water:    attributes.Hydro,
		model.ElementType_Grass:    attributes.Dendro,
		model.ElementType_Electric: attributes.Electro,
		model.ElementType_Ice:      attributes.Cryo,
		model.ElementType_Wind:     attributes.Anemo,
		model.ElementType_Rock:     attributes.Geo,
	}[in]; ok {
		return out
	}
	return attributes.NoElement
}

func ConvertRegion(in model.AssocType) ZoneType {
	if out, ok := map[model.AssocType]ZoneType{
		model.AssocType_ASSOC_TYPE_MONDSTADT:      ZoneMondstadt,
		model.AssocType_ASSOC_TYPE_LIYUE:          ZoneLiyue,
		model.AssocType_ASSOC_TYPE_INAZUMA:        ZoneInazuma,
		model.AssocType_ASSOC_TYPE_SUMERU:         ZoneSumeru,
		model.AssocType_ASSOC_TYPE_FONTAINE:       ZoneFontaine,
		model.AssocType_ASSOC_TYPE_NATLAN:         ZoneNatlan,
		model.AssocType_ASSOC_TYPE_SNEZHNAYA:      ZoneSnezhnaya,
		model.AssocType_ASSOC_TYPE_NODKRAI:        ZoneNodKrai,
		model.AssocType_ASSOC_TYPE_NODKRAI_ZIBAI:  ZoneNodKrai, // TODO: is zibai not liyue?
		model.AssocType_ASSOC_TYPE_SNEZHNAYA_STAR: ZoneSnezhnaya,
	}[in]; ok {
		return out
	}
	return ZoneUnknown
}

func ConvertWeaponClass(in model.WeaponType) WeaponClass {
	if out, ok := map[model.WeaponType]WeaponClass{
		model.WeaponType_WEAPON_SWORD_ONE_HAND: WeaponClassSword,
		model.WeaponType_WEAPON_CATALYST:       WeaponClassCatalyst,
		model.WeaponType_WEAPON_CLAYMORE:       WeaponClassClaymore,
		model.WeaponType_WEAPON_BOW:            WeaponClassBow,
		model.WeaponType_WEAPON_POLE:           WeaponClassSpear,
	}[in]; ok {
		return out
	}
	// TODO: no invalid
	return WeaponClassSword
}

func ConvertBodyType(in model.BodyType) BodyType {
	if out, ok := map[model.BodyType]BodyType{
		model.BodyType_BODY_BOY:  BodyBoy,
		model.BodyType_BODY_GIRL: BodyGirl,
		model.BodyType_BODY_LADY: BodyLady,
		model.BodyType_BODY_MALE: BodyMale,
		model.BodyType_BODY_LOLI: BodyLoli,
	}[in]; ok {
		return out
	}
	// TODO: no invalid
	return BodyMale
}

func ConvertRarity(in model.QualityType) int {
	if out, ok := map[model.QualityType]int{
		model.QualityType_QUALITY_WHITE:     1,
		model.QualityType_QUALITY_GREEN:     2,
		model.QualityType_QUALITY_BLUE:      3,
		model.QualityType_QUALITY_PURPLE:    4,
		model.QualityType_QUALITY_ORANGE:    5,
		model.QualityType_QUALITY_ORANGE_SP: 5,
	}[in]; ok {
		return out
	}
	return 0
}
