import TextMap from "./GenshinData/EnkaTextMapEN.json";
import CharacterMap from "./GenshinData/EnkaCharacterMap.json";

import {
  GOODArtifact,
  GOODArtifactSetKey,
  GOODCharacter,
  GOODCharacterKey,
  GOODSlotKey,
  GOODStatKey,
  GOODWeapon,
  GOODWeaponKey,
  IGOOD,
  ISubstat,
} from "../GOOD/GOODTypes";
import {
  EnkaData,
  FightProp,
  GenshinItemReliquary,
  GenshinItemWeapon,
  ReliquaryEquipType,
} from "./EnkaTypes";

export default function EnkaToGOOD(enkaData: EnkaData): IGOOD {
  let characters: GOODCharacter[] = [];
  let artifacts: GOODArtifact[] = [];
  let weapons: GOODWeapon[] = [];

  enkaData.avatarInfoList.forEach(
    ({ avatarId, propMap, talentIdList, skillLevelMap, equipList }) => {
      const character: GOODCharacter = {
        key: getGOODKeyFromAvatarId(avatarId),
        level: parseInt(propMap["4001"].val),
        ascension: parseInt(propMap["1002"].val),
        constellation: talentIdList?.length || 0,
        //Characters with 7 talents like AYAKA might be bugged (kokomi is fine?)//
        // characters with unodered talents like traveler(i htink) will not work
        talent: {
          auto: Object.entries(skillLevelMap)[0][1] as number,
          skill: Object.entries(skillLevelMap)[1][1] as number,
          burst: Object.entries(skillLevelMap)[2][1] as number,
        },
      };
      characters.push(character);

      equipList.forEach((equip) => {
        if (equip.flat.itemType == "ITEM_WEAPON") {
          const enkaWeapon = equip as GenshinItemWeapon;

          const weapon: GOODWeapon = {
            key: getGOODKeyFromWeaponNameTextMapHash(
              enkaWeapon.flat.nameTextMapHash
            ),
            level: enkaWeapon.weapon.level,
            ascension: enkaWeapon.weapon.promoteLevel,
            refinement:
              (Object.entries(enkaWeapon.weapon.affixMap)[0] != null
                ? Object.entries(enkaWeapon.weapon.affixMap)[0][1]
                : 0) + 1,
            location: character.key,
            lock: false,
          };
          weapons.push(weapon);
        } else {
          const enkaReliquary = equip as GenshinItemReliquary;
          const reliquary: GOODArtifact = {
            setKey: getGOODKeyFromReliquaryNameTextMapHash(
              enkaReliquary.flat.setNameTextMapHash
            ),
            level: enkaReliquary.reliquary.level - 1,
            slotKey: reliquaryTypeToGOODKey(enkaReliquary.flat.equipType),
            rarity: enkaReliquary.flat.rankLevel,
            location: character.key,
            lock: false,
            mainStatKey: fightPropToGOODKey(
              enkaReliquary.flat.reliquaryMainstat.mainPropId
            ),
            substats: getGOODSubstatsFromReliquarySubstats(
              enkaReliquary.flat.reliquarySubstats
            ),
          };
          artifacts.push(reliquary);
        }
      });
    }
  );

  return {
    format: "GOOD" as IGOOD["format"],
    version: 1,
    source: "gcsimFromEnka",
    characters,
    weapons,
    artifacts,
  };
}

function textToGOODKey(string: string) {
  function toTitleCase(str: string) {
    return str.replace(/-/g, " ").replace(/\w\S*/g, function (txt: string) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    });
  }
  return toTitleCase(string || "").replace(/[^A-Za-z]/g, "") as
    | GOODCharacterKey
    | GOODWeaponKey
    | GOODArtifactSetKey;
}
function reliquaryTypeToGOODKey(
  reliquaryType: ReliquaryEquipType
): GOODSlotKey {
  switch (reliquaryType) {
    case ReliquaryEquipType.EQUIP_BRACER:
      return "flower";
    case ReliquaryEquipType.EQUIP_NECKLACE:
      return "plume";
    case ReliquaryEquipType.EQUIP_SHOES:
      return "sands";
    case ReliquaryEquipType.EQUIP_RING:
      return "goblet";
    case ReliquaryEquipType.EQUIP_DRESS:
      return "circlet";
  }
}

function fightPropToGOODKey(fightProp: FightProp): GOODStatKey {
  switch (fightProp) {
    case FightProp.FIGHT_PROP_HP:
      return "hp";
    case FightProp.FIGHT_PROP_HP_PERCENT:
      return "hp_";
    case FightProp.FIGHT_PROP_ATTACK:
      return "atk";
    case FightProp.FIGHT_PROP_ATTACK_PERCENT:
      return "atk_";
    case FightProp.FIGHT_PROP_DEFENSE:
      return "def";
    case FightProp.FIGHT_PROP_DEFENSE_PERCENT:
      return "def_";
    case FightProp.FIGHT_PROP_CHARGE_EFFICIENCY:
      return "enerRech_";
    case FightProp.FIGHT_PROP_ELEMENT_MASTERY:
      return "eleMas";
    case FightProp.FIGHT_PROP_CRITICAL:
      return "critRate_";
    case FightProp.FIGHT_PROP_CRITICAL_HURT:
      return "critDMG_";
    case FightProp.FIGHT_PROP_HEAL_ADD:
      return "heal_";
    case FightProp.FIGHT_PROP_FIRE_ADD_HURT:
      return "pyro_dmg_";
    case FightProp.FIGHT_PROP_ELEC_ADD_HURT:
      return "electro_dmg_";
    case FightProp.FIGHT_PROP_ICE_ADD_HURT:
      return "cryo_dmg_";
    case FightProp.FIGHT_PROP_WATER_ADD_HURT:
      return "hydro_dmg_";
    case FightProp.FIGHT_PROP_WIND_ADD_HURT:
      return "anemo_dmg_";
    case FightProp.FIGHT_PROP_ROCK_ADD_HURT:
      return "geo_dmg_";
    //soon tm
    //   case FightProp.FIGHT_PROP_GRASS_ADD_HURT:
    //     return "dendro_dmg_";
    case FightProp.FIGHT_PROP_PHYSICAL_ADD_HURT:
      return "physical_dmg_";
    default:
      return "";
  }
}

const characterMap: ICharacterMap = CharacterMap;
const textMap: IENTextMap = TextMap.en;

function getGOODKeyFromAvatarId(avatarId: number): GOODCharacterKey {
  return textToGOODKey(
    textMap[characterMap[avatarId].NameTextMapHash.toString()]
  ) as GOODCharacterKey;
}

function getGOODKeyFromWeaponNameTextMapHash(
  weaponNameTextMapHash: string
): GOODWeaponKey {
  return textToGOODKey(textMap[weaponNameTextMapHash]) as GOODWeaponKey;
}

function getGOODKeyFromReliquaryNameTextMapHash(
  reliquaryNameTextMapHash: string
): GOODArtifactSetKey {
  const res = textToGOODKey(
    textMap[reliquaryNameTextMapHash]
  ) as GOODArtifactSetKey;
  return res;
}

function getGOODSubstatsFromReliquarySubstats(
  reliquarySubstats: {
    appendPropId: FightProp;
    statValue: number;
  }[]
): ISubstat[] {
  if (reliquarySubstats.length == 0 || reliquarySubstats.length > 4) {
    return [];
  }
  let GOODSubstats: ISubstat[] = [];
  for (let reliquarySubstat of reliquarySubstats) {
    GOODSubstats.push({
      key: fightPropToGOODKey(reliquarySubstat.appendPropId),
      value: reliquarySubstat.statValue,
    });
  }
  return GOODSubstats;
}

interface ICharacterMap {
  [key: number]: {
    Element: string;
    Consts: string[];
    SkillOrder: number[];
    Skills: {
      [key: string]: string;
    };
    //useless
    ProudMap: {
      [key: string]: number;
    };
    NameTextMapHash: number;
    SideIconName: string;
    QualityType: string;
  };
}

interface IENTextMap {
  [key: string]: string;
}
