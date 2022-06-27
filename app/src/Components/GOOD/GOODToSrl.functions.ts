import { Character, Weapon } from "~src/types";
import {
  GOODArtifact,
  GOODArtifactSetKey,
  GOODCharacter,
  GOODCharacterKey,
  GOODStatKey,
  GOODWeapon,
  GOODWeaponKey,
} from "./GOODTypes";
import ArtifactMainStatsData from "~src/Components/Artifacts/artifact_main_gen.json";
import { characterKeyToICharacter } from "~src/Components/Character";
import { ascLvlMax, StatToIndexMap } from "~src/util";
type rarityValue = "1" | "2" | "3" | "4" | "5";
const convertRarity: rarityValue[] = ["1", "2", "3", "4", "5"];

export function GOODWeapontoSrlWeapon(weapon: GOODWeapon): Weapon {
  return {
    name: GOODKeytoGCSIMKey(weapon.key),
    level: weapon.level,
    max_level: ascLvlMax(weapon.ascension),
    refine: weapon.refinement,
  };
}

export function sumArtifactStats(artifacts: GOODArtifact[]): number[] {
  const totalStats = [
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
  ].slice();

  artifacts.forEach((artifact) => {
    if (artifact.mainStatKey !== "" && artifact.mainStatKey !== "def") {
      const mainStatValue =
        ArtifactMainStatsData[convertRarity[artifact.rarity - 1]][
          artifact.mainStatKey
        ][artifact.level];
      const srlStat = GOODStattoSrlStat(artifact.mainStatKey);
      if (srlStat === undefined) return;
      totalStats[StatToIndexMap[srlStat]] += mainStatValue;
    } else {
      console.log("pepegaW artifact");
      return;
    }

    artifact.substats.forEach((substat) => {
      const srlStat = GOODStattoSrlStat(substat.key);
      if (srlStat === undefined) return;
      if (substat.key.includes("_")) {
        totalStats[StatToIndexMap[srlStat]] += substat.value / 100;
      } else {
        totalStats[StatToIndexMap[srlStat]] += substat.value;
      }
    });
  });
  return totalStats;
}

export function GOODStattoSrlStat(goodStat: GOODStatKey): string | undefined {
  switch (goodStat) {
    case "hp":
      return "HP";
    case "hp_":
      return "HPP";
    case "atk":
      return "ATK";
    case "atk_":
      return "ATKP";
    case "def":
      return "DEF";
    case "def_":
      return "DEFP";
    case "eleMas":
      return "EM";
    case "enerRech_":
      return "ER";
    case "heal_":
      return "Heal";
    case "critRate_":
      return "CR";
    case "critDMG_":
      return "CD";
    case "physical_dmg_":
      return "PhyP";
    case "anemo_dmg_":
      return "AnemoP";
    case "geo_dmg_":
      return "GeoP";
    case "electro_dmg_":
      return "ElectroP";
    case "hydro_dmg_":
      return "HydroP";
    case "pyro_dmg_":
      return "PyroP";
    case "cryo_dmg_":
      return "CryoP";
  }
}

export function tallyArtifactSet(artifacts: GOODArtifact[]): {
  [key: string]: number;
} {
  const setKeyTally: { [key: string]: number } = {};
  if (artifacts === undefined) {
    return {};
  }
  artifacts
    .map((artifact) => {
      return artifact.setKey.toLowerCase();
    })
    .map((setKey) => {
      if (Object.keys(setKeyTally).includes(setKey)) {
        setKeyTally[setKey] += 1;
      } else if (setKey != "") {
        setKeyTally[setKey] = 1;
      }
    }); // Tallies the set keys

  // Clamps artifact set value for better handling down the line #blamesrl
  Object.keys(setKeyTally).forEach((setKey) => {
    if (setKeyTally[setKey] < 2) {
      delete setKeyTally[setKey];
    } else if (setKeyTally[setKey] > 2 && setKeyTally[setKey] < 4) {
      setKeyTally[setKey] = 2;
    } else if (setKeyTally[setKey] > 4) {
      setKeyTally[setKey] = 4;
    }
  });
  return setKeyTally;
}

export function equipArtifacts(
  char: Character,
  charArtifacts: GOODArtifact[] | undefined
): Character {
  if (charArtifacts === undefined || charArtifacts.length === 0) {
    return char;
  } else {
    const sets = tallyArtifactSet(charArtifacts);
    const stats = sumArtifactStats(charArtifacts);
    return {
      ...char,
      stats,
      sets,
    };
  }
}

export function GOODChartoSrlChar(
  goodChar: GOODCharacter,
  weapon: Weapon | undefined
): Character | undefined {
  let today = new Date();
  const name = GOODKeytoGCSIMKey(goodChar.key);
  const iChar = characterKeyToICharacter[name];
  if (iChar == undefined) {
    return undefined;
  }

  return {
    name: name,
    level: goodChar.level,
    max_level: ascLvlMax(goodChar.ascension),
    element: iChar.element,
    cons: goodChar.constellation,
    weapon: weapon ?? {
      name: "dullblade",
      refine: 1,
      level: 1,
      max_level: 20,
    },
    talents: {
      attack: goodChar.talent.auto,
      skill: goodChar.talent.skill,
      burst: goodChar.talent.burst,
    },
    //need to sum stats
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: {},
    date_added: today.toLocaleDateString(),
  };
}

export function GOODKeytoGCSIMKey(
  goodKey: GOODArtifactSetKey | GOODCharacterKey | GOODWeaponKey
) {
  switch (goodKey) {
    case "KaedeharaKazuha":
      return "kazuha";
    case "KamisatoAyaka":
      return "ayaka";
    case "KamisatoAyato":
      return "ayato";
    case "KujouSara":
      return "sara";
    case "RaidenShogun":
      return "raiden";
    case "SangonomiyaKokomi":
      return "kokomi";
    case "YaeMiko":
      return "yaemiko";
    case "AratakiItto":
      return "itto";
  }

  const result = goodKey
    .toString()
    .replace(/[^0-9a-z]/gi, "")
    .toLowerCase();
  return result;
}
