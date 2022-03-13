import { Character, defaultStats, Weapon } from "~/src/types";

import ArtifactMainStatsData from "~src/Components/Artifacts/artifact_main_gen.json";
import { characterKeyToICharacter } from "~src/Components/Character";
import { ascLvlMax, StatToIndexMap } from "~src/util";

import { ICharacter, IGOOD, GOODArtifact, StatKey } from "./goodTypes";

export interface IGOODImport {
  err: string;
  characters: Character[];
}
type rarityValue = "1" | "2" | "3" | "4" | "5";
const convertRarity: rarityValue[] = ["1", "2", "3", "4", "5"];

interface GOODGearBank {
  [key: string]: { weapon: Weapon; artifact: GOODArtifact[] };
}

export function parseFromGO(val: string): IGOODImport {
  let result: {
    err: string;
    characters: Character[];
  } = {
    err: "",
    characters: [],
  };

  if (val === "") {
    result.err = "Please paste JSON from Genshin Optimizer to continue";
    return result;
  }

  //try parsing
  let data: IGOOD;
  try {
    data = JSON.parse(val);
  } catch (e) {
    if (val === "") {
      result.err = "Please enter JSON";
      return result;
    }

    result.err = "Invalid JSON";
    return result;
  }
  if (data.source !== "Genshin Optimizer") {
    result.err = "Only databases from Genshin Optimzer accepted";
    return result;
  }
  const goodGearBank: GOODGearBank = {};

  if (data.weapons) {
    data.weapons.forEach((goodweapon) => {
      let charKey = goodKeytoSrlKey(goodweapon.location);
      if (charKey === "") {
        //skip this weapon
        return;
      }

      let importedWeapon: Weapon = {
        name: goodKeytoSrlKey(goodweapon.key),
        level: goodweapon.level,
        max_level: ascLvlMax(goodweapon.ascension),
        refine: goodweapon.refinement,
      };
      goodGearBank[charKey] = {
        weapon: importedWeapon,
        artifact: [],
      };
    });
  } else {
    result.err = "No weapons found";
  }

  //Store artifacts based on character
  if (data.artifacts) {
    data.artifacts.forEach((artifact) => {
      let charKey = goodKeytoSrlKey(artifact.location);
      if (Object.keys(goodGearBank).includes(charKey)) {
        if (goodGearBank[charKey].artifact.length < 5) {
          goodGearBank[charKey].artifact.push(artifact);
        } else {
          result.err = `Too many artifacts on ${charKey} `;
          return result;
        }
      } else if (charKey === "") {
      } else {
        // goodGearBank[charKey] = {
        //   weapon: {
        //     // SRL uses {name} field like a key for action list
        //     name: "dullblade",
        //     refine: 1,
        //     level: 1,
        //     max_level: ascLvlMax(1),
        //   },
        //   artifact: [artifact],
        // };
        goodGearBank[charKey].artifact = [artifact];
      }
    });
  }
  //build the characters
  let chars: Character[] = [];
  if (!data.characters) {
    return {
      err: "No Characters Found",
      characters: [],
    };
  }
  data.characters.forEach((c) => {
    //convert GOOD key to our key
    let char = importCharFromGOOD(c, goodGearBank);
    if (char === undefined) {
      //skip char
      return;
    }
    chars.push(char);
  });

  //sort chars by element -> name
  chars.sort((a, b) => {
    if (b.name > a.name) {
      return -1;
    }
    if (b.name < a.name) {
      return 1;
    }
    return 0;
  });

  result.characters = chars;
  return result;
}

const sumArtifactStats = (artifacts: GOODArtifact[]): number[] => {
  const totalStats = [
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
  ].slice();

  artifacts.forEach((artifact) => {
    if (artifact.mainStatKey !== "" && artifact.mainStatKey !== "def") {
      const mainStatValue =
        ArtifactMainStatsData[convertRarity[artifact.rarity - 1]][
          artifact.mainStatKey
        ][artifact.level];
      totalStats[StatToIndexMap[goodStattoSrlStat(artifact.mainStatKey)]] +=
        mainStatValue;
    } else {
      console.log("pepegaW artifact");
    }

    artifact.substats.forEach((substat) => {
      if (substat.key.includes("_")) {
        totalStats[StatToIndexMap[goodStattoSrlStat(substat.key)]] +=
          substat.value / 100;
      } else {
        totalStats[StatToIndexMap[goodStattoSrlStat(substat.key)]] +=
          substat.value;
      }
    });
  });
  return totalStats;
};

function goodStattoSrlStat(goodStat: StatKey): string {
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

const tallyArtifactSet = (
  artifacts: GOODArtifact[]
): { [key: string]: number } => {
  const setKeyTally: { [key: string]: number } = {};
  if (artifacts === undefined) {
    return {};
  }
  artifacts
    .map((artifact) => {
      return artifact.setKey;
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
};

export function importCharFromGOOD(
  goodObj: ICharacter,
  goodGearBank: GOODGearBank
): Character | undefined {
  //find char

  if (goodObj === undefined) {
    //stop here
    return undefined;
  }
  let today = new Date();
  //copy over all the attributes we care about; ignore anything
  //we don't need
  const name = goodKeytoSrlKey(goodObj.key);
  let setCount, statTotal;

  if (goodGearBank[name].artifact === undefined) {
    setCount = {};
    statTotal = defaultStats;
  } else {
    setCount = tallyArtifactSet(goodGearBank[name].artifact);
    statTotal = sumArtifactStats(goodGearBank[name].artifact);
  }
  let char = {
    name: name,
    level: goodObj.level,
    max_level: ascLvlMax(goodObj.ascension),
    element: characterKeyToICharacter[goodKeytoSrlKey(goodObj.key)].element,
    cons: goodObj.constellation,
    weapon: goodGearBank[name].weapon,
    talents: {
      attack: goodObj.talent.auto,
      skill: goodObj.talent.skill,
      burst: goodObj.talent.burst,
    },
    //need to sum stats
    stats: statTotal,
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: setCount,
    date_added: today.toLocaleDateString(),
  };

  return char;
}

export function goodKeytoSrlKey(s: string) {
  switch (s) {
    case "KaedeharaKazuha":
      return "kazuha";
    case "KamisatoAyaka":
      return "ayaka";
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
  const result = s
    .toString()
    .replace(/[^0-9a-z]/gi, "")
    .toLowerCase();
  return result;
}
