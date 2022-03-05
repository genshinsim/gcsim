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