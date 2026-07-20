/* eslint-disable */

export enum QualityType {
  QUALITY_NONE = 0,
  QUALITY_WHITE = 1,
  QUALITY_GREEN = 2,
  QUALITY_BLUE = 3,
  QUALITY_PURPLE = 4,
  QUALITY_ORANGE = 5,
  /** QUALITY_ORANGE_SP - the special 6 star aka aloy */
  QUALITY_ORANGE_SP = 105,
  UNRECOGNIZED = -1,
}

export function qualityTypeFromJSON(object: any): QualityType {
  switch (object) {
    case 0:
    case "QUALITY_NONE":
      return QualityType.QUALITY_NONE;
    case 1:
    case "QUALITY_WHITE":
      return QualityType.QUALITY_WHITE;
    case 2:
    case "QUALITY_GREEN":
      return QualityType.QUALITY_GREEN;
    case 3:
    case "QUALITY_BLUE":
      return QualityType.QUALITY_BLUE;
    case 4:
    case "QUALITY_PURPLE":
      return QualityType.QUALITY_PURPLE;
    case 5:
    case "QUALITY_ORANGE":
      return QualityType.QUALITY_ORANGE;
    case 105:
    case "QUALITY_ORANGE_SP":
      return QualityType.QUALITY_ORANGE_SP;
    case -1:
    case "UNRECOGNIZED":
    default:
      return QualityType.UNRECOGNIZED;
  }
}

export function qualityTypeToJSON(object: QualityType): string {
  switch (object) {
    case QualityType.QUALITY_NONE:
      return "QUALITY_NONE";
    case QualityType.QUALITY_WHITE:
      return "QUALITY_WHITE";
    case QualityType.QUALITY_GREEN:
      return "QUALITY_GREEN";
    case QualityType.QUALITY_BLUE:
      return "QUALITY_BLUE";
    case QualityType.QUALITY_PURPLE:
      return "QUALITY_PURPLE";
    case QualityType.QUALITY_ORANGE:
      return "QUALITY_ORANGE";
    case QualityType.QUALITY_ORANGE_SP:
      return "QUALITY_ORANGE_SP";
    case QualityType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum WeaponType {
  WEAPON_NONE = 0,
  WEAPON_SWORD_ONE_HAND = 1,
  /**
   * WEAPON_CATALYST - WEAPON_CROSSBOW = 2;
   * WEAPON_STAFF = 3;
   * WEAPON_DOUBLE_DAGGER = 4;
   * WEAPON_KATANA = 5;
   * WEAPON_SHURIKEN = 6;
   * WEAPON_STICK = 7;
   * WEAPON_SPEAR = 8;
   * WEAPON_SHIELD_SMALL = 9;
   */
  WEAPON_CATALYST = 10,
  WEAPON_CLAYMORE = 11,
  WEAPON_BOW = 12,
  WEAPON_POLE = 13,
  UNRECOGNIZED = -1,
}

export function weaponTypeFromJSON(object: any): WeaponType {
  switch (object) {
    case 0:
    case "WEAPON_NONE":
      return WeaponType.WEAPON_NONE;
    case 1:
    case "WEAPON_SWORD_ONE_HAND":
      return WeaponType.WEAPON_SWORD_ONE_HAND;
    case 10:
    case "WEAPON_CATALYST":
      return WeaponType.WEAPON_CATALYST;
    case 11:
    case "WEAPON_CLAYMORE":
      return WeaponType.WEAPON_CLAYMORE;
    case 12:
    case "WEAPON_BOW":
      return WeaponType.WEAPON_BOW;
    case 13:
    case "WEAPON_POLE":
      return WeaponType.WEAPON_POLE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return WeaponType.UNRECOGNIZED;
  }
}

export function weaponTypeToJSON(object: WeaponType): string {
  switch (object) {
    case WeaponType.WEAPON_NONE:
      return "WEAPON_NONE";
    case WeaponType.WEAPON_SWORD_ONE_HAND:
      return "WEAPON_SWORD_ONE_HAND";
    case WeaponType.WEAPON_CATALYST:
      return "WEAPON_CATALYST";
    case WeaponType.WEAPON_CLAYMORE:
      return "WEAPON_CLAYMORE";
    case WeaponType.WEAPON_BOW:
      return "WEAPON_BOW";
    case WeaponType.WEAPON_POLE:
      return "WEAPON_POLE";
    case WeaponType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum BodyType {
  BODY_NONE = 0,
  BODY_BOY = 1,
  BODY_GIRL = 2,
  BODY_LADY = 3,
  BODY_MALE = 4,
  BODY_LOLI = 5,
  UNRECOGNIZED = -1,
}

export function bodyTypeFromJSON(object: any): BodyType {
  switch (object) {
    case 0:
    case "BODY_NONE":
      return BodyType.BODY_NONE;
    case 1:
    case "BODY_BOY":
      return BodyType.BODY_BOY;
    case 2:
    case "BODY_GIRL":
      return BodyType.BODY_GIRL;
    case 3:
    case "BODY_LADY":
      return BodyType.BODY_LADY;
    case 4:
    case "BODY_MALE":
      return BodyType.BODY_MALE;
    case 5:
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
    case BodyType.BODY_NONE:
      return "BODY_NONE";
    case BodyType.BODY_BOY:
      return "BODY_BOY";
    case BodyType.BODY_GIRL:
      return "BODY_GIRL";
    case BodyType.BODY_LADY:
      return "BODY_LADY";
    case BodyType.BODY_MALE:
      return "BODY_MALE";
    case BodyType.BODY_LOLI:
      return "BODY_LOLI";
    case BodyType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum AssocType {
  ASSOC_TYPE_NONE = 0,
  ASSOC_TYPE_MONDSTADT = 1,
  ASSOC_TYPE_LIYUE = 2,
  /** ASSOC_TYPE_MAINACTOR - traveler/manekin */
  ASSOC_TYPE_MAINACTOR = 3,
  ASSOC_TYPE_FATUI = 4,
  ASSOC_TYPE_INAZUMA = 5,
  /** ASSOC_TYPE_RANGER - aloy pls */
  ASSOC_TYPE_RANGER = 6,
  ASSOC_TYPE_SUMERU = 7,
  ASSOC_TYPE_FONTAINE = 8,
  ASSOC_TYPE_NATLAN = 9,
  ASSOC_TYPE_SNEZHNAYA = 10,
  ASSOC_TYPE_OMNI_SCOURGE = 11,
  ASSOC_TYPE_NODKRAI = 12,
  ASSOC_TYPE_NODKRAI_ZIBAI = 13,
  ASSOC_TYPE_HVISION = 14,
  ASSOC_TYPE_SNEZHNAYA_STAR = 15,
  UNRECOGNIZED = -1,
}

export function assocTypeFromJSON(object: any): AssocType {
  switch (object) {
    case 0:
    case "ASSOC_TYPE_NONE":
      return AssocType.ASSOC_TYPE_NONE;
    case 1:
    case "ASSOC_TYPE_MONDSTADT":
      return AssocType.ASSOC_TYPE_MONDSTADT;
    case 2:
    case "ASSOC_TYPE_LIYUE":
      return AssocType.ASSOC_TYPE_LIYUE;
    case 3:
    case "ASSOC_TYPE_MAINACTOR":
      return AssocType.ASSOC_TYPE_MAINACTOR;
    case 4:
    case "ASSOC_TYPE_FATUI":
      return AssocType.ASSOC_TYPE_FATUI;
    case 5:
    case "ASSOC_TYPE_INAZUMA":
      return AssocType.ASSOC_TYPE_INAZUMA;
    case 6:
    case "ASSOC_TYPE_RANGER":
      return AssocType.ASSOC_TYPE_RANGER;
    case 7:
    case "ASSOC_TYPE_SUMERU":
      return AssocType.ASSOC_TYPE_SUMERU;
    case 8:
    case "ASSOC_TYPE_FONTAINE":
      return AssocType.ASSOC_TYPE_FONTAINE;
    case 9:
    case "ASSOC_TYPE_NATLAN":
      return AssocType.ASSOC_TYPE_NATLAN;
    case 10:
    case "ASSOC_TYPE_SNEZHNAYA":
      return AssocType.ASSOC_TYPE_SNEZHNAYA;
    case 11:
    case "ASSOC_TYPE_OMNI_SCOURGE":
      return AssocType.ASSOC_TYPE_OMNI_SCOURGE;
    case 12:
    case "ASSOC_TYPE_NODKRAI":
      return AssocType.ASSOC_TYPE_NODKRAI;
    case 13:
    case "ASSOC_TYPE_NODKRAI_ZIBAI":
      return AssocType.ASSOC_TYPE_NODKRAI_ZIBAI;
    case 14:
    case "ASSOC_TYPE_HVISION":
      return AssocType.ASSOC_TYPE_HVISION;
    case 15:
    case "ASSOC_TYPE_SNEZHNAYA_STAR":
      return AssocType.ASSOC_TYPE_SNEZHNAYA_STAR;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AssocType.UNRECOGNIZED;
  }
}

export function assocTypeToJSON(object: AssocType): string {
  switch (object) {
    case AssocType.ASSOC_TYPE_NONE:
      return "ASSOC_TYPE_NONE";
    case AssocType.ASSOC_TYPE_MONDSTADT:
      return "ASSOC_TYPE_MONDSTADT";
    case AssocType.ASSOC_TYPE_LIYUE:
      return "ASSOC_TYPE_LIYUE";
    case AssocType.ASSOC_TYPE_MAINACTOR:
      return "ASSOC_TYPE_MAINACTOR";
    case AssocType.ASSOC_TYPE_FATUI:
      return "ASSOC_TYPE_FATUI";
    case AssocType.ASSOC_TYPE_INAZUMA:
      return "ASSOC_TYPE_INAZUMA";
    case AssocType.ASSOC_TYPE_RANGER:
      return "ASSOC_TYPE_RANGER";
    case AssocType.ASSOC_TYPE_SUMERU:
      return "ASSOC_TYPE_SUMERU";
    case AssocType.ASSOC_TYPE_FONTAINE:
      return "ASSOC_TYPE_FONTAINE";
    case AssocType.ASSOC_TYPE_NATLAN:
      return "ASSOC_TYPE_NATLAN";
    case AssocType.ASSOC_TYPE_SNEZHNAYA:
      return "ASSOC_TYPE_SNEZHNAYA";
    case AssocType.ASSOC_TYPE_OMNI_SCOURGE:
      return "ASSOC_TYPE_OMNI_SCOURGE";
    case AssocType.ASSOC_TYPE_NODKRAI:
      return "ASSOC_TYPE_NODKRAI";
    case AssocType.ASSOC_TYPE_NODKRAI_ZIBAI:
      return "ASSOC_TYPE_NODKRAI_ZIBAI";
    case AssocType.ASSOC_TYPE_HVISION:
      return "ASSOC_TYPE_HVISION";
    case AssocType.ASSOC_TYPE_SNEZHNAYA_STAR:
      return "ASSOC_TYPE_SNEZHNAYA_STAR";
    case AssocType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ElementType {
  None = 0,
  Fire = 1,
  Water = 2,
  Grass = 3,
  Electric = 4,
  Ice = 5,
  Frozen = 6,
  Wind = 7,
  Rock = 8,
  /**
   * Mushroom - AntiFire = 9;
   * VehicleMuteIce = 10;
   */
  Mushroom = 11,
  /**
   * Overdose - Wood = 13;
   * LiquidPhlogiston = 14;
   * SolidPhlogiston = 15;
   * SolidifyPhlogiston = 16;
   */
  Overdose = 12,
  UNRECOGNIZED = -1,
}

export function elementTypeFromJSON(object: any): ElementType {
  switch (object) {
    case 0:
    case "None":
      return ElementType.None;
    case 1:
    case "Fire":
      return ElementType.Fire;
    case 2:
    case "Water":
      return ElementType.Water;
    case 3:
    case "Grass":
      return ElementType.Grass;
    case 4:
    case "Electric":
      return ElementType.Electric;
    case 5:
    case "Ice":
      return ElementType.Ice;
    case 6:
    case "Frozen":
      return ElementType.Frozen;
    case 7:
    case "Wind":
      return ElementType.Wind;
    case 8:
    case "Rock":
      return ElementType.Rock;
    case 11:
    case "Mushroom":
      return ElementType.Mushroom;
    case 12:
    case "Overdose":
      return ElementType.Overdose;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ElementType.UNRECOGNIZED;
  }
}

export function elementTypeToJSON(object: ElementType): string {
  switch (object) {
    case ElementType.None:
      return "None";
    case ElementType.Fire:
      return "Fire";
    case ElementType.Water:
      return "Water";
    case ElementType.Grass:
      return "Grass";
    case ElementType.Electric:
      return "Electric";
    case ElementType.Ice:
      return "Ice";
    case ElementType.Frozen:
      return "Frozen";
    case ElementType.Wind:
      return "Wind";
    case ElementType.Rock:
      return "Rock";
    case ElementType.Mushroom:
      return "Mushroom";
    case ElementType.Overdose:
      return "Overdose";
    case ElementType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum EquipType {
  /** EQUIP_NONE - none */
  EQUIP_NONE = 0,
  /** EQUIP_BRACER - flower */
  EQUIP_BRACER = 1,
  /** EQUIP_NECKLACE - plume */
  EQUIP_NECKLACE = 2,
  /** EQUIP_SHOES - sands */
  EQUIP_SHOES = 3,
  /** EQUIP_RING - goblet */
  EQUIP_RING = 4,
  /** EQUIP_DRESS - circlet */
  EQUIP_DRESS = 5,
  /** EQUIP_WEAPON - weapon */
  EQUIP_WEAPON = 6,
  UNRECOGNIZED = -1,
}

export function equipTypeFromJSON(object: any): EquipType {
  switch (object) {
    case 0:
    case "EQUIP_NONE":
      return EquipType.EQUIP_NONE;
    case 1:
    case "EQUIP_BRACER":
      return EquipType.EQUIP_BRACER;
    case 2:
    case "EQUIP_NECKLACE":
      return EquipType.EQUIP_NECKLACE;
    case 3:
    case "EQUIP_SHOES":
      return EquipType.EQUIP_SHOES;
    case 4:
    case "EQUIP_RING":
      return EquipType.EQUIP_RING;
    case 5:
    case "EQUIP_DRESS":
      return EquipType.EQUIP_DRESS;
    case 6:
    case "EQUIP_WEAPON":
      return EquipType.EQUIP_WEAPON;
    case -1:
    case "UNRECOGNIZED":
    default:
      return EquipType.UNRECOGNIZED;
  }
}

export function equipTypeToJSON(object: EquipType): string {
  switch (object) {
    case EquipType.EQUIP_NONE:
      return "EQUIP_NONE";
    case EquipType.EQUIP_BRACER:
      return "EQUIP_BRACER";
    case EquipType.EQUIP_NECKLACE:
      return "EQUIP_NECKLACE";
    case EquipType.EQUIP_SHOES:
      return "EQUIP_SHOES";
    case EquipType.EQUIP_RING:
      return "EQUIP_RING";
    case EquipType.EQUIP_DRESS:
      return "EQUIP_DRESS";
    case EquipType.EQUIP_WEAPON:
      return "EQUIP_WEAPON";
    case EquipType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum FightPropType {
  FIGHT_PROP_NONE = 0,
  FIGHT_PROP_BASE_HP = 1,
  FIGHT_PROP_HP = 2,
  FIGHT_PROP_HP_PERCENT = 3,
  FIGHT_PROP_BASE_ATTACK = 4,
  FIGHT_PROP_ATTACK = 5,
  FIGHT_PROP_ATTACK_PERCENT = 6,
  FIGHT_PROP_BASE_DEFENSE = 7,
  FIGHT_PROP_DEFENSE = 8,
  FIGHT_PROP_DEFENSE_PERCENT = 9,
  FIGHT_PROP_BASE_SPEED = 10,
  FIGHT_PROP_SPEED_PERCENT = 11,
  FIGHT_PROP_HP_MP_PERCENT = 12,
  FIGHT_PROP_ATTACK_MP_PERCENT = 13,
  FIGHT_PROP_CRITICAL = 20,
  FIGHT_PROP_ANTI_CRITICAL = 21,
  FIGHT_PROP_CRITICAL_HURT = 22,
  FIGHT_PROP_CHARGE_EFFICIENCY = 23,
  FIGHT_PROP_ADD_HURT = 24,
  FIGHT_PROP_SUB_HURT = 25,
  FIGHT_PROP_HEAL_ADD = 26,
  FIGHT_PROP_HEALED_ADD = 27,
  FIGHT_PROP_ELEMENT_MASTERY = 28,
  FIGHT_PROP_PHYSICAL_SUB_HURT = 29,
  FIGHT_PROP_PHYSICAL_ADD_HURT = 30,
  FIGHT_PROP_DEFENCE_IGNORE_RATIO = 31,
  FIGHT_PROP_DEFENCE_IGNORE_DELTA = 32,
  FIGHT_PROP_FIRE_ADD_HURT = 40,
  FIGHT_PROP_ELEC_ADD_HURT = 41,
  FIGHT_PROP_WATER_ADD_HURT = 42,
  FIGHT_PROP_GRASS_ADD_HURT = 43,
  FIGHT_PROP_WIND_ADD_HURT = 44,
  FIGHT_PROP_ROCK_ADD_HURT = 45,
  FIGHT_PROP_ICE_ADD_HURT = 46,
  FIGHT_PROP_HIT_HEAD_ADD_HURT = 47,
  FIGHT_PROP_FIRE_SUB_HURT = 50,
  FIGHT_PROP_ELEC_SUB_HURT = 51,
  FIGHT_PROP_WATER_SUB_HURT = 52,
  FIGHT_PROP_GRASS_SUB_HURT = 53,
  FIGHT_PROP_WIND_SUB_HURT = 54,
  FIGHT_PROP_ROCK_SUB_HURT = 55,
  FIGHT_PROP_ICE_SUB_HURT = 56,
  FIGHT_PROP_EFFECT_HIT = 60,
  FIGHT_PROP_EFFECT_RESIST = 61,
  FIGHT_PROP_FREEZE_RESIST = 62,
  FIGHT_PROP_DIZZY_RESIST = 64,
  FIGHT_PROP_FREEZE_SHORTEN = 65,
  FIGHT_PROP_DIZZY_SHORTEN = 67,
  FIGHT_PROP_MAX_FIRE_ENERGY = 70,
  FIGHT_PROP_MAX_ELEC_ENERGY = 71,
  FIGHT_PROP_MAX_WATER_ENERGY = 72,
  FIGHT_PROP_MAX_GRASS_ENERGY = 73,
  FIGHT_PROP_MAX_WIND_ENERGY = 74,
  FIGHT_PROP_MAX_ICE_ENERGY = 75,
  FIGHT_PROP_MAX_ROCK_ENERGY = 76,
  FIGHT_PROP_MAX_SPECIAL_ENERGY = 77,
  FIGHT_PROP_START_SPECIAL_ENERGY = 78,
  FIGHT_PROP_SKILL_CD_MINUS_RATIO = 80,
  FIGHT_PROP_SHIELD_COST_MINUS_RATIO = 81,
  FIGHT_PROP_CUR_FIRE_ENERGY = 1000,
  FIGHT_PROP_CUR_ELEC_ENERGY = 1001,
  FIGHT_PROP_CUR_WATER_ENERGY = 1002,
  FIGHT_PROP_CUR_GRASS_ENERGY = 1003,
  FIGHT_PROP_CUR_WIND_ENERGY = 1004,
  FIGHT_PROP_CUR_ICE_ENERGY = 1005,
  FIGHT_PROP_CUR_ROCK_ENERGY = 1006,
  FIGHT_PROP_CUR_SPECIAL_ENERGY = 1007,
  FIGHT_PROP_CUR_HP = 1010,
  FIGHT_PROP_MAX_HP = 2000,
  FIGHT_PROP_CUR_ATTACK = 2001,
  FIGHT_PROP_CUR_DEFENSE = 2002,
  FIGHT_PROP_CUR_SPEED = 2003,
  FIGHT_PROP_CUR_HP_DEBTS = 2004,
  FIGHT_PROP_CUR_HP_PAID_DEBTS = 2005,
  FIGHT_PROP_CUR_NATLAN_HP = 2006,
  FIGHT_PROP_NONEXTRA_ATTACK = 3000,
  FIGHT_PROP_NONEXTRA_DEFENSE = 3001,
  FIGHT_PROP_NONEXTRA_CRITICAL = 3002,
  FIGHT_PROP_NONEXTRA_ANTI_CRITICAL = 3003,
  FIGHT_PROP_NONEXTRA_CRITICAL_HURT = 3004,
  FIGHT_PROP_NONEXTRA_CHARGE_EFFICIENCY = 3005,
  FIGHT_PROP_NONEXTRA_ELEMENT_MASTERY = 3006,
  FIGHT_PROP_NONEXTRA_PHYSICAL_SUB_HURT = 3007,
  FIGHT_PROP_NONEXTRA_FIRE_ADD_HURT = 3008,
  FIGHT_PROP_NONEXTRA_ELEC_ADD_HURT = 3009,
  FIGHT_PROP_NONEXTRA_WATER_ADD_HURT = 3010,
  FIGHT_PROP_NONEXTRA_GRASS_ADD_HURT = 3011,
  FIGHT_PROP_NONEXTRA_WIND_ADD_HURT = 3012,
  FIGHT_PROP_NONEXTRA_ROCK_ADD_HURT = 3013,
  FIGHT_PROP_NONEXTRA_ICE_ADD_HURT = 3014,
  FIGHT_PROP_NONEXTRA_FIRE_SUB_HURT = 3015,
  FIGHT_PROP_NONEXTRA_ELEC_SUB_HURT = 3016,
  FIGHT_PROP_NONEXTRA_WATER_SUB_HURT = 3017,
  FIGHT_PROP_NONEXTRA_GRASS_SUB_HURT = 3018,
  FIGHT_PROP_NONEXTRA_WIND_SUB_HURT = 3019,
  FIGHT_PROP_NONEXTRA_ROCK_SUB_HURT = 3020,
  FIGHT_PROP_NONEXTRA_ICE_SUB_HURT = 3021,
  FIGHT_PROP_NONEXTRA_SKILL_CD_MINUS_RATIO = 3022,
  FIGHT_PROP_NONEXTRA_SHIELD_COST_MINUS_RATIO = 3023,
  FIGHT_PROP_NONEXTRA_PHYSICAL_ADD_HURT = 3024,
  FIGHT_PROP_BASE_ELEM_REACT_CRITICAL = 3045,
  FIGHT_PROP_BASE_ELEM_REACT_CRITICAL_HURT = 3046,
  FIGHT_PROP_ELEM_REACT_CRITICAL = 3025,
  FIGHT_PROP_ELEM_REACT_CRITICAL_HURT = 3026,
  FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL = 3027,
  FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL_HURT = 3028,
  FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL = 3029,
  FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL_HURT = 3030,
  FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL = 3031,
  FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL_HURT = 3032,
  FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL = 3033,
  FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL_HURT = 3034,
  FIGHT_PROP_ELEM_REACT_BURN_CRITICAL = 3035,
  FIGHT_PROP_ELEM_REACT_BURN_CRITICAL_HURT = 3036,
  FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL = 3037,
  FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL_HURT = 3038,
  FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL = 3039,
  FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL_HURT = 3040,
  FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL = 3041,
  FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL_HURT = 3042,
  FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL = 3043,
  FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL_HURT = 3044,
  UNRECOGNIZED = -1,
}

export function fightPropTypeFromJSON(object: any): FightPropType {
  switch (object) {
    case 0:
    case "FIGHT_PROP_NONE":
      return FightPropType.FIGHT_PROP_NONE;
    case 1:
    case "FIGHT_PROP_BASE_HP":
      return FightPropType.FIGHT_PROP_BASE_HP;
    case 2:
    case "FIGHT_PROP_HP":
      return FightPropType.FIGHT_PROP_HP;
    case 3:
    case "FIGHT_PROP_HP_PERCENT":
      return FightPropType.FIGHT_PROP_HP_PERCENT;
    case 4:
    case "FIGHT_PROP_BASE_ATTACK":
      return FightPropType.FIGHT_PROP_BASE_ATTACK;
    case 5:
    case "FIGHT_PROP_ATTACK":
      return FightPropType.FIGHT_PROP_ATTACK;
    case 6:
    case "FIGHT_PROP_ATTACK_PERCENT":
      return FightPropType.FIGHT_PROP_ATTACK_PERCENT;
    case 7:
    case "FIGHT_PROP_BASE_DEFENSE":
      return FightPropType.FIGHT_PROP_BASE_DEFENSE;
    case 8:
    case "FIGHT_PROP_DEFENSE":
      return FightPropType.FIGHT_PROP_DEFENSE;
    case 9:
    case "FIGHT_PROP_DEFENSE_PERCENT":
      return FightPropType.FIGHT_PROP_DEFENSE_PERCENT;
    case 10:
    case "FIGHT_PROP_BASE_SPEED":
      return FightPropType.FIGHT_PROP_BASE_SPEED;
    case 11:
    case "FIGHT_PROP_SPEED_PERCENT":
      return FightPropType.FIGHT_PROP_SPEED_PERCENT;
    case 12:
    case "FIGHT_PROP_HP_MP_PERCENT":
      return FightPropType.FIGHT_PROP_HP_MP_PERCENT;
    case 13:
    case "FIGHT_PROP_ATTACK_MP_PERCENT":
      return FightPropType.FIGHT_PROP_ATTACK_MP_PERCENT;
    case 20:
    case "FIGHT_PROP_CRITICAL":
      return FightPropType.FIGHT_PROP_CRITICAL;
    case 21:
    case "FIGHT_PROP_ANTI_CRITICAL":
      return FightPropType.FIGHT_PROP_ANTI_CRITICAL;
    case 22:
    case "FIGHT_PROP_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_CRITICAL_HURT;
    case 23:
    case "FIGHT_PROP_CHARGE_EFFICIENCY":
      return FightPropType.FIGHT_PROP_CHARGE_EFFICIENCY;
    case 24:
    case "FIGHT_PROP_ADD_HURT":
      return FightPropType.FIGHT_PROP_ADD_HURT;
    case 25:
    case "FIGHT_PROP_SUB_HURT":
      return FightPropType.FIGHT_PROP_SUB_HURT;
    case 26:
    case "FIGHT_PROP_HEAL_ADD":
      return FightPropType.FIGHT_PROP_HEAL_ADD;
    case 27:
    case "FIGHT_PROP_HEALED_ADD":
      return FightPropType.FIGHT_PROP_HEALED_ADD;
    case 28:
    case "FIGHT_PROP_ELEMENT_MASTERY":
      return FightPropType.FIGHT_PROP_ELEMENT_MASTERY;
    case 29:
    case "FIGHT_PROP_PHYSICAL_SUB_HURT":
      return FightPropType.FIGHT_PROP_PHYSICAL_SUB_HURT;
    case 30:
    case "FIGHT_PROP_PHYSICAL_ADD_HURT":
      return FightPropType.FIGHT_PROP_PHYSICAL_ADD_HURT;
    case 31:
    case "FIGHT_PROP_DEFENCE_IGNORE_RATIO":
      return FightPropType.FIGHT_PROP_DEFENCE_IGNORE_RATIO;
    case 32:
    case "FIGHT_PROP_DEFENCE_IGNORE_DELTA":
      return FightPropType.FIGHT_PROP_DEFENCE_IGNORE_DELTA;
    case 40:
    case "FIGHT_PROP_FIRE_ADD_HURT":
      return FightPropType.FIGHT_PROP_FIRE_ADD_HURT;
    case 41:
    case "FIGHT_PROP_ELEC_ADD_HURT":
      return FightPropType.FIGHT_PROP_ELEC_ADD_HURT;
    case 42:
    case "FIGHT_PROP_WATER_ADD_HURT":
      return FightPropType.FIGHT_PROP_WATER_ADD_HURT;
    case 43:
    case "FIGHT_PROP_GRASS_ADD_HURT":
      return FightPropType.FIGHT_PROP_GRASS_ADD_HURT;
    case 44:
    case "FIGHT_PROP_WIND_ADD_HURT":
      return FightPropType.FIGHT_PROP_WIND_ADD_HURT;
    case 45:
    case "FIGHT_PROP_ROCK_ADD_HURT":
      return FightPropType.FIGHT_PROP_ROCK_ADD_HURT;
    case 46:
    case "FIGHT_PROP_ICE_ADD_HURT":
      return FightPropType.FIGHT_PROP_ICE_ADD_HURT;
    case 47:
    case "FIGHT_PROP_HIT_HEAD_ADD_HURT":
      return FightPropType.FIGHT_PROP_HIT_HEAD_ADD_HURT;
    case 50:
    case "FIGHT_PROP_FIRE_SUB_HURT":
      return FightPropType.FIGHT_PROP_FIRE_SUB_HURT;
    case 51:
    case "FIGHT_PROP_ELEC_SUB_HURT":
      return FightPropType.FIGHT_PROP_ELEC_SUB_HURT;
    case 52:
    case "FIGHT_PROP_WATER_SUB_HURT":
      return FightPropType.FIGHT_PROP_WATER_SUB_HURT;
    case 53:
    case "FIGHT_PROP_GRASS_SUB_HURT":
      return FightPropType.FIGHT_PROP_GRASS_SUB_HURT;
    case 54:
    case "FIGHT_PROP_WIND_SUB_HURT":
      return FightPropType.FIGHT_PROP_WIND_SUB_HURT;
    case 55:
    case "FIGHT_PROP_ROCK_SUB_HURT":
      return FightPropType.FIGHT_PROP_ROCK_SUB_HURT;
    case 56:
    case "FIGHT_PROP_ICE_SUB_HURT":
      return FightPropType.FIGHT_PROP_ICE_SUB_HURT;
    case 60:
    case "FIGHT_PROP_EFFECT_HIT":
      return FightPropType.FIGHT_PROP_EFFECT_HIT;
    case 61:
    case "FIGHT_PROP_EFFECT_RESIST":
      return FightPropType.FIGHT_PROP_EFFECT_RESIST;
    case 62:
    case "FIGHT_PROP_FREEZE_RESIST":
      return FightPropType.FIGHT_PROP_FREEZE_RESIST;
    case 64:
    case "FIGHT_PROP_DIZZY_RESIST":
      return FightPropType.FIGHT_PROP_DIZZY_RESIST;
    case 65:
    case "FIGHT_PROP_FREEZE_SHORTEN":
      return FightPropType.FIGHT_PROP_FREEZE_SHORTEN;
    case 67:
    case "FIGHT_PROP_DIZZY_SHORTEN":
      return FightPropType.FIGHT_PROP_DIZZY_SHORTEN;
    case 70:
    case "FIGHT_PROP_MAX_FIRE_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_FIRE_ENERGY;
    case 71:
    case "FIGHT_PROP_MAX_ELEC_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_ELEC_ENERGY;
    case 72:
    case "FIGHT_PROP_MAX_WATER_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_WATER_ENERGY;
    case 73:
    case "FIGHT_PROP_MAX_GRASS_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_GRASS_ENERGY;
    case 74:
    case "FIGHT_PROP_MAX_WIND_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_WIND_ENERGY;
    case 75:
    case "FIGHT_PROP_MAX_ICE_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_ICE_ENERGY;
    case 76:
    case "FIGHT_PROP_MAX_ROCK_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_ROCK_ENERGY;
    case 77:
    case "FIGHT_PROP_MAX_SPECIAL_ENERGY":
      return FightPropType.FIGHT_PROP_MAX_SPECIAL_ENERGY;
    case 78:
    case "FIGHT_PROP_START_SPECIAL_ENERGY":
      return FightPropType.FIGHT_PROP_START_SPECIAL_ENERGY;
    case 80:
    case "FIGHT_PROP_SKILL_CD_MINUS_RATIO":
      return FightPropType.FIGHT_PROP_SKILL_CD_MINUS_RATIO;
    case 81:
    case "FIGHT_PROP_SHIELD_COST_MINUS_RATIO":
      return FightPropType.FIGHT_PROP_SHIELD_COST_MINUS_RATIO;
    case 1000:
    case "FIGHT_PROP_CUR_FIRE_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_FIRE_ENERGY;
    case 1001:
    case "FIGHT_PROP_CUR_ELEC_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_ELEC_ENERGY;
    case 1002:
    case "FIGHT_PROP_CUR_WATER_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_WATER_ENERGY;
    case 1003:
    case "FIGHT_PROP_CUR_GRASS_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_GRASS_ENERGY;
    case 1004:
    case "FIGHT_PROP_CUR_WIND_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_WIND_ENERGY;
    case 1005:
    case "FIGHT_PROP_CUR_ICE_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_ICE_ENERGY;
    case 1006:
    case "FIGHT_PROP_CUR_ROCK_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_ROCK_ENERGY;
    case 1007:
    case "FIGHT_PROP_CUR_SPECIAL_ENERGY":
      return FightPropType.FIGHT_PROP_CUR_SPECIAL_ENERGY;
    case 1010:
    case "FIGHT_PROP_CUR_HP":
      return FightPropType.FIGHT_PROP_CUR_HP;
    case 2000:
    case "FIGHT_PROP_MAX_HP":
      return FightPropType.FIGHT_PROP_MAX_HP;
    case 2001:
    case "FIGHT_PROP_CUR_ATTACK":
      return FightPropType.FIGHT_PROP_CUR_ATTACK;
    case 2002:
    case "FIGHT_PROP_CUR_DEFENSE":
      return FightPropType.FIGHT_PROP_CUR_DEFENSE;
    case 2003:
    case "FIGHT_PROP_CUR_SPEED":
      return FightPropType.FIGHT_PROP_CUR_SPEED;
    case 2004:
    case "FIGHT_PROP_CUR_HP_DEBTS":
      return FightPropType.FIGHT_PROP_CUR_HP_DEBTS;
    case 2005:
    case "FIGHT_PROP_CUR_HP_PAID_DEBTS":
      return FightPropType.FIGHT_PROP_CUR_HP_PAID_DEBTS;
    case 2006:
    case "FIGHT_PROP_CUR_NATLAN_HP":
      return FightPropType.FIGHT_PROP_CUR_NATLAN_HP;
    case 3000:
    case "FIGHT_PROP_NONEXTRA_ATTACK":
      return FightPropType.FIGHT_PROP_NONEXTRA_ATTACK;
    case 3001:
    case "FIGHT_PROP_NONEXTRA_DEFENSE":
      return FightPropType.FIGHT_PROP_NONEXTRA_DEFENSE;
    case 3002:
    case "FIGHT_PROP_NONEXTRA_CRITICAL":
      return FightPropType.FIGHT_PROP_NONEXTRA_CRITICAL;
    case 3003:
    case "FIGHT_PROP_NONEXTRA_ANTI_CRITICAL":
      return FightPropType.FIGHT_PROP_NONEXTRA_ANTI_CRITICAL;
    case 3004:
    case "FIGHT_PROP_NONEXTRA_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_CRITICAL_HURT;
    case 3005:
    case "FIGHT_PROP_NONEXTRA_CHARGE_EFFICIENCY":
      return FightPropType.FIGHT_PROP_NONEXTRA_CHARGE_EFFICIENCY;
    case 3006:
    case "FIGHT_PROP_NONEXTRA_ELEMENT_MASTERY":
      return FightPropType.FIGHT_PROP_NONEXTRA_ELEMENT_MASTERY;
    case 3007:
    case "FIGHT_PROP_NONEXTRA_PHYSICAL_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_PHYSICAL_SUB_HURT;
    case 3008:
    case "FIGHT_PROP_NONEXTRA_FIRE_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_FIRE_ADD_HURT;
    case 3009:
    case "FIGHT_PROP_NONEXTRA_ELEC_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ELEC_ADD_HURT;
    case 3010:
    case "FIGHT_PROP_NONEXTRA_WATER_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_WATER_ADD_HURT;
    case 3011:
    case "FIGHT_PROP_NONEXTRA_GRASS_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_GRASS_ADD_HURT;
    case 3012:
    case "FIGHT_PROP_NONEXTRA_WIND_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_WIND_ADD_HURT;
    case 3013:
    case "FIGHT_PROP_NONEXTRA_ROCK_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ROCK_ADD_HURT;
    case 3014:
    case "FIGHT_PROP_NONEXTRA_ICE_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ICE_ADD_HURT;
    case 3015:
    case "FIGHT_PROP_NONEXTRA_FIRE_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_FIRE_SUB_HURT;
    case 3016:
    case "FIGHT_PROP_NONEXTRA_ELEC_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ELEC_SUB_HURT;
    case 3017:
    case "FIGHT_PROP_NONEXTRA_WATER_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_WATER_SUB_HURT;
    case 3018:
    case "FIGHT_PROP_NONEXTRA_GRASS_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_GRASS_SUB_HURT;
    case 3019:
    case "FIGHT_PROP_NONEXTRA_WIND_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_WIND_SUB_HURT;
    case 3020:
    case "FIGHT_PROP_NONEXTRA_ROCK_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ROCK_SUB_HURT;
    case 3021:
    case "FIGHT_PROP_NONEXTRA_ICE_SUB_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_ICE_SUB_HURT;
    case 3022:
    case "FIGHT_PROP_NONEXTRA_SKILL_CD_MINUS_RATIO":
      return FightPropType.FIGHT_PROP_NONEXTRA_SKILL_CD_MINUS_RATIO;
    case 3023:
    case "FIGHT_PROP_NONEXTRA_SHIELD_COST_MINUS_RATIO":
      return FightPropType.FIGHT_PROP_NONEXTRA_SHIELD_COST_MINUS_RATIO;
    case 3024:
    case "FIGHT_PROP_NONEXTRA_PHYSICAL_ADD_HURT":
      return FightPropType.FIGHT_PROP_NONEXTRA_PHYSICAL_ADD_HURT;
    case 3045:
    case "FIGHT_PROP_BASE_ELEM_REACT_CRITICAL":
      return FightPropType.FIGHT_PROP_BASE_ELEM_REACT_CRITICAL;
    case 3046:
    case "FIGHT_PROP_BASE_ELEM_REACT_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_BASE_ELEM_REACT_CRITICAL_HURT;
    case 3025:
    case "FIGHT_PROP_ELEM_REACT_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_CRITICAL;
    case 3026:
    case "FIGHT_PROP_ELEM_REACT_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_CRITICAL_HURT;
    case 3027:
    case "FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL;
    case 3028:
    case "FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL_HURT;
    case 3029:
    case "FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL;
    case 3030:
    case "FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL_HURT;
    case 3031:
    case "FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL;
    case 3032:
    case "FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL_HURT;
    case 3033:
    case "FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL;
    case 3034:
    case "FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL_HURT;
    case 3035:
    case "FIGHT_PROP_ELEM_REACT_BURN_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_BURN_CRITICAL;
    case 3036:
    case "FIGHT_PROP_ELEM_REACT_BURN_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_BURN_CRITICAL_HURT;
    case 3037:
    case "FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL;
    case 3038:
    case "FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL_HURT;
    case 3039:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL;
    case 3040:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL_HURT;
    case 3041:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL;
    case 3042:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL_HURT;
    case 3043:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL;
    case 3044:
    case "FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL_HURT":
      return FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL_HURT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FightPropType.UNRECOGNIZED;
  }
}

export function fightPropTypeToJSON(object: FightPropType): string {
  switch (object) {
    case FightPropType.FIGHT_PROP_NONE:
      return "FIGHT_PROP_NONE";
    case FightPropType.FIGHT_PROP_BASE_HP:
      return "FIGHT_PROP_BASE_HP";
    case FightPropType.FIGHT_PROP_HP:
      return "FIGHT_PROP_HP";
    case FightPropType.FIGHT_PROP_HP_PERCENT:
      return "FIGHT_PROP_HP_PERCENT";
    case FightPropType.FIGHT_PROP_BASE_ATTACK:
      return "FIGHT_PROP_BASE_ATTACK";
    case FightPropType.FIGHT_PROP_ATTACK:
      return "FIGHT_PROP_ATTACK";
    case FightPropType.FIGHT_PROP_ATTACK_PERCENT:
      return "FIGHT_PROP_ATTACK_PERCENT";
    case FightPropType.FIGHT_PROP_BASE_DEFENSE:
      return "FIGHT_PROP_BASE_DEFENSE";
    case FightPropType.FIGHT_PROP_DEFENSE:
      return "FIGHT_PROP_DEFENSE";
    case FightPropType.FIGHT_PROP_DEFENSE_PERCENT:
      return "FIGHT_PROP_DEFENSE_PERCENT";
    case FightPropType.FIGHT_PROP_BASE_SPEED:
      return "FIGHT_PROP_BASE_SPEED";
    case FightPropType.FIGHT_PROP_SPEED_PERCENT:
      return "FIGHT_PROP_SPEED_PERCENT";
    case FightPropType.FIGHT_PROP_HP_MP_PERCENT:
      return "FIGHT_PROP_HP_MP_PERCENT";
    case FightPropType.FIGHT_PROP_ATTACK_MP_PERCENT:
      return "FIGHT_PROP_ATTACK_MP_PERCENT";
    case FightPropType.FIGHT_PROP_CRITICAL:
      return "FIGHT_PROP_CRITICAL";
    case FightPropType.FIGHT_PROP_ANTI_CRITICAL:
      return "FIGHT_PROP_ANTI_CRITICAL";
    case FightPropType.FIGHT_PROP_CRITICAL_HURT:
      return "FIGHT_PROP_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_CHARGE_EFFICIENCY:
      return "FIGHT_PROP_CHARGE_EFFICIENCY";
    case FightPropType.FIGHT_PROP_ADD_HURT:
      return "FIGHT_PROP_ADD_HURT";
    case FightPropType.FIGHT_PROP_SUB_HURT:
      return "FIGHT_PROP_SUB_HURT";
    case FightPropType.FIGHT_PROP_HEAL_ADD:
      return "FIGHT_PROP_HEAL_ADD";
    case FightPropType.FIGHT_PROP_HEALED_ADD:
      return "FIGHT_PROP_HEALED_ADD";
    case FightPropType.FIGHT_PROP_ELEMENT_MASTERY:
      return "FIGHT_PROP_ELEMENT_MASTERY";
    case FightPropType.FIGHT_PROP_PHYSICAL_SUB_HURT:
      return "FIGHT_PROP_PHYSICAL_SUB_HURT";
    case FightPropType.FIGHT_PROP_PHYSICAL_ADD_HURT:
      return "FIGHT_PROP_PHYSICAL_ADD_HURT";
    case FightPropType.FIGHT_PROP_DEFENCE_IGNORE_RATIO:
      return "FIGHT_PROP_DEFENCE_IGNORE_RATIO";
    case FightPropType.FIGHT_PROP_DEFENCE_IGNORE_DELTA:
      return "FIGHT_PROP_DEFENCE_IGNORE_DELTA";
    case FightPropType.FIGHT_PROP_FIRE_ADD_HURT:
      return "FIGHT_PROP_FIRE_ADD_HURT";
    case FightPropType.FIGHT_PROP_ELEC_ADD_HURT:
      return "FIGHT_PROP_ELEC_ADD_HURT";
    case FightPropType.FIGHT_PROP_WATER_ADD_HURT:
      return "FIGHT_PROP_WATER_ADD_HURT";
    case FightPropType.FIGHT_PROP_GRASS_ADD_HURT:
      return "FIGHT_PROP_GRASS_ADD_HURT";
    case FightPropType.FIGHT_PROP_WIND_ADD_HURT:
      return "FIGHT_PROP_WIND_ADD_HURT";
    case FightPropType.FIGHT_PROP_ROCK_ADD_HURT:
      return "FIGHT_PROP_ROCK_ADD_HURT";
    case FightPropType.FIGHT_PROP_ICE_ADD_HURT:
      return "FIGHT_PROP_ICE_ADD_HURT";
    case FightPropType.FIGHT_PROP_HIT_HEAD_ADD_HURT:
      return "FIGHT_PROP_HIT_HEAD_ADD_HURT";
    case FightPropType.FIGHT_PROP_FIRE_SUB_HURT:
      return "FIGHT_PROP_FIRE_SUB_HURT";
    case FightPropType.FIGHT_PROP_ELEC_SUB_HURT:
      return "FIGHT_PROP_ELEC_SUB_HURT";
    case FightPropType.FIGHT_PROP_WATER_SUB_HURT:
      return "FIGHT_PROP_WATER_SUB_HURT";
    case FightPropType.FIGHT_PROP_GRASS_SUB_HURT:
      return "FIGHT_PROP_GRASS_SUB_HURT";
    case FightPropType.FIGHT_PROP_WIND_SUB_HURT:
      return "FIGHT_PROP_WIND_SUB_HURT";
    case FightPropType.FIGHT_PROP_ROCK_SUB_HURT:
      return "FIGHT_PROP_ROCK_SUB_HURT";
    case FightPropType.FIGHT_PROP_ICE_SUB_HURT:
      return "FIGHT_PROP_ICE_SUB_HURT";
    case FightPropType.FIGHT_PROP_EFFECT_HIT:
      return "FIGHT_PROP_EFFECT_HIT";
    case FightPropType.FIGHT_PROP_EFFECT_RESIST:
      return "FIGHT_PROP_EFFECT_RESIST";
    case FightPropType.FIGHT_PROP_FREEZE_RESIST:
      return "FIGHT_PROP_FREEZE_RESIST";
    case FightPropType.FIGHT_PROP_DIZZY_RESIST:
      return "FIGHT_PROP_DIZZY_RESIST";
    case FightPropType.FIGHT_PROP_FREEZE_SHORTEN:
      return "FIGHT_PROP_FREEZE_SHORTEN";
    case FightPropType.FIGHT_PROP_DIZZY_SHORTEN:
      return "FIGHT_PROP_DIZZY_SHORTEN";
    case FightPropType.FIGHT_PROP_MAX_FIRE_ENERGY:
      return "FIGHT_PROP_MAX_FIRE_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_ELEC_ENERGY:
      return "FIGHT_PROP_MAX_ELEC_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_WATER_ENERGY:
      return "FIGHT_PROP_MAX_WATER_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_GRASS_ENERGY:
      return "FIGHT_PROP_MAX_GRASS_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_WIND_ENERGY:
      return "FIGHT_PROP_MAX_WIND_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_ICE_ENERGY:
      return "FIGHT_PROP_MAX_ICE_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_ROCK_ENERGY:
      return "FIGHT_PROP_MAX_ROCK_ENERGY";
    case FightPropType.FIGHT_PROP_MAX_SPECIAL_ENERGY:
      return "FIGHT_PROP_MAX_SPECIAL_ENERGY";
    case FightPropType.FIGHT_PROP_START_SPECIAL_ENERGY:
      return "FIGHT_PROP_START_SPECIAL_ENERGY";
    case FightPropType.FIGHT_PROP_SKILL_CD_MINUS_RATIO:
      return "FIGHT_PROP_SKILL_CD_MINUS_RATIO";
    case FightPropType.FIGHT_PROP_SHIELD_COST_MINUS_RATIO:
      return "FIGHT_PROP_SHIELD_COST_MINUS_RATIO";
    case FightPropType.FIGHT_PROP_CUR_FIRE_ENERGY:
      return "FIGHT_PROP_CUR_FIRE_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_ELEC_ENERGY:
      return "FIGHT_PROP_CUR_ELEC_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_WATER_ENERGY:
      return "FIGHT_PROP_CUR_WATER_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_GRASS_ENERGY:
      return "FIGHT_PROP_CUR_GRASS_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_WIND_ENERGY:
      return "FIGHT_PROP_CUR_WIND_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_ICE_ENERGY:
      return "FIGHT_PROP_CUR_ICE_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_ROCK_ENERGY:
      return "FIGHT_PROP_CUR_ROCK_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_SPECIAL_ENERGY:
      return "FIGHT_PROP_CUR_SPECIAL_ENERGY";
    case FightPropType.FIGHT_PROP_CUR_HP:
      return "FIGHT_PROP_CUR_HP";
    case FightPropType.FIGHT_PROP_MAX_HP:
      return "FIGHT_PROP_MAX_HP";
    case FightPropType.FIGHT_PROP_CUR_ATTACK:
      return "FIGHT_PROP_CUR_ATTACK";
    case FightPropType.FIGHT_PROP_CUR_DEFENSE:
      return "FIGHT_PROP_CUR_DEFENSE";
    case FightPropType.FIGHT_PROP_CUR_SPEED:
      return "FIGHT_PROP_CUR_SPEED";
    case FightPropType.FIGHT_PROP_CUR_HP_DEBTS:
      return "FIGHT_PROP_CUR_HP_DEBTS";
    case FightPropType.FIGHT_PROP_CUR_HP_PAID_DEBTS:
      return "FIGHT_PROP_CUR_HP_PAID_DEBTS";
    case FightPropType.FIGHT_PROP_CUR_NATLAN_HP:
      return "FIGHT_PROP_CUR_NATLAN_HP";
    case FightPropType.FIGHT_PROP_NONEXTRA_ATTACK:
      return "FIGHT_PROP_NONEXTRA_ATTACK";
    case FightPropType.FIGHT_PROP_NONEXTRA_DEFENSE:
      return "FIGHT_PROP_NONEXTRA_DEFENSE";
    case FightPropType.FIGHT_PROP_NONEXTRA_CRITICAL:
      return "FIGHT_PROP_NONEXTRA_CRITICAL";
    case FightPropType.FIGHT_PROP_NONEXTRA_ANTI_CRITICAL:
      return "FIGHT_PROP_NONEXTRA_ANTI_CRITICAL";
    case FightPropType.FIGHT_PROP_NONEXTRA_CRITICAL_HURT:
      return "FIGHT_PROP_NONEXTRA_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_CHARGE_EFFICIENCY:
      return "FIGHT_PROP_NONEXTRA_CHARGE_EFFICIENCY";
    case FightPropType.FIGHT_PROP_NONEXTRA_ELEMENT_MASTERY:
      return "FIGHT_PROP_NONEXTRA_ELEMENT_MASTERY";
    case FightPropType.FIGHT_PROP_NONEXTRA_PHYSICAL_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_PHYSICAL_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_FIRE_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_FIRE_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ELEC_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_ELEC_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_WATER_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_WATER_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_GRASS_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_GRASS_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_WIND_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_WIND_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ROCK_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_ROCK_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ICE_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_ICE_ADD_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_FIRE_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_FIRE_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ELEC_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_ELEC_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_WATER_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_WATER_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_GRASS_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_GRASS_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_WIND_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_WIND_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ROCK_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_ROCK_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_ICE_SUB_HURT:
      return "FIGHT_PROP_NONEXTRA_ICE_SUB_HURT";
    case FightPropType.FIGHT_PROP_NONEXTRA_SKILL_CD_MINUS_RATIO:
      return "FIGHT_PROP_NONEXTRA_SKILL_CD_MINUS_RATIO";
    case FightPropType.FIGHT_PROP_NONEXTRA_SHIELD_COST_MINUS_RATIO:
      return "FIGHT_PROP_NONEXTRA_SHIELD_COST_MINUS_RATIO";
    case FightPropType.FIGHT_PROP_NONEXTRA_PHYSICAL_ADD_HURT:
      return "FIGHT_PROP_NONEXTRA_PHYSICAL_ADD_HURT";
    case FightPropType.FIGHT_PROP_BASE_ELEM_REACT_CRITICAL:
      return "FIGHT_PROP_BASE_ELEM_REACT_CRITICAL";
    case FightPropType.FIGHT_PROP_BASE_ELEM_REACT_CRITICAL_HURT:
      return "FIGHT_PROP_BASE_ELEM_REACT_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_EXPLODE_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_SWIRL_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_ELECTRIC_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_SCONDUCT_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_BURN_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_BURN_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_BURN_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_BURN_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_FROZENBROKEN_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_FIRE_CRITICAL_HURT";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL";
    case FightPropType.FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL_HURT:
      return "FIGHT_PROP_ELEM_REACT_OVERGROW_ELECTRIC_CRITICAL_HURT";
    case FightPropType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum GrowCurveType {
  GROW_CURVE_NONE = 0,
  GROW_CURVE_HP = 1,
  GROW_CURVE_ATTACK = 2,
  GROW_CURVE_STAMINA = 3,
  GROW_CURVE_STRIKE = 4,
  GROW_CURVE_ANTI_STRIKE = 5,
  GROW_CURVE_ANTI_STRIKE1 = 6,
  GROW_CURVE_ANTI_STRIKE2 = 7,
  GROW_CURVE_ANTI_STRIKE3 = 8,
  GROW_CURVE_STRIKE_HURT = 9,
  GROW_CURVE_ELEMENT = 10,
  GROW_CURVE_KILL_EXP = 11,
  GROW_CURVE_DEFENSE = 12,
  GROW_CURVE_ATTACK_BOMB = 13,
  GROW_CURVE_HP_LITTLEMONSTER = 14,
  GROW_CURVE_ELEMENT_MASTERY = 15,
  GROW_CURVE_PROGRESSION = 16,
  GROW_CURVE_DEFENDING = 17,
  GROW_CURVE_MHP = 18,
  GROW_CURVE_MATK = 19,
  GROW_CURVE_TOWERATK = 20,
  GROW_CURVE_HP_S5 = 21,
  GROW_CURVE_HP_S4 = 22,
  GROW_CURVE_HP_2 = 23,
  GROW_CURVE_ATTACK_2 = 24,
  GROW_CURVE_HP_ENVIRONMENT = 25,
  GROW_CURVE_ATTACK_S5 = 31,
  GROW_CURVE_ATTACK_S4 = 32,
  GROW_CURVE_ATTACK_S3 = 33,
  GROW_CURVE_STRIKE_S5 = 34,
  GROW_CURVE_DEFENSE_S5 = 41,
  GROW_CURVE_DEFENSE_S4 = 42,
  GROW_CURVE_ATTACK_101 = 1101,
  GROW_CURVE_ATTACK_102 = 1102,
  GROW_CURVE_ATTACK_103 = 1103,
  GROW_CURVE_ATTACK_104 = 1104,
  GROW_CURVE_ATTACK_105 = 1105,
  GROW_CURVE_ATTACK_201 = 1201,
  GROW_CURVE_ATTACK_202 = 1202,
  GROW_CURVE_ATTACK_203 = 1203,
  GROW_CURVE_ATTACK_204 = 1204,
  GROW_CURVE_ATTACK_205 = 1205,
  GROW_CURVE_ATTACK_301 = 1301,
  GROW_CURVE_ATTACK_302 = 1302,
  GROW_CURVE_ATTACK_303 = 1303,
  GROW_CURVE_ATTACK_304 = 1304,
  GROW_CURVE_ATTACK_305 = 1305,
  GROW_CURVE_CRITICAL_101 = 2101,
  GROW_CURVE_CRITICAL_102 = 2102,
  GROW_CURVE_CRITICAL_103 = 2103,
  GROW_CURVE_CRITICAL_104 = 2104,
  GROW_CURVE_CRITICAL_105 = 2105,
  GROW_CURVE_CRITICAL_201 = 2201,
  GROW_CURVE_CRITICAL_202 = 2202,
  GROW_CURVE_CRITICAL_203 = 2203,
  GROW_CURVE_CRITICAL_204 = 2204,
  GROW_CURVE_CRITICAL_205 = 2205,
  GROW_CURVE_CRITICAL_301 = 2301,
  GROW_CURVE_CRITICAL_302 = 2302,
  GROW_CURVE_CRITICAL_303 = 2303,
  GROW_CURVE_CRITICAL_304 = 2304,
  GROW_CURVE_CRITICAL_305 = 2305,
  GROW_CURVE_ACTIVITY_ATTACK_1 = 5201,
  GROW_CURVE_ACTIVITY_HP_1 = 5202,
  GROW_CURVE_ACTIVITY_ATTACK_2 = 5701,
  GROW_CURVE_ACTIVITY_HP_2 = 5702,
  GROW_CURVE_ACTIVITY_DEFENSE_2 = 5703,
  UNRECOGNIZED = -1,
}

export function growCurveTypeFromJSON(object: any): GrowCurveType {
  switch (object) {
    case 0:
    case "GROW_CURVE_NONE":
      return GrowCurveType.GROW_CURVE_NONE;
    case 1:
    case "GROW_CURVE_HP":
      return GrowCurveType.GROW_CURVE_HP;
    case 2:
    case "GROW_CURVE_ATTACK":
      return GrowCurveType.GROW_CURVE_ATTACK;
    case 3:
    case "GROW_CURVE_STAMINA":
      return GrowCurveType.GROW_CURVE_STAMINA;
    case 4:
    case "GROW_CURVE_STRIKE":
      return GrowCurveType.GROW_CURVE_STRIKE;
    case 5:
    case "GROW_CURVE_ANTI_STRIKE":
      return GrowCurveType.GROW_CURVE_ANTI_STRIKE;
    case 6:
    case "GROW_CURVE_ANTI_STRIKE1":
      return GrowCurveType.GROW_CURVE_ANTI_STRIKE1;
    case 7:
    case "GROW_CURVE_ANTI_STRIKE2":
      return GrowCurveType.GROW_CURVE_ANTI_STRIKE2;
    case 8:
    case "GROW_CURVE_ANTI_STRIKE3":
      return GrowCurveType.GROW_CURVE_ANTI_STRIKE3;
    case 9:
    case "GROW_CURVE_STRIKE_HURT":
      return GrowCurveType.GROW_CURVE_STRIKE_HURT;
    case 10:
    case "GROW_CURVE_ELEMENT":
      return GrowCurveType.GROW_CURVE_ELEMENT;
    case 11:
    case "GROW_CURVE_KILL_EXP":
      return GrowCurveType.GROW_CURVE_KILL_EXP;
    case 12:
    case "GROW_CURVE_DEFENSE":
      return GrowCurveType.GROW_CURVE_DEFENSE;
    case 13:
    case "GROW_CURVE_ATTACK_BOMB":
      return GrowCurveType.GROW_CURVE_ATTACK_BOMB;
    case 14:
    case "GROW_CURVE_HP_LITTLEMONSTER":
      return GrowCurveType.GROW_CURVE_HP_LITTLEMONSTER;
    case 15:
    case "GROW_CURVE_ELEMENT_MASTERY":
      return GrowCurveType.GROW_CURVE_ELEMENT_MASTERY;
    case 16:
    case "GROW_CURVE_PROGRESSION":
      return GrowCurveType.GROW_CURVE_PROGRESSION;
    case 17:
    case "GROW_CURVE_DEFENDING":
      return GrowCurveType.GROW_CURVE_DEFENDING;
    case 18:
    case "GROW_CURVE_MHP":
      return GrowCurveType.GROW_CURVE_MHP;
    case 19:
    case "GROW_CURVE_MATK":
      return GrowCurveType.GROW_CURVE_MATK;
    case 20:
    case "GROW_CURVE_TOWERATK":
      return GrowCurveType.GROW_CURVE_TOWERATK;
    case 21:
    case "GROW_CURVE_HP_S5":
      return GrowCurveType.GROW_CURVE_HP_S5;
    case 22:
    case "GROW_CURVE_HP_S4":
      return GrowCurveType.GROW_CURVE_HP_S4;
    case 23:
    case "GROW_CURVE_HP_2":
      return GrowCurveType.GROW_CURVE_HP_2;
    case 24:
    case "GROW_CURVE_ATTACK_2":
      return GrowCurveType.GROW_CURVE_ATTACK_2;
    case 25:
    case "GROW_CURVE_HP_ENVIRONMENT":
      return GrowCurveType.GROW_CURVE_HP_ENVIRONMENT;
    case 31:
    case "GROW_CURVE_ATTACK_S5":
      return GrowCurveType.GROW_CURVE_ATTACK_S5;
    case 32:
    case "GROW_CURVE_ATTACK_S4":
      return GrowCurveType.GROW_CURVE_ATTACK_S4;
    case 33:
    case "GROW_CURVE_ATTACK_S3":
      return GrowCurveType.GROW_CURVE_ATTACK_S3;
    case 34:
    case "GROW_CURVE_STRIKE_S5":
      return GrowCurveType.GROW_CURVE_STRIKE_S5;
    case 41:
    case "GROW_CURVE_DEFENSE_S5":
      return GrowCurveType.GROW_CURVE_DEFENSE_S5;
    case 42:
    case "GROW_CURVE_DEFENSE_S4":
      return GrowCurveType.GROW_CURVE_DEFENSE_S4;
    case 1101:
    case "GROW_CURVE_ATTACK_101":
      return GrowCurveType.GROW_CURVE_ATTACK_101;
    case 1102:
    case "GROW_CURVE_ATTACK_102":
      return GrowCurveType.GROW_CURVE_ATTACK_102;
    case 1103:
    case "GROW_CURVE_ATTACK_103":
      return GrowCurveType.GROW_CURVE_ATTACK_103;
    case 1104:
    case "GROW_CURVE_ATTACK_104":
      return GrowCurveType.GROW_CURVE_ATTACK_104;
    case 1105:
    case "GROW_CURVE_ATTACK_105":
      return GrowCurveType.GROW_CURVE_ATTACK_105;
    case 1201:
    case "GROW_CURVE_ATTACK_201":
      return GrowCurveType.GROW_CURVE_ATTACK_201;
    case 1202:
    case "GROW_CURVE_ATTACK_202":
      return GrowCurveType.GROW_CURVE_ATTACK_202;
    case 1203:
    case "GROW_CURVE_ATTACK_203":
      return GrowCurveType.GROW_CURVE_ATTACK_203;
    case 1204:
    case "GROW_CURVE_ATTACK_204":
      return GrowCurveType.GROW_CURVE_ATTACK_204;
    case 1205:
    case "GROW_CURVE_ATTACK_205":
      return GrowCurveType.GROW_CURVE_ATTACK_205;
    case 1301:
    case "GROW_CURVE_ATTACK_301":
      return GrowCurveType.GROW_CURVE_ATTACK_301;
    case 1302:
    case "GROW_CURVE_ATTACK_302":
      return GrowCurveType.GROW_CURVE_ATTACK_302;
    case 1303:
    case "GROW_CURVE_ATTACK_303":
      return GrowCurveType.GROW_CURVE_ATTACK_303;
    case 1304:
    case "GROW_CURVE_ATTACK_304":
      return GrowCurveType.GROW_CURVE_ATTACK_304;
    case 1305:
    case "GROW_CURVE_ATTACK_305":
      return GrowCurveType.GROW_CURVE_ATTACK_305;
    case 2101:
    case "GROW_CURVE_CRITICAL_101":
      return GrowCurveType.GROW_CURVE_CRITICAL_101;
    case 2102:
    case "GROW_CURVE_CRITICAL_102":
      return GrowCurveType.GROW_CURVE_CRITICAL_102;
    case 2103:
    case "GROW_CURVE_CRITICAL_103":
      return GrowCurveType.GROW_CURVE_CRITICAL_103;
    case 2104:
    case "GROW_CURVE_CRITICAL_104":
      return GrowCurveType.GROW_CURVE_CRITICAL_104;
    case 2105:
    case "GROW_CURVE_CRITICAL_105":
      return GrowCurveType.GROW_CURVE_CRITICAL_105;
    case 2201:
    case "GROW_CURVE_CRITICAL_201":
      return GrowCurveType.GROW_CURVE_CRITICAL_201;
    case 2202:
    case "GROW_CURVE_CRITICAL_202":
      return GrowCurveType.GROW_CURVE_CRITICAL_202;
    case 2203:
    case "GROW_CURVE_CRITICAL_203":
      return GrowCurveType.GROW_CURVE_CRITICAL_203;
    case 2204:
    case "GROW_CURVE_CRITICAL_204":
      return GrowCurveType.GROW_CURVE_CRITICAL_204;
    case 2205:
    case "GROW_CURVE_CRITICAL_205":
      return GrowCurveType.GROW_CURVE_CRITICAL_205;
    case 2301:
    case "GROW_CURVE_CRITICAL_301":
      return GrowCurveType.GROW_CURVE_CRITICAL_301;
    case 2302:
    case "GROW_CURVE_CRITICAL_302":
      return GrowCurveType.GROW_CURVE_CRITICAL_302;
    case 2303:
    case "GROW_CURVE_CRITICAL_303":
      return GrowCurveType.GROW_CURVE_CRITICAL_303;
    case 2304:
    case "GROW_CURVE_CRITICAL_304":
      return GrowCurveType.GROW_CURVE_CRITICAL_304;
    case 2305:
    case "GROW_CURVE_CRITICAL_305":
      return GrowCurveType.GROW_CURVE_CRITICAL_305;
    case 5201:
    case "GROW_CURVE_ACTIVITY_ATTACK_1":
      return GrowCurveType.GROW_CURVE_ACTIVITY_ATTACK_1;
    case 5202:
    case "GROW_CURVE_ACTIVITY_HP_1":
      return GrowCurveType.GROW_CURVE_ACTIVITY_HP_1;
    case 5701:
    case "GROW_CURVE_ACTIVITY_ATTACK_2":
      return GrowCurveType.GROW_CURVE_ACTIVITY_ATTACK_2;
    case 5702:
    case "GROW_CURVE_ACTIVITY_HP_2":
      return GrowCurveType.GROW_CURVE_ACTIVITY_HP_2;
    case 5703:
    case "GROW_CURVE_ACTIVITY_DEFENSE_2":
      return GrowCurveType.GROW_CURVE_ACTIVITY_DEFENSE_2;
    case -1:
    case "UNRECOGNIZED":
    default:
      return GrowCurveType.UNRECOGNIZED;
  }
}

export function growCurveTypeToJSON(object: GrowCurveType): string {
  switch (object) {
    case GrowCurveType.GROW_CURVE_NONE:
      return "GROW_CURVE_NONE";
    case GrowCurveType.GROW_CURVE_HP:
      return "GROW_CURVE_HP";
    case GrowCurveType.GROW_CURVE_ATTACK:
      return "GROW_CURVE_ATTACK";
    case GrowCurveType.GROW_CURVE_STAMINA:
      return "GROW_CURVE_STAMINA";
    case GrowCurveType.GROW_CURVE_STRIKE:
      return "GROW_CURVE_STRIKE";
    case GrowCurveType.GROW_CURVE_ANTI_STRIKE:
      return "GROW_CURVE_ANTI_STRIKE";
    case GrowCurveType.GROW_CURVE_ANTI_STRIKE1:
      return "GROW_CURVE_ANTI_STRIKE1";
    case GrowCurveType.GROW_CURVE_ANTI_STRIKE2:
      return "GROW_CURVE_ANTI_STRIKE2";
    case GrowCurveType.GROW_CURVE_ANTI_STRIKE3:
      return "GROW_CURVE_ANTI_STRIKE3";
    case GrowCurveType.GROW_CURVE_STRIKE_HURT:
      return "GROW_CURVE_STRIKE_HURT";
    case GrowCurveType.GROW_CURVE_ELEMENT:
      return "GROW_CURVE_ELEMENT";
    case GrowCurveType.GROW_CURVE_KILL_EXP:
      return "GROW_CURVE_KILL_EXP";
    case GrowCurveType.GROW_CURVE_DEFENSE:
      return "GROW_CURVE_DEFENSE";
    case GrowCurveType.GROW_CURVE_ATTACK_BOMB:
      return "GROW_CURVE_ATTACK_BOMB";
    case GrowCurveType.GROW_CURVE_HP_LITTLEMONSTER:
      return "GROW_CURVE_HP_LITTLEMONSTER";
    case GrowCurveType.GROW_CURVE_ELEMENT_MASTERY:
      return "GROW_CURVE_ELEMENT_MASTERY";
    case GrowCurveType.GROW_CURVE_PROGRESSION:
      return "GROW_CURVE_PROGRESSION";
    case GrowCurveType.GROW_CURVE_DEFENDING:
      return "GROW_CURVE_DEFENDING";
    case GrowCurveType.GROW_CURVE_MHP:
      return "GROW_CURVE_MHP";
    case GrowCurveType.GROW_CURVE_MATK:
      return "GROW_CURVE_MATK";
    case GrowCurveType.GROW_CURVE_TOWERATK:
      return "GROW_CURVE_TOWERATK";
    case GrowCurveType.GROW_CURVE_HP_S5:
      return "GROW_CURVE_HP_S5";
    case GrowCurveType.GROW_CURVE_HP_S4:
      return "GROW_CURVE_HP_S4";
    case GrowCurveType.GROW_CURVE_HP_2:
      return "GROW_CURVE_HP_2";
    case GrowCurveType.GROW_CURVE_ATTACK_2:
      return "GROW_CURVE_ATTACK_2";
    case GrowCurveType.GROW_CURVE_HP_ENVIRONMENT:
      return "GROW_CURVE_HP_ENVIRONMENT";
    case GrowCurveType.GROW_CURVE_ATTACK_S5:
      return "GROW_CURVE_ATTACK_S5";
    case GrowCurveType.GROW_CURVE_ATTACK_S4:
      return "GROW_CURVE_ATTACK_S4";
    case GrowCurveType.GROW_CURVE_ATTACK_S3:
      return "GROW_CURVE_ATTACK_S3";
    case GrowCurveType.GROW_CURVE_STRIKE_S5:
      return "GROW_CURVE_STRIKE_S5";
    case GrowCurveType.GROW_CURVE_DEFENSE_S5:
      return "GROW_CURVE_DEFENSE_S5";
    case GrowCurveType.GROW_CURVE_DEFENSE_S4:
      return "GROW_CURVE_DEFENSE_S4";
    case GrowCurveType.GROW_CURVE_ATTACK_101:
      return "GROW_CURVE_ATTACK_101";
    case GrowCurveType.GROW_CURVE_ATTACK_102:
      return "GROW_CURVE_ATTACK_102";
    case GrowCurveType.GROW_CURVE_ATTACK_103:
      return "GROW_CURVE_ATTACK_103";
    case GrowCurveType.GROW_CURVE_ATTACK_104:
      return "GROW_CURVE_ATTACK_104";
    case GrowCurveType.GROW_CURVE_ATTACK_105:
      return "GROW_CURVE_ATTACK_105";
    case GrowCurveType.GROW_CURVE_ATTACK_201:
      return "GROW_CURVE_ATTACK_201";
    case GrowCurveType.GROW_CURVE_ATTACK_202:
      return "GROW_CURVE_ATTACK_202";
    case GrowCurveType.GROW_CURVE_ATTACK_203:
      return "GROW_CURVE_ATTACK_203";
    case GrowCurveType.GROW_CURVE_ATTACK_204:
      return "GROW_CURVE_ATTACK_204";
    case GrowCurveType.GROW_CURVE_ATTACK_205:
      return "GROW_CURVE_ATTACK_205";
    case GrowCurveType.GROW_CURVE_ATTACK_301:
      return "GROW_CURVE_ATTACK_301";
    case GrowCurveType.GROW_CURVE_ATTACK_302:
      return "GROW_CURVE_ATTACK_302";
    case GrowCurveType.GROW_CURVE_ATTACK_303:
      return "GROW_CURVE_ATTACK_303";
    case GrowCurveType.GROW_CURVE_ATTACK_304:
      return "GROW_CURVE_ATTACK_304";
    case GrowCurveType.GROW_CURVE_ATTACK_305:
      return "GROW_CURVE_ATTACK_305";
    case GrowCurveType.GROW_CURVE_CRITICAL_101:
      return "GROW_CURVE_CRITICAL_101";
    case GrowCurveType.GROW_CURVE_CRITICAL_102:
      return "GROW_CURVE_CRITICAL_102";
    case GrowCurveType.GROW_CURVE_CRITICAL_103:
      return "GROW_CURVE_CRITICAL_103";
    case GrowCurveType.GROW_CURVE_CRITICAL_104:
      return "GROW_CURVE_CRITICAL_104";
    case GrowCurveType.GROW_CURVE_CRITICAL_105:
      return "GROW_CURVE_CRITICAL_105";
    case GrowCurveType.GROW_CURVE_CRITICAL_201:
      return "GROW_CURVE_CRITICAL_201";
    case GrowCurveType.GROW_CURVE_CRITICAL_202:
      return "GROW_CURVE_CRITICAL_202";
    case GrowCurveType.GROW_CURVE_CRITICAL_203:
      return "GROW_CURVE_CRITICAL_203";
    case GrowCurveType.GROW_CURVE_CRITICAL_204:
      return "GROW_CURVE_CRITICAL_204";
    case GrowCurveType.GROW_CURVE_CRITICAL_205:
      return "GROW_CURVE_CRITICAL_205";
    case GrowCurveType.GROW_CURVE_CRITICAL_301:
      return "GROW_CURVE_CRITICAL_301";
    case GrowCurveType.GROW_CURVE_CRITICAL_302:
      return "GROW_CURVE_CRITICAL_302";
    case GrowCurveType.GROW_CURVE_CRITICAL_303:
      return "GROW_CURVE_CRITICAL_303";
    case GrowCurveType.GROW_CURVE_CRITICAL_304:
      return "GROW_CURVE_CRITICAL_304";
    case GrowCurveType.GROW_CURVE_CRITICAL_305:
      return "GROW_CURVE_CRITICAL_305";
    case GrowCurveType.GROW_CURVE_ACTIVITY_ATTACK_1:
      return "GROW_CURVE_ACTIVITY_ATTACK_1";
    case GrowCurveType.GROW_CURVE_ACTIVITY_HP_1:
      return "GROW_CURVE_ACTIVITY_HP_1";
    case GrowCurveType.GROW_CURVE_ACTIVITY_ATTACK_2:
      return "GROW_CURVE_ACTIVITY_ATTACK_2";
    case GrowCurveType.GROW_CURVE_ACTIVITY_HP_2:
      return "GROW_CURVE_ACTIVITY_HP_2";
    case GrowCurveType.GROW_CURVE_ACTIVITY_DEFENSE_2:
      return "GROW_CURVE_ACTIVITY_DEFENSE_2";
    case GrowCurveType.UNRECOGNIZED:
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
