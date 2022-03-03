import genshindb from "genshin-db";
// import { Artifact, Weapon, Character } from "./types";
import { Character, maxStatLength, Talent, Weapon } from "~/src/types";
import { characterKeyToICharacter } from "~src/Components/Character";
import { ascLvlMax } from "~src/util";

import { ICharacter, IGOOD, GOODArtifact } from "./goodTypes"

export const staticPath = {
  avatar: "/images/avatar",
  cons: "/images/avatar/cons",
  weapons: "/images/weapons",
  artifacts: "/images/artifacts",
};

interface IartifactBank {[srlcharkey : string]:GOODArtifact[] }

export interface IGOODImport {
  err: string;
  characters: Character[];
}

interface GOODGearBank {
  [key: string] : {weapon: Weapon, artifact: GOODArtifact[]}
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
console.log("parse", data)
//Store artifacts based on character
//add artifacts if any
const artifactBank: IartifactBank={}
if (data.artifacts) {
    console.log("parsing artifacts ", data.artifacts);
    data.artifacts.forEach((artifact) => {
        let charKey  = convertFromGOODKey(artifact.location)
        if (Object.keys(artifactBank).includes(charKey) ){
          
          if(artifactBank[charKey].length <5)
          {artifactBank[charKey].push(artifact) }
          else{
            result.err = `Too many artifacts on ${charKey} `
            return result
          }
        }
        else if(charKey===""){        }
          else{
            artifactBank[charKey] = [artifact];
          }
        // //special check for traveler
        // if (e.location === "Traveler") {
        //   index = pos.get("Aether");
        //   chars[index].artifact[e.slotKey] = JSON.parse(JSON.stringify(art));
        // }
      });
      console.log("parsed results arts: ", artifactBank);
  }

  //build the characters
  // let pos = new Map();
  let chars: Character[] = [];
  if (!data.characters) {
    return {
      err: "No Characters Found",
      characters: [],
    };
  }
  console.log("parsing characters ", data.characters);
  let trav = "";
  data.characters.forEach((c, i) => {
    //convert GOOD key to our key
    let char = importCharFromGOOD(c,artifactBank);
    if (char === undefined) {
      //skip char
      return;
    }

    //we use the imported name as key since this is what
    //should match on weapon and artifacts
    //this should handle traveler as well
    // pos.set(c.key, i);

    //special check for traveler
    if (c.key === "Traveler") {
      //@ts-ignore
      c.element = c.elementKey;
      //temporarily store this so we can come back to it and add it to
      //end of the list
      trav = JSON.stringify(c);
    }

    chars.push(char);

  });
console.log("after Char", chars)
    if (trav !== "") {
    // let c = JSON.parse(trav);
    // c.name = "Aether";
    // let char = charFromGOOD(c);
    // if (char !== undefined) {
    //   //shouldn't really happen other wise
    //   chars.push(char);
    //   // pos.set("Aether", chars.length - 1);
    } else {
      console.log("Unexpected error parsing traveler");
    }
  
  

  // //add weapons if any
  // if (data.weapons) {
  //   console.log("parsing weapons ", data.weapons);
  //   data.weapons.forEach((e) => {
  //     if (pos.has(e.location)) {
  //       // console.log("adding weapon for ", e.location);
  //       //grab index
  //       let index = pos.get(e.location);
  //       let w = chars[index].weapon;

  //       if (d === undefined) {
  //         //skip this weapon
  //         return;
  //       }

  //       w.name = convertFromGOODKey(e.key);
  //       w.name = d.name;
  //       w.level = e.level;
  //       w.refine = e.refinement;

  //       chars[index].weapon = w;

  //       //special check for traveler
  //       if (e.location === "Traveler") {
  //         index = pos.get("Aether");
  //         chars[index].weapon = Object.assign({}, w);
  //       }
  //     }
  //   });
  // }

  

  // console.log("parsed results: ", chars);
  // // 

  // //sort chars by element -> name
  // chars.sort((a, b) => {
  //   if (b.name > a.name) {
  //     return -1;
  //   }
  //   if (b.name < a.name) {
  //     return 1;
  //   }
  //   return 0;
  // });

  // result.characters = chars;
  // result.selected = sel;
  return result;
}

const tallyArtifactSet = (artifacts: GOODArtifact[]): {[key: string]: number}=>{
  const setKeyTally: {[key: string]: number} = {};
  if(artifacts === undefined ){
    return {}
  }
  artifacts.map((artifact) => {return artifact.setKey}) 
  .map((setKey) => {
    if (Object.keys(setKeyTally).includes(setKey) ){
    setKeyTally[setKey] += 1 }
    else if(setKey!=""){
      setKeyTally[setKey] = 1;
    }
  });// Tallies the set keys

  // Clamps artifact set value for better handling down the line #blamesrl
  Object.keys(setKeyTally).forEach(setKey => {
    if(setKeyTally[setKey] < 2){
      delete setKeyTally[setKey]
    }
    else if(setKeyTally[setKey] > 2 && setKeyTally[setKey] < 4 ){
      setKeyTally[setKey]= 2
    }
    else if(setKeyTally[setKey] > 4){
      setKeyTally[setKey]= 4
    }
  });
  return setKeyTally
}


export function importCharFromGOOD(goodObj: ICharacter, artifactBank: IartifactBank): Character | undefined {
  //find char
  if (goodObj === undefined) {
    //stop here
    return undefined;
  }
  //copy over all the attributes we care about; ignore anything
  //we don't need
  const name = convertFromGOODKey(goodObj.key)
  const setcount = tallyArtifactSet(artifactBank[name])
  let char = {name: name,
    level: goodObj.level,
    max_level: ascLvlMax(goodObj.level),
    element: characterKeyToICharacter[convertFromGOODKey(goodObj.key)].element,
    cons: goodObj.constellation,
    weapon: {
      // SRL uses {name} field like a key for action list
      name: "dullblade",
      refine: 1,
      level: 1,
      max_level: ascLvlMax(1),
    },
    talents: {
      attack: goodObj.talent.auto,
      skill: goodObj.talent.skill,
      burst: goodObj.talent.burst,
    },
    //need to sum stats
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    // need to sum arti sets
    sets: setcount,}

  return char;
}
const newChar = (name: string): Character => {
  const c = characterKeyToICharacter[name];
  //default weapons
  return {
    name: name,
    level: 80,
    max_level: 90,
    element: c.element,
    cons: 0,
    weapon: {
      name: "dullblade",
      refine: 1,
      level: 1,
      max_level: 20,
    },
    talents: {
      attack: 6,
      skill: 6,
      burst: 6,
    },
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: {},
  };
};
// export function blankChar(): Character {
//   return {
//     key: "", //
//     name: "", //display name
//     element: "",
//     icon: "",
//     level: 80,
//     constellation: 0,
//     ascension: 6,
//     talent: {
//       auto: 6,
//       skill: 6,
//       burst: 6,
//     },
//     weapontype: "",
//     weapon: blankWeapon(),
//     artifact: blankArtifacts(),
//   };
// }

// export function blankWeapon(): Weapon {
//   return {
//     // key: "", //"CrescentPike"
//     name: "",
//     // icon: "",
//     level: 1, //1-90 inclusive
//     // ascension: 0, //0-6 inclusive. need to disambiguate 80/90 or 80/80
//     refine: 1, //1-5 inclusive
//     // location: "", //where "" means not equipped.
//     // lock: false, //Whether the weapon is locked in game.
//   };
// }

// export function blankArtifacts(): {
//   flower: GOODArtifact;
//   plume: GOODArtifact;
//   sands: GOODArtifact;
//   goblet: GOODArtifact;
//   circlet: GOODArtifact;
// } {
//   return {
//     flower: {
//       setKey: "", //e.g. "GladiatorsFinale"
//       slotKey: "flower", //e.g. "plume"
//       icon: "",
//       level: 20, //0-20 inclusive
//       rarity: 5, //1-5 inclusive
//       mainStatKey: "hp",
//       // location: "", //where "" means not equipped.
//       // lock: false, //Whether the artifact is locked in game.
//       substats: [],
//     },
//     plume: {
//       setKey: "", //e.g. "GladiatorsFinale"
//       slotKey: "plume", //e.g. "plume"
//       icon: "",
//       level: 20, //0-20 inclusive
//       rarity: 5, //1-5 inclusive
//       mainStatKey: "atk",
//       // location: "", //where "" means not equipped.
//       // lock: false, //Whether the artifact is locked in game.
//       substats: [],
//     },
//     sands: {
//       setKey: "", //e.g. "GladiatorsFinale"
//       slotKey: "sands", //e.g. "plume"
//       icon: "",
//       level: 20, //0-20 inclusive
//       rarity: 5, //1-5 inclusive
//       mainStatKey: "",
//       // location: "", //where "" means not equipped.
//       // lock: false, //Whether the artifact is locked in game.
//       substats: [],
//     },
//     goblet: {
//       setKey: "", //e.g. "GladiatorsFinale"
//       slotKey: "goblet", //e.g. "plume"
//       icon: "",
//       level: 20, //0-20 inclusive
//       rarity: 5, //1-5 inclusive
//       mainStatKey: "",
//       // location: "", //where "" means not equipped.
//       // lock: false, //Whether the artifact is locked in game.
//       substats: [],
//     },
//     circlet: {
//       setKey: "", //e.g. "GladiatorsFinale"
//       slotKey: "circlet", //e.g. "plume"
//       icon: "",
//       level: 20, //0-20 inclusive
//       rarity: 5, //1-5 inclusive
//       mainStatKey: "",
//       // location: "", //where "" means not equipped.
//       // lock: false, //Whether the artifact is locked in game.
//       substats: [],
//     },
//   };
// }

export const defaultWeapons = {
  Sword: "Dull Blade",
  Claymore: "Waster Greatsword",
  Bow: "Hunter's Bow",
  Catalyst: "Apprentice's Notes	",
  Polearm: "Beginner's Protector",
};

export function convertFromGOODKey(s: string) {
  switch (s){
  case "KaedeharaKazuha":
    return "kazuha"
  case "KamisatoAyaka":
    return "ayaka"
  case "KujouSara":
    return "sara"
  case "RaidenShogun":
    return "raiden"
  case "SangonomiyaKokomi":
      return "kokomi"
  case "YaeMiko":
    return "yaemiko"
  case "AratakiItto":
    return "itto"
    }
    const result = s.toString().replace(/[^0-9a-z]/gi, "").toLowerCase();
  return result
}
