/* eslint-disable */

export enum AvatarCurveType {
  INVALID_AVATAR_CURVE = 0,
  GROW_CURVE_HP_S4 = 1,
  GROW_CURVE_ATTACK_S4 = 2,
  GROW_CURVE_HP_S5 = 3,
  GROW_CURVE_ATTACK_S5 = 4,
  UNRECOGNIZED = -1,
}

export function avatarCurveTypeFromJSON(object: any): AvatarCurveType {
  switch (object) {
    case 0:
    case "INVALID_AVATAR_CURVE":
      return AvatarCurveType.INVALID_AVATAR_CURVE;
    case 1:
    case "GROW_CURVE_HP_S4":
      return AvatarCurveType.GROW_CURVE_HP_S4;
    case 2:
    case "GROW_CURVE_ATTACK_S4":
      return AvatarCurveType.GROW_CURVE_ATTACK_S4;
    case 3:
    case "GROW_CURVE_HP_S5":
      return AvatarCurveType.GROW_CURVE_HP_S5;
    case 4:
    case "GROW_CURVE_ATTACK_S5":
      return AvatarCurveType.GROW_CURVE_ATTACK_S5;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AvatarCurveType.UNRECOGNIZED;
  }
}

export function avatarCurveTypeToJSON(object: AvatarCurveType): string {
  switch (object) {
    case AvatarCurveType.INVALID_AVATAR_CURVE:
      return "INVALID_AVATAR_CURVE";
    case AvatarCurveType.GROW_CURVE_HP_S4:
      return "GROW_CURVE_HP_S4";
    case AvatarCurveType.GROW_CURVE_ATTACK_S4:
      return "GROW_CURVE_ATTACK_S4";
    case AvatarCurveType.GROW_CURVE_HP_S5:
      return "GROW_CURVE_HP_S5";
    case AvatarCurveType.GROW_CURVE_ATTACK_S5:
      return "GROW_CURVE_ATTACK_S5";
    case AvatarCurveType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum QualityType {
  INVALID_QUALITY_TYPE = 0,
  /** QUALITY_ORANGE_SP - the special 6 star aka aloy */
  QUALITY_ORANGE_SP = 6,
  QUALITY_ORANGE = 5,
  QUALITY_PURPLE = 4,
  UNRECOGNIZED = -1,
}

export function qualityTypeFromJSON(object: any): QualityType {
  switch (object) {
    case 0:
    case "INVALID_QUALITY_TYPE":
      return QualityType.INVALID_QUALITY_TYPE;
    case 6:
    case "QUALITY_ORANGE_SP":
      return QualityType.QUALITY_ORANGE_SP;
    case 5:
    case "QUALITY_ORANGE":
      return QualityType.QUALITY_ORANGE;
    case 4:
    case "QUALITY_PURPLE":
      return QualityType.QUALITY_PURPLE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return QualityType.UNRECOGNIZED;
  }
}

export function qualityTypeToJSON(object: QualityType): string {
  switch (object) {
    case QualityType.INVALID_QUALITY_TYPE:
      return "INVALID_QUALITY_TYPE";
    case QualityType.QUALITY_ORANGE_SP:
      return "QUALITY_ORANGE_SP";
    case QualityType.QUALITY_ORANGE:
      return "QUALITY_ORANGE";
    case QualityType.QUALITY_PURPLE:
      return "QUALITY_PURPLE";
    case QualityType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum WeaponCurveType {
  INVALID_WEAPON_CURVE = 0,
  GROW_CURVE_ATTACK_101 = 1,
  GROW_CURVE_ATTACK_102 = 2,
  GROW_CURVE_ATTACK_103 = 3,
  GROW_CURVE_ATTACK_104 = 4,
  GROW_CURVE_ATTACK_105 = 5,
  GROW_CURVE_CRITICAL_101 = 6,
  GROW_CURVE_ATTACK_201 = 7,
  GROW_CURVE_ATTACK_202 = 8,
  GROW_CURVE_ATTACK_203 = 9,
  GROW_CURVE_ATTACK_204 = 10,
  GROW_CURVE_ATTACK_205 = 11,
  GROW_CURVE_CRITICAL_201 = 12,
  GROW_CURVE_ATTACK_301 = 13,
  GROW_CURVE_ATTACK_302 = 14,
  GROW_CURVE_ATTACK_303 = 15,
  GROW_CURVE_ATTACK_304 = 16,
  GROW_CURVE_ATTACK_305 = 17,
  GROW_CURVE_CRITICAL_301 = 18,
  UNRECOGNIZED = -1,
}

export function weaponCurveTypeFromJSON(object: any): WeaponCurveType {
  switch (object) {
    case 0:
    case "INVALID_WEAPON_CURVE":
      return WeaponCurveType.INVALID_WEAPON_CURVE;
    case 1:
    case "GROW_CURVE_ATTACK_101":
      return WeaponCurveType.GROW_CURVE_ATTACK_101;
    case 2:
    case "GROW_CURVE_ATTACK_102":
      return WeaponCurveType.GROW_CURVE_ATTACK_102;
    case 3:
    case "GROW_CURVE_ATTACK_103":
      return WeaponCurveType.GROW_CURVE_ATTACK_103;
    case 4:
    case "GROW_CURVE_ATTACK_104":
      return WeaponCurveType.GROW_CURVE_ATTACK_104;
    case 5:
    case "GROW_CURVE_ATTACK_105":
      return WeaponCurveType.GROW_CURVE_ATTACK_105;
    case 6:
    case "GROW_CURVE_CRITICAL_101":
      return WeaponCurveType.GROW_CURVE_CRITICAL_101;
    case 7:
    case "GROW_CURVE_ATTACK_201":
      return WeaponCurveType.GROW_CURVE_ATTACK_201;
    case 8:
    case "GROW_CURVE_ATTACK_202":
      return WeaponCurveType.GROW_CURVE_ATTACK_202;
    case 9:
    case "GROW_CURVE_ATTACK_203":
      return WeaponCurveType.GROW_CURVE_ATTACK_203;
    case 10:
    case "GROW_CURVE_ATTACK_204":
      return WeaponCurveType.GROW_CURVE_ATTACK_204;
    case 11:
    case "GROW_CURVE_ATTACK_205":
      return WeaponCurveType.GROW_CURVE_ATTACK_205;
    case 12:
    case "GROW_CURVE_CRITICAL_201":
      return WeaponCurveType.GROW_CURVE_CRITICAL_201;
    case 13:
    case "GROW_CURVE_ATTACK_301":
      return WeaponCurveType.GROW_CURVE_ATTACK_301;
    case 14:
    case "GROW_CURVE_ATTACK_302":
      return WeaponCurveType.GROW_CURVE_ATTACK_302;
    case 15:
    case "GROW_CURVE_ATTACK_303":
      return WeaponCurveType.GROW_CURVE_ATTACK_303;
    case 16:
    case "GROW_CURVE_ATTACK_304":
      return WeaponCurveType.GROW_CURVE_ATTACK_304;
    case 17:
    case "GROW_CURVE_ATTACK_305":
      return WeaponCurveType.GROW_CURVE_ATTACK_305;
    case 18:
    case "GROW_CURVE_CRITICAL_301":
      return WeaponCurveType.GROW_CURVE_CRITICAL_301;
    case -1:
    case "UNRECOGNIZED":
    default:
      return WeaponCurveType.UNRECOGNIZED;
  }
}

export function weaponCurveTypeToJSON(object: WeaponCurveType): string {
  switch (object) {
    case WeaponCurveType.INVALID_WEAPON_CURVE:
      return "INVALID_WEAPON_CURVE";
    case WeaponCurveType.GROW_CURVE_ATTACK_101:
      return "GROW_CURVE_ATTACK_101";
    case WeaponCurveType.GROW_CURVE_ATTACK_102:
      return "GROW_CURVE_ATTACK_102";
    case WeaponCurveType.GROW_CURVE_ATTACK_103:
      return "GROW_CURVE_ATTACK_103";
    case WeaponCurveType.GROW_CURVE_ATTACK_104:
      return "GROW_CURVE_ATTACK_104";
    case WeaponCurveType.GROW_CURVE_ATTACK_105:
      return "GROW_CURVE_ATTACK_105";
    case WeaponCurveType.GROW_CURVE_CRITICAL_101:
      return "GROW_CURVE_CRITICAL_101";
    case WeaponCurveType.GROW_CURVE_ATTACK_201:
      return "GROW_CURVE_ATTACK_201";
    case WeaponCurveType.GROW_CURVE_ATTACK_202:
      return "GROW_CURVE_ATTACK_202";
    case WeaponCurveType.GROW_CURVE_ATTACK_203:
      return "GROW_CURVE_ATTACK_203";
    case WeaponCurveType.GROW_CURVE_ATTACK_204:
      return "GROW_CURVE_ATTACK_204";
    case WeaponCurveType.GROW_CURVE_ATTACK_205:
      return "GROW_CURVE_ATTACK_205";
    case WeaponCurveType.GROW_CURVE_CRITICAL_201:
      return "GROW_CURVE_CRITICAL_201";
    case WeaponCurveType.GROW_CURVE_ATTACK_301:
      return "GROW_CURVE_ATTACK_301";
    case WeaponCurveType.GROW_CURVE_ATTACK_302:
      return "GROW_CURVE_ATTACK_302";
    case WeaponCurveType.GROW_CURVE_ATTACK_303:
      return "GROW_CURVE_ATTACK_303";
    case WeaponCurveType.GROW_CURVE_ATTACK_304:
      return "GROW_CURVE_ATTACK_304";
    case WeaponCurveType.GROW_CURVE_ATTACK_305:
      return "GROW_CURVE_ATTACK_305";
    case WeaponCurveType.GROW_CURVE_CRITICAL_301:
      return "GROW_CURVE_CRITICAL_301";
    case WeaponCurveType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum WeaponClass {
  INVALID_WEAPON_CLASS = 0,
  WEAPON_SWORD_ONE_HAND = 1,
  WEAPON_CLAYMORE = 2,
  WEAPON_POLE = 3,
  WEAPON_BOW = 4,
  WEAPON_CATALYST = 5,
  UNRECOGNIZED = -1,
}

export function weaponClassFromJSON(object: any): WeaponClass {
  switch (object) {
    case 0:
    case "INVALID_WEAPON_CLASS":
      return WeaponClass.INVALID_WEAPON_CLASS;
    case 1:
    case "WEAPON_SWORD_ONE_HAND":
      return WeaponClass.WEAPON_SWORD_ONE_HAND;
    case 2:
    case "WEAPON_CLAYMORE":
      return WeaponClass.WEAPON_CLAYMORE;
    case 3:
    case "WEAPON_POLE":
      return WeaponClass.WEAPON_POLE;
    case 4:
    case "WEAPON_BOW":
      return WeaponClass.WEAPON_BOW;
    case 5:
    case "WEAPON_CATALYST":
      return WeaponClass.WEAPON_CATALYST;
    case -1:
    case "UNRECOGNIZED":
    default:
      return WeaponClass.UNRECOGNIZED;
  }
}

export function weaponClassToJSON(object: WeaponClass): string {
  switch (object) {
    case WeaponClass.INVALID_WEAPON_CLASS:
      return "INVALID_WEAPON_CLASS";
    case WeaponClass.WEAPON_SWORD_ONE_HAND:
      return "WEAPON_SWORD_ONE_HAND";
    case WeaponClass.WEAPON_CLAYMORE:
      return "WEAPON_CLAYMORE";
    case WeaponClass.WEAPON_POLE:
      return "WEAPON_POLE";
    case WeaponClass.WEAPON_BOW:
      return "WEAPON_BOW";
    case WeaponClass.WEAPON_CATALYST:
      return "WEAPON_CATALYST";
    case WeaponClass.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum MonsterCurveType {
  INVALID_MONSTER_CURVE = 0,
  GROW_CURVE_HP = 1,
  GROW_CURVE_HP_2 = 2,
  GROW_CURVE_HP_ENVIRONMENT = 3,
  UNRECOGNIZED = -1,
}

export function monsterCurveTypeFromJSON(object: any): MonsterCurveType {
  switch (object) {
    case 0:
    case "INVALID_MONSTER_CURVE":
      return MonsterCurveType.INVALID_MONSTER_CURVE;
    case 1:
    case "GROW_CURVE_HP":
      return MonsterCurveType.GROW_CURVE_HP;
    case 2:
    case "GROW_CURVE_HP_2":
      return MonsterCurveType.GROW_CURVE_HP_2;
    case 3:
    case "GROW_CURVE_HP_ENVIRONMENT":
      return MonsterCurveType.GROW_CURVE_HP_ENVIRONMENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return MonsterCurveType.UNRECOGNIZED;
  }
}

export function monsterCurveTypeToJSON(object: MonsterCurveType): string {
  switch (object) {
    case MonsterCurveType.INVALID_MONSTER_CURVE:
      return "INVALID_MONSTER_CURVE";
    case MonsterCurveType.GROW_CURVE_HP:
      return "GROW_CURVE_HP";
    case MonsterCurveType.GROW_CURVE_HP_2:
      return "GROW_CURVE_HP_2";
    case MonsterCurveType.GROW_CURVE_HP_ENVIRONMENT:
      return "GROW_CURVE_HP_ENVIRONMENT";
    case MonsterCurveType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum BodyType {
  INVALID_BODY_TYPE = 0,
  BODY_UNKNOWN = 1,
  BODY_BOY = 2,
  BODY_GIRL = 3,
  BODY_MALE = 4,
  BODY_LADY = 5,
  BODY_LOLI = 6,
  UNRECOGNIZED = -1,
}

export function bodyTypeFromJSON(object: any): BodyType {
  switch (object) {
    case 0:
    case "INVALID_BODY_TYPE":
      return BodyType.INVALID_BODY_TYPE;
    case 1:
    case "BODY_UNKNOWN":
      return BodyType.BODY_UNKNOWN;
    case 2:
    case "BODY_BOY":
      return BodyType.BODY_BOY;
    case 3:
    case "BODY_GIRL":
      return BodyType.BODY_GIRL;
    case 4:
    case "BODY_MALE":
      return BodyType.BODY_MALE;
    case 5:
    case "BODY_LADY":
      return BodyType.BODY_LADY;
    case 6:
    case "BODY_LOLI":
      return BodyType.BODY_LOLI;
    case -1:
    case "UNRECOGNIZED":
    default:
      return BodyType.UNRECOGNIZED;
  }
}

export function bodyTypeToJSON(object: BodyType): string {
  switch (object) {
    case BodyType.INVALID_BODY_TYPE:
      return "INVALID_BODY_TYPE";
    case BodyType.BODY_UNKNOWN:
      return "BODY_UNKNOWN";
    case BodyType.BODY_BOY:
      return "BODY_BOY";
    case BodyType.BODY_GIRL:
      return "BODY_GIRL";
    case BodyType.BODY_MALE:
      return "BODY_MALE";
    case BodyType.BODY_LADY:
      return "BODY_LADY";
    case BodyType.BODY_LOLI:
      return "BODY_LOLI";
    case BodyType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ZoneType {
  INVALID_ZONE_TYPE = 0,
  ASSOC_TYPE_UNKNOWN = 1,
  ASSOC_TYPE_MONDSTADT = 2,
  ASSOC_TYPE_LIYUE = 3,
  ASSOC_TYPE_INAZUMA = 4,
  ASSOC_TYPE_SUMERU = 5,
  ASSOC_TYPE_FATUI = 6,
  /** ASSOC_TYPE_RANGER - aloy pls */
  ASSOC_TYPE_RANGER = 7,
  /** ASSOC_TYPE_MAINACTOR - traveler is cool */
  ASSOC_TYPE_MAINACTOR = 8,
  ASSOC_TYPE_FONTAINE = 9,
  UNRECOGNIZED = -1,
}

export function zoneTypeFromJSON(object: any): ZoneType {
  switch (object) {
    case 0:
    case "INVALID_ZONE_TYPE":
      return ZoneType.INVALID_ZONE_TYPE;
    case 1:
    case "ASSOC_TYPE_UNKNOWN":
      return ZoneType.ASSOC_TYPE_UNKNOWN;
    case 2:
    case "ASSOC_TYPE_MONDSTADT":
      return ZoneType.ASSOC_TYPE_MONDSTADT;
    case 3:
    case "ASSOC_TYPE_LIYUE":
      return ZoneType.ASSOC_TYPE_LIYUE;
    case 4:
    case "ASSOC_TYPE_INAZUMA":
      return ZoneType.ASSOC_TYPE_INAZUMA;
    case 5:
    case "ASSOC_TYPE_SUMERU":
      return ZoneType.ASSOC_TYPE_SUMERU;
    case 6:
    case "ASSOC_TYPE_FATUI":
      return ZoneType.ASSOC_TYPE_FATUI;
    case 7:
    case "ASSOC_TYPE_RANGER":
      return ZoneType.ASSOC_TYPE_RANGER;
    case 8:
    case "ASSOC_TYPE_MAINACTOR":
      return ZoneType.ASSOC_TYPE_MAINACTOR;
    case 9:
    case "ASSOC_TYPE_FONTAINE":
      return ZoneType.ASSOC_TYPE_FONTAINE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ZoneType.UNRECOGNIZED;
  }
}

export function zoneTypeToJSON(object: ZoneType): string {
  switch (object) {
    case ZoneType.INVALID_ZONE_TYPE:
      return "INVALID_ZONE_TYPE";
    case ZoneType.ASSOC_TYPE_UNKNOWN:
      return "ASSOC_TYPE_UNKNOWN";
    case ZoneType.ASSOC_TYPE_MONDSTADT:
      return "ASSOC_TYPE_MONDSTADT";
    case ZoneType.ASSOC_TYPE_LIYUE:
      return "ASSOC_TYPE_LIYUE";
    case ZoneType.ASSOC_TYPE_INAZUMA:
      return "ASSOC_TYPE_INAZUMA";
    case ZoneType.ASSOC_TYPE_SUMERU:
      return "ASSOC_TYPE_SUMERU";
    case ZoneType.ASSOC_TYPE_FATUI:
      return "ASSOC_TYPE_FATUI";
    case ZoneType.ASSOC_TYPE_RANGER:
      return "ASSOC_TYPE_RANGER";
    case ZoneType.ASSOC_TYPE_MAINACTOR:
      return "ASSOC_TYPE_MAINACTOR";
    case ZoneType.ASSOC_TYPE_FONTAINE:
      return "ASSOC_TYPE_FONTAINE";
    case ZoneType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum Element {
  INVALID_ELEMENT = 0,
  Electric = 1,
  Fire = 2,
  Ice = 3,
  Water = 4,
  Grass = 5,
  ELEMENT_QUICKEN = 6,
  ELEMENT_FROZEN = 7,
  Wind = 8,
  Rock = 9,
  ELEMENT_NONE = 10,
  ELEMENT_PHYSICAL = 11,
  ELEMENT_UNKNOWN = 12,
  UNRECOGNIZED = -1,
}

export function elementFromJSON(object: any): Element {
  switch (object) {
    case 0:
    case "INVALID_ELEMENT":
      return Element.INVALID_ELEMENT;
    case 1:
    case "Electric":
      return Element.Electric;
    case 2:
    case "Fire":
      return Element.Fire;
    case 3:
    case "Ice":
      return Element.Ice;
    case 4:
    case "Water":
      return Element.Water;
    case 5:
    case "Grass":
      return Element.Grass;
    case 6:
    case "ELEMENT_QUICKEN":
      return Element.ELEMENT_QUICKEN;
    case 7:
    case "ELEMENT_FROZEN":
      return Element.ELEMENT_FROZEN;
    case 8:
    case "Wind":
      return Element.Wind;
    case 9:
    case "Rock":
      return Element.Rock;
    case 10:
    case "ELEMENT_NONE":
      return Element.ELEMENT_NONE;
    case 11:
    case "ELEMENT_PHYSICAL":
      return Element.ELEMENT_PHYSICAL;
    case 12:
    case "ELEMENT_UNKNOWN":
      return Element.ELEMENT_UNKNOWN;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Element.UNRECOGNIZED;
  }
}

export function elementToJSON(object: Element): string {
  switch (object) {
    case Element.INVALID_ELEMENT:
      return "INVALID_ELEMENT";
    case Element.Electric:
      return "Electric";
    case Element.Fire:
      return "Fire";
    case Element.Ice:
      return "Ice";
    case Element.Water:
      return "Water";
    case Element.Grass:
      return "Grass";
    case Element.ELEMENT_QUICKEN:
      return "ELEMENT_QUICKEN";
    case Element.ELEMENT_FROZEN:
      return "ELEMENT_FROZEN";
    case Element.Wind:
      return "Wind";
    case Element.Rock:
      return "Rock";
    case Element.ELEMENT_NONE:
      return "ELEMENT_NONE";
    case Element.ELEMENT_PHYSICAL:
      return "ELEMENT_PHYSICAL";
    case Element.ELEMENT_UNKNOWN:
      return "ELEMENT_UNKNOWN";
    case Element.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum StatType {
  INVALID_STAT_TYPE = 0,
  FIGHT_PROP_DEFENSE_PERCENT = 1,
  FIGHT_PROP_DEFENSE = 2,
  FIGHT_PROP_HP = 3,
  FIGHT_PROP_HP_PERCENT = 4,
  FIGHT_PROP_ATTACK = 5,
  FIGHT_PROP_ATTACK_PERCENT = 6,
  FIGHT_PROP_CHARGE_EFFICIENCY = 7,
  FIGHT_PROP_ELEMENT_MASTERY = 8,
  FIGHT_PROP_CRITICAL = 9,
  FIGHT_PROP_CRITICAL_HURT = 10,
  FIGHT_PROP_HEAL_ADD = 11,
  FIGHT_PROP_FIRE_ADD_HURT = 12,
  FIGHT_PROP_WATER_ADD_HURT = 13,
  FIGHT_PROP_GRASS_ADD_HURT = 14,
  FIGHT_PROP_ELEC_ADD_HURT = 15,
  FIGHT_PROP_WIND_ADD_HURT = 16,
  FIGHT_PROP_ICE_ADD_HURT = 17,
  FIGHT_PROP_ROCK_ADD_HURT = 18,
  FIGHT_PROP_PHYSICAL_ADD_HURT = 19,
  FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT = 20,
  /** FIGHT_PROP_HEALED_ADD - healing bonus */
  FIGHT_PROP_HEALED_ADD = 21,
  /** FIGHT_PROP_BASE_HP - base hp */
  FIGHT_PROP_BASE_HP = 22,
  /** FIGHT_PROP_BASE_ATTACK - base attack */
  FIGHT_PROP_BASE_ATTACK = 23,
  /** FIGHT_PROP_BASE_DEFENSE - base defense */
  FIGHT_PROP_BASE_DEFENSE = 24,
  /** FIGHT_PROP_MAX_HP - max hp; not really used? */
  FIGHT_PROP_MAX_HP = 25,
  UNRECOGNIZED = -1,
}

export function statTypeFromJSON(object: any): StatType {
  switch (object) {
    case 0:
    case "INVALID_STAT_TYPE":
      return StatType.INVALID_STAT_TYPE;
    case 1:
    case "FIGHT_PROP_DEFENSE_PERCENT":
      return StatType.FIGHT_PROP_DEFENSE_PERCENT;
    case 2:
    case "FIGHT_PROP_DEFENSE":
      return StatType.FIGHT_PROP_DEFENSE;
    case 3:
    case "FIGHT_PROP_HP":
      return StatType.FIGHT_PROP_HP;
    case 4:
    case "FIGHT_PROP_HP_PERCENT":
      return StatType.FIGHT_PROP_HP_PERCENT;
    case 5:
    case "FIGHT_PROP_ATTACK":
      return StatType.FIGHT_PROP_ATTACK;
    case 6:
    case "FIGHT_PROP_ATTACK_PERCENT":
      return StatType.FIGHT_PROP_ATTACK_PERCENT;
    case 7:
    case "FIGHT_PROP_CHARGE_EFFICIENCY":
      return StatType.FIGHT_PROP_CHARGE_EFFICIENCY;
    case 8:
    case "FIGHT_PROP_ELEMENT_MASTERY":
      return StatType.FIGHT_PROP_ELEMENT_MASTERY;
    case 9:
    case "FIGHT_PROP_CRITICAL":
      return StatType.FIGHT_PROP_CRITICAL;
    case 10:
    case "FIGHT_PROP_CRITICAL_HURT":
      return StatType.FIGHT_PROP_CRITICAL_HURT;
    case 11:
    case "FIGHT_PROP_HEAL_ADD":
      return StatType.FIGHT_PROP_HEAL_ADD;
    case 12:
    case "FIGHT_PROP_FIRE_ADD_HURT":
      return StatType.FIGHT_PROP_FIRE_ADD_HURT;
    case 13:
    case "FIGHT_PROP_WATER_ADD_HURT":
      return StatType.FIGHT_PROP_WATER_ADD_HURT;
    case 14:
    case "FIGHT_PROP_GRASS_ADD_HURT":
      return StatType.FIGHT_PROP_GRASS_ADD_HURT;
    case 15:
    case "FIGHT_PROP_ELEC_ADD_HURT":
      return StatType.FIGHT_PROP_ELEC_ADD_HURT;
    case 16:
    case "FIGHT_PROP_WIND_ADD_HURT":
      return StatType.FIGHT_PROP_WIND_ADD_HURT;
    case 17:
    case "FIGHT_PROP_ICE_ADD_HURT":
      return StatType.FIGHT_PROP_ICE_ADD_HURT;
    case 18:
    case "FIGHT_PROP_ROCK_ADD_HURT":
      return StatType.FIGHT_PROP_ROCK_ADD_HURT;
    case 19:
    case "FIGHT_PROP_PHYSICAL_ADD_HURT":
      return StatType.FIGHT_PROP_PHYSICAL_ADD_HURT;
    case 20:
    case "FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT":
      return StatType.FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT;
    case 21:
    case "FIGHT_PROP_HEALED_ADD":
      return StatType.FIGHT_PROP_HEALED_ADD;
    case 22:
    case "FIGHT_PROP_BASE_HP":
      return StatType.FIGHT_PROP_BASE_HP;
    case 23:
    case "FIGHT_PROP_BASE_ATTACK":
      return StatType.FIGHT_PROP_BASE_ATTACK;
    case 24:
    case "FIGHT_PROP_BASE_DEFENSE":
      return StatType.FIGHT_PROP_BASE_DEFENSE;
    case 25:
    case "FIGHT_PROP_MAX_HP":
      return StatType.FIGHT_PROP_MAX_HP;
    case -1:
    case "UNRECOGNIZED":
    default:
      return StatType.UNRECOGNIZED;
  }
}

export function statTypeToJSON(object: StatType): string {
  switch (object) {
    case StatType.INVALID_STAT_TYPE:
      return "INVALID_STAT_TYPE";
    case StatType.FIGHT_PROP_DEFENSE_PERCENT:
      return "FIGHT_PROP_DEFENSE_PERCENT";
    case StatType.FIGHT_PROP_DEFENSE:
      return "FIGHT_PROP_DEFENSE";
    case StatType.FIGHT_PROP_HP:
      return "FIGHT_PROP_HP";
    case StatType.FIGHT_PROP_HP_PERCENT:
      return "FIGHT_PROP_HP_PERCENT";
    case StatType.FIGHT_PROP_ATTACK:
      return "FIGHT_PROP_ATTACK";
    case StatType.FIGHT_PROP_ATTACK_PERCENT:
      return "FIGHT_PROP_ATTACK_PERCENT";
    case StatType.FIGHT_PROP_CHARGE_EFFICIENCY:
      return "FIGHT_PROP_CHARGE_EFFICIENCY";
    case StatType.FIGHT_PROP_ELEMENT_MASTERY:
      return "FIGHT_PROP_ELEMENT_MASTERY";
    case StatType.FIGHT_PROP_CRITICAL:
      return "FIGHT_PROP_CRITICAL";
    case StatType.FIGHT_PROP_CRITICAL_HURT:
      return "FIGHT_PROP_CRITICAL_HURT";
    case StatType.FIGHT_PROP_HEAL_ADD:
      return "FIGHT_PROP_HEAL_ADD";
    case StatType.FIGHT_PROP_FIRE_ADD_HURT:
      return "FIGHT_PROP_FIRE_ADD_HURT";
    case StatType.FIGHT_PROP_WATER_ADD_HURT:
      return "FIGHT_PROP_WATER_ADD_HURT";
    case StatType.FIGHT_PROP_GRASS_ADD_HURT:
      return "FIGHT_PROP_GRASS_ADD_HURT";
    case StatType.FIGHT_PROP_ELEC_ADD_HURT:
      return "FIGHT_PROP_ELEC_ADD_HURT";
    case StatType.FIGHT_PROP_WIND_ADD_HURT:
      return "FIGHT_PROP_WIND_ADD_HURT";
    case StatType.FIGHT_PROP_ICE_ADD_HURT:
      return "FIGHT_PROP_ICE_ADD_HURT";
    case StatType.FIGHT_PROP_ROCK_ADD_HURT:
      return "FIGHT_PROP_ROCK_ADD_HURT";
    case StatType.FIGHT_PROP_PHYSICAL_ADD_HURT:
      return "FIGHT_PROP_PHYSICAL_ADD_HURT";
    case StatType.FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT:
      return "FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT";
    case StatType.FIGHT_PROP_HEALED_ADD:
      return "FIGHT_PROP_HEALED_ADD";
    case StatType.FIGHT_PROP_BASE_HP:
      return "FIGHT_PROP_BASE_HP";
    case StatType.FIGHT_PROP_BASE_ATTACK:
      return "FIGHT_PROP_BASE_ATTACK";
    case StatType.FIGHT_PROP_BASE_DEFENSE:
      return "FIGHT_PROP_BASE_DEFENSE";
    case StatType.FIGHT_PROP_MAX_HP:
      return "FIGHT_PROP_MAX_HP";
    case StatType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SimMode {
  INVALID_SIM_MODE = 0,
  DURATION_MODE = 1,
  TTK_MODE = 2,
  UNRECOGNIZED = -1,
}

export function simModeFromJSON(object: any): SimMode {
  switch (object) {
    case 0:
    case "INVALID_SIM_MODE":
      return SimMode.INVALID_SIM_MODE;
    case 1:
    case "DURATION_MODE":
      return SimMode.DURATION_MODE;
    case 2:
    case "TTK_MODE":
      return SimMode.TTK_MODE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SimMode.UNRECOGNIZED;
  }
}

export function simModeToJSON(object: SimMode): string {
  switch (object) {
    case SimMode.INVALID_SIM_MODE:
      return "INVALID_SIM_MODE";
    case SimMode.DURATION_MODE:
      return "DURATION_MODE";
    case SimMode.TTK_MODE:
      return "TTK_MODE";
    case SimMode.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ComputeWorkSource {
  InvalidWork = 0,
  DBWork = 1,
  SubmissionWork = 2,
  UNRECOGNIZED = -1,
}

export function computeWorkSourceFromJSON(object: any): ComputeWorkSource {
  switch (object) {
    case 0:
    case "InvalidWork":
      return ComputeWorkSource.InvalidWork;
    case 1:
    case "DBWork":
      return ComputeWorkSource.DBWork;
    case 2:
    case "SubmissionWork":
      return ComputeWorkSource.SubmissionWork;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ComputeWorkSource.UNRECOGNIZED;
  }
}

export function computeWorkSourceToJSON(object: ComputeWorkSource): string {
  switch (object) {
    case ComputeWorkSource.InvalidWork:
      return "InvalidWork";
    case ComputeWorkSource.DBWork:
      return "DBWork";
    case ComputeWorkSource.SubmissionWork:
      return "SubmissionWork";
    case ComputeWorkSource.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum DBTag {
  DB_TAG_INVALID = 0,
  DB_TAG_GCSIM = 1,
  DB_TAG_TESTING = 2,
  DB_TAG_ITTO_SIMPS = 5,
  DB_TAG_RANDOM_DELAYS = 6,
  /** DB_TAG_ARFOIRE_NEWBIES - reissue4917 (1070054618895233035) tag for newbie players */
  DB_TAG_ARFOIRE_NEWBIES = 7,
  DB_TAG_APL = 8,
  DB_TAG_GUIDES = 9,
  DB_TAG_ADMIN_DO_NOT_USE = 99999999,
  UNRECOGNIZED = -1,
}

export function dBTagFromJSON(object: any): DBTag {
  switch (object) {
    case 0:
    case "DB_TAG_INVALID":
      return DBTag.DB_TAG_INVALID;
    case 1:
    case "DB_TAG_GCSIM":
      return DBTag.DB_TAG_GCSIM;
    case 2:
    case "DB_TAG_TESTING":
      return DBTag.DB_TAG_TESTING;
    case 5:
    case "DB_TAG_ITTO_SIMPS":
      return DBTag.DB_TAG_ITTO_SIMPS;
    case 6:
    case "DB_TAG_RANDOM_DELAYS":
      return DBTag.DB_TAG_RANDOM_DELAYS;
    case 7:
    case "DB_TAG_ARFOIRE_NEWBIES":
      return DBTag.DB_TAG_ARFOIRE_NEWBIES;
    case 8:
    case "DB_TAG_APL":
      return DBTag.DB_TAG_APL;
    case 9:
    case "DB_TAG_GUIDES":
      return DBTag.DB_TAG_GUIDES;
    case 99999999:
    case "DB_TAG_ADMIN_DO_NOT_USE":
      return DBTag.DB_TAG_ADMIN_DO_NOT_USE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return DBTag.UNRECOGNIZED;
  }
}

export function dBTagToJSON(object: DBTag): string {
  switch (object) {
    case DBTag.DB_TAG_INVALID:
      return "DB_TAG_INVALID";
    case DBTag.DB_TAG_GCSIM:
      return "DB_TAG_GCSIM";
    case DBTag.DB_TAG_TESTING:
      return "DB_TAG_TESTING";
    case DBTag.DB_TAG_ITTO_SIMPS:
      return "DB_TAG_ITTO_SIMPS";
    case DBTag.DB_TAG_RANDOM_DELAYS:
      return "DB_TAG_RANDOM_DELAYS";
    case DBTag.DB_TAG_ARFOIRE_NEWBIES:
      return "DB_TAG_ARFOIRE_NEWBIES";
    case DBTag.DB_TAG_APL:
      return "DB_TAG_APL";
    case DBTag.DB_TAG_GUIDES:
      return "DB_TAG_GUIDES";
    case DBTag.DB_TAG_ADMIN_DO_NOT_USE:
      return "DB_TAG_ADMIN_DO_NOT_USE";
    case DBTag.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
