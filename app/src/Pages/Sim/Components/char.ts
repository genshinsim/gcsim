import genshindb from "genshin-db";
import { Artifact, Weapon, Character } from "./types";
import { ICharacter, IGOOD } from "./goodTypes"

export const staticPath = {
  avatar: "/images/avatar",
  cons: "/images/avatar/cons",
  weapons: "/images/weapons",
  artifacts: "/images/artifacts",
};

export interface IGOODImport {
  err: string;
  characters: Character[];
  // selected: boolean[];
}

export function parseFromGO(val: string): IGOODImport {
  let result: {
    err: string;
    characters: Character[];
    // selected: boolean[];
  } = {
    err: "",
    characters: [],
    // selected: [],
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

  //build the characters
  let pos = new Map();
  let chars: Character[] = [];
  let sel = [];
  if (!data.characters) {
    return {
      err: "",
      characters: [],
      // selected: [],
    };
  }
  console.log("parsing characters ", data.characters);
  let trav = "";
  data.characters.forEach((c, i) => {
    //convert GOOD key to our key
    let char = charFromGOOD(c);
    if (char === undefined) {
      //skip char
      return;
    }

    //we use the imported name as key since this is what
    //should match on weapon and artifacts
    //this should handle traveler as well
    pos.set(c.key, i);

    //special check for traveler
    if (c.key === "Traveler") {
      //@ts-ignore
      c.element = c.elementKey;
      //temporarily store this so we can come back to it and add it to
      //end of the list
      trav = JSON.stringify(c);
    }

    chars.push(char);
    sel.push(false);
  });

  if (trav !== "") {
    let c = JSON.parse(trav);
    c.name = "Aether";
    let char = charFromGOOD(c);
    if (char !== undefined) {
      //shouldn't really happen other wise
      chars.push(char);
      pos.set("Aether", chars.length - 1);
      sel.push(false);
    } else {
      console.log("Unexpected error parsing traveler");
    }
  }

  //add weapons if any
  if (data.weapons) {
    console.log("parsing weapons ", data.weapons);
    data.weapons.forEach((e) => {
      if (pos.has(e.location)) {
        // console.log("adding weapon for ", e.location);
        //grab index
        let index = pos.get(e.location);
        let w = chars[index].weapon;
        let d = genshindb.weapons(e.key);

        if (d === undefined) {
          //skip this weapon
          return;
        }

        w.key = toKey(d.name);
        w.name = d.name;
        w.level = e.level;
        w.ascension = e.ascension;
        w.refinement = e.refinement;
        w.icon = `${staticPath.weapons}/${w.key}.png`;

        chars[index].weapon = w;

        //special check for traveler
        if (e.location === "Traveler") {
          index = pos.get("Aether");
          chars[index].weapon = Object.assign({}, w);
        }
      }
    });
  }

  //add artifacts if any
  if (data.artifacts) {
    console.log("parsing artifacts ", data.artifacts);
    data.artifacts.forEach((e) => {
      if (pos.has(e.location)) {
        //grab index
        let index = pos.get(e.location);
        let art = chars[index].artifact[e.slotKey];

        let d = genshindb.artifacts(e.setKey);
        if (d === undefined) {
          return;
        }

        //copy - hope this works?
        art = JSON.parse(JSON.stringify(e));
        //change set key
        art.setKey = toKey(d.name);
        art.icon = `${staticPath.artifacts}/${art.setKey}_${art.slotKey}.png`;
        // delete art.location;
        // delete art.lock;

        chars[index].artifact[e.slotKey] = art;

        //special check for traveler
        if (e.location === "Traveler") {
          index = pos.get("Aether");
          chars[index].artifact[e.slotKey] = JSON.parse(JSON.stringify(art));
        }
      }
    });
  }

  console.log("parsed results: ", chars);

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
  // result.selected = sel;
  return result;
}

export function charFromGOOD(goodObj: ICharacter): Character | undefined {
  //find char
  let d = genshindb.characters(goodObj.key);
  if (d === undefined) {
    //stop here
    return undefined;
  }
  //copy over all the attributes we care about; ignore anything
  //we don't need
  let char = blankChar();
  //convert key from db name
  char.key = toKey(d.name); //note that traveler will be set as lumine here
  char.icon = `${staticPath.avatar}/${char.key}.png`;
  //copy db data
  char.name = d.name;
  char.element = d.element;
  char.weapontype = d.weapontype;
  //copy attributes from GOOD import
  char.level = goodObj.level;
  char.constellation = goodObj.constellation;
  char.ascension = goodObj.ascension;
  char.talent.auto = goodObj.talent.auto;
  char.talent.skill = goodObj.talent.skill;
  char.talent.burst = goodObj.talent.burst;

  return char;
}

export function blankChar(): Character {
  return {
    key: "", //
    name: "", //display name
    element: "",
    icon: "",
    level: 80,
    constellation: 0,
    ascension: 6,
    talent: {
      auto: 6,
      skill: 6,
      burst: 6,
    },
    weapontype: "",
    weapon: blankWeapon(),
    artifact: blankArtifacts(),
  };
}

export function blankWeapon(): Weapon {
  return {
    key: "", //"CrescentPike"
    name: "",
    icon: "",
    level: 1, //1-90 inclusive
    ascension: 0, //0-6 inclusive. need to disambiguate 80/90 or 80/80
    refinement: 1, //1-5 inclusive
    // location: "", //where "" means not equipped.
    // lock: false, //Whether the weapon is locked in game.
  };
}

export function blankArtifacts(): {
  flower: Artifact;
  plume: Artifact;
  sands: Artifact;
  goblet: Artifact;
  circlet: Artifact;
} {
  return {
    flower: {
      setKey: "", //e.g. "GladiatorsFinale"
      slotKey: "flower", //e.g. "plume"
      icon: "",
      level: 20, //0-20 inclusive
      rarity: 5, //1-5 inclusive
      mainStatKey: "hp",
      // location: "", //where "" means not equipped.
      // lock: false, //Whether the artifact is locked in game.
      substats: [],
    },
    plume: {
      setKey: "", //e.g. "GladiatorsFinale"
      slotKey: "plume", //e.g. "plume"
      icon: "",
      level: 20, //0-20 inclusive
      rarity: 5, //1-5 inclusive
      mainStatKey: "atk",
      // location: "", //where "" means not equipped.
      // lock: false, //Whether the artifact is locked in game.
      substats: [],
    },
    sands: {
      setKey: "", //e.g. "GladiatorsFinale"
      slotKey: "sands", //e.g. "plume"
      icon: "",
      level: 20, //0-20 inclusive
      rarity: 5, //1-5 inclusive
      mainStatKey: "",
      // location: "", //where "" means not equipped.
      // lock: false, //Whether the artifact is locked in game.
      substats: [],
    },
    goblet: {
      setKey: "", //e.g. "GladiatorsFinale"
      slotKey: "goblet", //e.g. "plume"
      icon: "",
      level: 20, //0-20 inclusive
      rarity: 5, //1-5 inclusive
      mainStatKey: "",
      // location: "", //where "" means not equipped.
      // lock: false, //Whether the artifact is locked in game.
      substats: [],
    },
    circlet: {
      setKey: "", //e.g. "GladiatorsFinale"
      slotKey: "circlet", //e.g. "plume"
      icon: "",
      level: 20, //0-20 inclusive
      rarity: 5, //1-5 inclusive
      mainStatKey: "",
      // location: "", //where "" means not equipped.
      // lock: false, //Whether the artifact is locked in game.
      substats: [],
    },
  };
}

export const defaultWeapons = {
  Sword: "Dull Blade",
  Claymore: "Waster Greatsword",
  Bow: "Hunter's Bow",
  Catalyst: "Apprentice's Notes	",
  Polearm: "Beginner's Protector",
};

export function toKey(s: string) {
  return s.replace(/[^0-9a-z]/gi, "").toLowerCase();
}
