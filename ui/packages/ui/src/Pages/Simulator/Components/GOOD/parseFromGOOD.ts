/* eslint-disable prefer-const */
import { GOODArtifact, GOODCharacter, GOODCharacterKey, GOODWeapon, IGOOD } from "./GOODTypes";
import { equipArtifacts, GOODChartoSrlChar, GOODWeapontoSrlWeapon } from "./GOODToSrl.functions";
import { Character, Weapon } from "@gcsim/types";

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
  const result: {
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

  const weaponBank: WeaponBank = extractWeapons(data.weapons);

  let artifactBank: GOODArtifactBank = {};
  if (data.artifacts) {
    artifactBank = extractArtifacts(data.artifacts);
  }

  result.characters = buildCharactersFromGOOD(data.characters, weaponBank, artifactBank);

  console.log(result)
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
    // eslint-disable-next-line prefer-const
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
  let travelerIdx: { [key in string]: number } = {}
  goodChars.forEach((goodChar, index) => {
    if (goodChar.key.startsWith("Traveler")) {
      travelerIdx[goodChar.key] = index
    }
    let char = GOODChartoSrlChar(goodChar, weaponBank[goodChar.key]);

    if (char === undefined) {
      //skip char
      return;
    }
    char = equipArtifacts(char, goodArtifactBank[goodChar.key]);

    result.push(char);
  });

  //this code sucks, kids do not do this
  for (const [goodkey, idx] of Object.entries(travelerIdx)) {
    console.log("parsing ", goodkey)
    const g = goodChars[idx]
    travelers.forEach(ck => {
      let key = goodkey 
      key = key.toLowerCase()
      //split the string between traveler and element; if no element
      key  = key.replace("traveler", ck) 

      console.log("adding: ", key)
  
      let copy: GOODCharacter = {
        ...g,
        talent: {
          ...g.talent
        },
        key: key,
      }
  
      //weapon and artifact bank uses Traveler as key ignoring element
      let char = GOODChartoSrlChar(copy, weaponBank["Traveler"])
      if (char === undefined) {
        console.log(key, "not found")
        //skip char
        return;
      }
      char = equipArtifacts(char, goodArtifactBank["Traveler"])
 
      result.push(char)
    })
  }

  return result;
}

const travelers : string[] = ["lumine", "aether"]