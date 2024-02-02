import ArtifactDataGen from "@ui/Data/artifact_data.generated.json";
import CharDataGen from "@ui/Data/char_data.generated.json";
import WeaponDataGen from "@ui/Data/weapon_data.generated.json";

import { Character, Set, Weapon } from "@gcsim/types";
import { ArtifactMainStatsData } from "@ui/Data";
import { GOODStatKey } from "../GOOD/GOODTypes";
import {
  ascLvlMax,
  ascToMaxLvl,
  DMElementToKey,
  GOODStatToIndexMap,
} from "../util";
import {
  EnkaData,
  FightProp,
  GenshinItemReliquary,
  GenshinItemWeapon,
} from "./EnkaTypes";

const charIDToKeyLookup = Object.values(CharDataGen.data).reduce((acc, val) => {
  //this concat here is to handle traveler which has the same id but diff sub_id
  const id_str = `${val.id}${"sub_id" in val ? "-" + val["sub_id"] : ""}`;
  acc[id_str] = val.key;
  return acc;
}, {} as { [id_str: string]: string });

const weaponIDToKeyLookup = Object.values(WeaponDataGen.data).reduce(
  (acc, val) => {
    acc[val.id] = val.key;
    return acc;
  },
  {} as { [id: number]: string }
);

let artifactMapByTextMapId: { [key in string]: string } = {};
for (const [k, v] of Object.entries(ArtifactDataGen.data)) {
  artifactMapByTextMapId[v.text_map_id] = k;
}

const stats_base = [
  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

type CharData = {
  key: string;
  element: string;
  skill_details: {
    skill: number;
    burst: number;
    attack: number;
    burst_energy_cost: number;
  };
};

//findCharDataFromEnka takes id and skillmap (required to identify traveler) from enka and
//converts it to our internal generated data
function findCharDataFromEnka(
  avatarId: number,
  skillDepotId: number
): CharData {
  let converted_id = avatarId.toString();
  //if traveler, then we need to find the subid
  if (avatarId === 10000007 || avatarId === 10000005) {
    converted_id = converted_id + "-" + skillDepotId;
    console.log("using id + skill depot id for traveler", converted_id);
  }
  //sanity check that this character is implemented
  if (!(converted_id in charIDToKeyLookup)) {
    throw `character with id ${converted_id} not imported; possibly not implemented`;
  }
  const key = charIDToKeyLookup[converted_id];
  const data = {
    key: key,
    skill_details: CharDataGen.data[key].skill_details,
    element: CharDataGen.data[key].element,
  };
  console.log("character data found", data);
  return data;
}

//extract character weapon from enka equip list; return null if not found
function extractWeapon(
  equipList: (GenshinItemWeapon | GenshinItemReliquary)[]
): Weapon {
  let result: Weapon | null = null;
  equipList.forEach((e) => {
    if (e.flat.itemType != "ITEM_WEAPON") {
      return;
    }
    const { weapon: enkaWeapon, itemId } = e as GenshinItemWeapon;
    if (!(itemId in weaponIDToKeyLookup)) {
      throw `unrecognized weapon (id ${itemId})`;
    }
    const key = weaponIDToKeyLookup[itemId];
    result = {
      name: key,
      refine: determineWeaponRefinement(enkaWeapon.affixMap),
      level: enkaWeapon.level,
      max_level: ascLvlMax(enkaWeapon.promoteLevel ?? 0),
    };
  });
  if (result === null) {
    throw `no weapon found`;
  }
  return result;
}

function extractArtifactSet(
  equipList: (GenshinItemWeapon | GenshinItemReliquary)[]
): Set {
  let result: Set = {};
  equipList.forEach((e) => {
    if (e.flat.itemType != "ITEM_RELIQUARY") {
      return;
    }
    //find set, throw error if we can an unrecognized set
    if (!(e.flat.setNameTextMapHash in artifactMapByTextMapId)) {
      throw `unrecognized artifact set (id: ${e.flat.setNameTextMapHash})`;
    }
    const key = artifactMapByTextMapId[e.flat.setNameTextMapHash];
    if (!(key in result)) {
      result[key] = 0;
    }
    result[key] = result[key] + 1;
  });
  return result;
}

function extractArtifactStats(
  equipList: (GenshinItemWeapon | GenshinItemReliquary)[]
): number[] {
  //TODO: using this here so we can in future return main + sub on sep lines

  //track total
  let total = stats_base.slice();
  //we'll want to use labels for the other stuff (although doesn't do anything atm)
  equipList.forEach((e) => {
    if (e.flat.itemType != "ITEM_RELIQUARY") {
      return;
    }
    const ms = extractMainStat(e as GenshinItemReliquary);
    //TODO: we really shouldn't be doing this so roundabout...
    const ms_idx =
      GOODStatToIndexMap[
        fightPropToGOODKey(e.flat.reliquaryMainstat.mainPropId)
      ];
    total[ms_idx] += ms;
    //TODO: in future this should go into a diff slice
    //add sub stats
    const subs = extractSubStats(e as GenshinItemReliquary);
    subs.forEach((v, i) => {
      total[i] += v;
    });
  });
  return total;
}

function extractSubStats(e: GenshinItemReliquary): number[] {
  let total = stats_base.slice();
  for (const sub of e.flat.reliquarySubstats) {
    const key = fightPropToGOODKey(sub.appendPropId);
    if (!(key in GOODStatToIndexMap)) {
      continue;
    }
    const idx = GOODStatToIndexMap[key];
    let val = sub.statValue;
    //this corrects for percentages
    if (key.includes("_")) {
      val = val / 100.0;
    }
    total[idx] += val;
  }
  console.log("substats extracted", total);
  return total;
}

function extractMainStat(e: GenshinItemReliquary): number {
  const { flat, reliquary: data } = e;
  //mainstat is calculated based on lvl + rarity + stat key
  const lvl = data.level - 1; //enka returns lvl as +1, so 20 is actually 21
  const rarity = e.flat.rankLevel.toString(); //keyed as string
  const ms_key = fightPropToGOODKey(flat.reliquaryMainstat.mainPropId);
  return ArtifactMainStatsData[rarity][ms_key][lvl];
}

export default function EnkaToGOOD(enkaData: EnkaData): {
  characters: Character[];
  errors: any[];
} {
  const characters: Character[] = [];
  let errors: any[] = [];
  const today = new Date();
  try {
    enkaData.avatarInfoList.forEach(
      ({
        avatarId,
        propMap,
        skillDepotId,
        talentIdList,
        skillLevelMap,
        equipList,
      }) => {
        //try getting character data, if failed then skip this character
        let characterData: CharData;
        try {
          characterData = findCharDataFromEnka(avatarId, skillDepotId);
        } catch (e) {
          errors.push(e);
          return;
        }

        let weapon: Weapon;
        try {
          weapon = extractWeapon(equipList);
        } catch (e) {
          errors.push(`failed to import ${avatarId}: ${e}`);
          return;
        }

        let set: Set;
        try {
          set = extractArtifactSet(equipList);
        } catch (e) {
          errors.push(`failed to import ${avatarId}: ${e}`);
          return;
        }

        let stats: number[];
        try {
          stats = extractArtifactStats(equipList);
        } catch (e) {
          errors.push(`failed to import ${avatarId}: ${e}`);
          return;
        }

        let result: Character = {
          name: characterData.key,
          level: parseInt(propMap["4001"].val) ?? 1,
          element: DMElementToKey[characterData.element],
          max_level: ascToMaxLvl(parseInt(propMap["1002"].val) ?? 1),
          cons: talentIdList?.length ?? 0,
          talents: getCharacterTalentV2(
            characterData.skill_details,
            skillLevelMap
          ),
          weapon: weapon,
          sets: set,
          stats: stats,
          snapshot: stats_base.slice(),
          date_added: today.toLocaleDateString(),
        };

        characters.push(result);

        console.log(`succesfully imported ${result.name} (id: ${avatarId})`);
      }
    );
  } catch (e: any) {
    console.log(e);
    errors.push(e);
  }

  return {
    characters,
    errors,
  };
}

function determineWeaponRefinement(affixMap?: { [key: number]: number }) {
  return affixMap
    ? (Object.entries(affixMap)[0] != null
        ? Object.entries(affixMap)[0][1]
        : 0) + 1
    : 1;
}

function getCharacterTalentV2(
  skill_details: {
    skill: number;
    burst: number;
    attack: number;
  },
  skillLevelMap: { [key: number]: number }
) {
  return {
    attack: skillLevelMap[skill_details.attack],
    skill: skillLevelMap[skill_details.skill],
    burst: skillLevelMap[skill_details.burst],
  };
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
    case FightProp.FIGHT_PROP_GRASS_ADD_HURT:
      return "dendro_dmg_";
    case FightProp.FIGHT_PROP_PHYSICAL_ADD_HURT:
      return "physical_dmg_";
    default:
      return "";
  }
}
