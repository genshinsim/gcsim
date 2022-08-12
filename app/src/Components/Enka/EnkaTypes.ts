// https://github.com/EnkaNetwork/API-docs/blob/master/api.md
export interface EnkaData {
  playerInfo: object;
  avatarInfoList: AvatarInfo[];
}

export interface AvatarInfo {
  //Name
  avatarId: number;

  //Constellation id
  talentIdList?: number[];

  //There are other objects in this field, but I don't know what they are for.
  //Denoting the useful ones here for now.
  propMap: {
    //Ascension
    1002: {
      type: 1002;
      ival: string;
      val: string;
    };
    //Level
    4001: {
      type: 4001;
      ival: string;
      val: string;
    };
  };

  //Total Stat Obj, can't use it as it combines base and artifact stats #blamesrl
  //   fightPropMap: any;

  //???
  skillDepotId: number;

  //Talents unlocked
  inherentProudSkillList: number[];

  //Talent Levels
  skillLevelMap: {
    [key: number]: number;
  };

  //TODO: Figure out what this is
  equipList: (GenshinItemWeapon | GenshinItemReliquary)[];
}

export interface GenshinItemWeapon {
  itemId: number;
  weapon: {
    level: number;
    //ascension
    promoteLevel: number;
    affixMap: {
      [key: number]: number;
    };
  };
  flat: {
    nameTextMapHash: string;
    rankLevel: number;
    weaponStats: [
      {
        // If this isn't attack, I've probably quit by then
        appendPropId: "FIGHT_PROP_BASE_ATTACK";
        statValue: number;
      },
      {
        appendPropId: FightProp;
        statValue: number;
      }
    ];
    itemType: "ITEM_WEAPON";
    icon: string;
  };
}

export interface GenshinItemReliquary {
  itemId: number;
  reliquary: {
    // ingame level + 1 (idk why)
    level: number;
    mainPropId: number;
    appendPropIdList: number[];
  };
  flat: {
    nameTextMapHash: string;
    // use this for setKey
    setNameTextMapHash: string;
    // rarity
    rankLevel: number;
    reliquaryMainstat: {
      mainPropId: FightProp;
      statValue: number;
    };
    reliquarySubstats: [
      {
        appendPropId: FightProp;
        statValue: number;
      },
      {
        appendPropId: FightProp;
        statValue: number;
      },
      {
        appendPropId: FightProp;
        statValue: number;
      },
      {
        appendPropId: FightProp;
        statValue: number;
      }
    ];
    itemType: "ITEM_RELIQUARY";
    icon: string;
    equipType: ReliquaryEquipType;
  };
}

export enum FightProp {
  FIGHT_PROP_HP = "FIGHT_PROP_HP",
  FIGHT_PROP_HP_PERCENT = "FIGHT_PROP_HP_PERCENT",
  FIGHT_PROP_ATTACK = "FIGHT_PROP_ATTACK",
  FIGHT_PROP_ATTACK_PERCENT = "FIGHT_PROP_ATTACK_PERCENT",
  FIGHT_PROP_DEFENSE = "FIGHT_PROP_DEFENSE",
  FIGHT_PROP_DEFENSE_PERCENT = "FIGHT_PROP_DEFENSE_PERCENT",
  FIGHT_PROP_ELEMENT_MASTERY = "FIGHT_PROP_ELEMENT_MASTERY",
  FIGHT_PROP_CRITICAL = "FIGHT_PROP_CRITICAL",
  FIGHT_PROP_CRITICAL_HURT = "FIGHT_PROP_CRITICAL_HURT",
  FIGHT_PROP_HEAL_ADD = "FIGHT_PROP_HEAL_ADD",
  FIGHT_PROP_HEALED_ADD = "FIGHT_PROP_HEALED_ADD",
  FIGHT_PROP_CHARGE_EFFICIENCY = "FIGHT_PROP_CHARGE_EFFICIENCY",
  FIGHT_PROP_SHIELD_COST_MINUS_RATIO = "FIGHT_PROP_SHIELD_COST_MINUS_RATIO",
  FIGHT_PROP_FIRE_ADD_HURT = "FIGHT_PROP_FIRE_ADD_HURT",
  FIGHT_PROP_WATER_ADD_HURT = "FIGHT_PROP_WATER_ADD_HURT",
  FIGHT_PROP_GRASS_ADD_HURT = "FIGHT_PROP_GRASS_ADD_HURT",
  FIGHT_PROP_ELEC_ADD_HURT = "FIGHT_PROP_ELEC_ADD_HURT",
  FIGHT_PROP_WIND_ADD_HURT = "FIGHT_PROP_WIND_ADD_HURT",
  FIGHT_PROP_ICE_ADD_HURT = "FIGHT_PROP_ICE_ADD_HURT",
  FIGHT_PROP_ROCK_ADD_HURT = "FIGHT_PROP_ROCK_ADD_HURT",
  FIGHT_PROP_PHYSICAL_ADD_HURT = "FIGHT_PROP_PHYSICAL_ADD_HURT",
}

export enum ReliquaryEquipType {
  EQUIP_BRACER = "EQUIP_BRACER",
  EQUIP_NECKLACE = "EQUIP_NECKLACE",
  EQUIP_SHOES = "EQUIP_SHOES",
  EQUIP_RING = "EQUIP_RING",
  EQUIP_DRESS = "EQUIP_DRESS",
}
