import { Character, Weapon } from "~/src/types";

import {
  GOODArtifact,
  GOODCharacter,
  GOODCharacterKey,
  GOODWeapon,
  IGOOD,
} from "./GOODTypes";
import {
  equipArtifacts,
  GOODChartoSrlChar,
  GOODWeapontoSrlWeapon,
} from "./GOODToSrl.functions";

export interface IGOODImport {
  err: string;
  characters: Character[];
}

type WeaponBank = {
  [char in GOODCharacterKey]?: Weapon;
};
//currently don't have inhouse artifact type
type GOODArtifactBank = {
  [char in GOODCharacterKey]?: GOODArtifact[];
};

export function parseFromGOOD(val: string): IGOODImport {
  let result: {
    err: string;
    characters: Character[];
  } = {
    err: "",
    characters: [],
  };

  if (val === "") {
    result.err = "Please paste JSON in GOOD format to continue";
    return result;
  }

  //try parsing
  let data: IGOOD;
  try {
    data = JSON.parse(val);
  } catch (e) {
    result.err = "Invalid JSON";
    return result;
  }
  if (!data.characters) {
    return {
      err: "No Characters Found",
      characters: [],
    };
  }
  if (!data.weapons) {
    return {
      err: "No Weapons Found",
      characters: [],
    };
  }

  let weaponBank: WeaponBank = extractWeapons(data.weapons);

  let artifactBank: GOODArtifactBank = {};
  if (data.artifacts) {
    artifactBank = extractArtifacts(data.artifacts);
  }

  result.characters = buildCharactersFromGOOD(
    data.characters,
    weaponBank,
    artifactBank
  );

  return result;
}

const extractWeapons = (weapons: GOODWeapon[]): WeaponBank => {
  const result: WeaponBank = {};
  weapons.forEach((goodWeapon) => {
    let GOODCharKey = goodWeapon.location;
    if (GOODCharKey !== "") {
      result[GOODCharKey] = GOODWeapontoSrlWeapon(goodWeapon);
    }
  });
  return result;
};

const extractArtifacts = (artifacts: GOODArtifact[]): GOODArtifactBank => {
  const result: GOODArtifactBank = {};
  artifacts.forEach((goodArtifact) => {
    let GOODCharKey = goodArtifact.location;
    if (GOODCharKey === "") {
      return;
    } else {
      if (result[GOODCharKey] === undefined) {
        result[GOODCharKey] = [goodArtifact];
      } else {
        result[GOODCharKey]?.push(goodArtifact);
      }
    }
  });
  return result;
};

function buildCharactersFromGOOD(
  goodChars: GOODCharacter[],
  weaponBank: WeaponBank,
  goodArtifactBank: GOODArtifactBank
) {
  const result: Character[] = [];
  goodChars.forEach((goodChar) => {
    let char = GOODChartoSrlChar(goodChar, weaponBank[goodChar.key]);
    if (char === undefined) {
      //skip char
      return;
    }
    char = equipArtifacts(char, goodArtifactBank[goodChar.key]);

    result.push(char);
  });
  return result;
}
