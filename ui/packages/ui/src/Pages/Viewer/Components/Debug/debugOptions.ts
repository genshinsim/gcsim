export const DefaultDebugOptions = [
  "damage",
  "element",
  "action",
  "energy",
  "pre_damage_mods",
  "status",
];

export const AllDebugOptions = [
  //basic stuff
  "action", //character actions
  "damage", //character damage
  "energy", //energy regen etc..
  "warning", // sim warnings; things went wrong
  //advanced stuff
  "status", // various status/buffs/debuffs
  "cooldown", // tracking things going on and off cooldown
  "element", // tracking element applications
  "shield", // shield creation
  "construct", // construct creation
  "player", // track player stuff such as stam/swapcd
  "user",
  //verbose stuff
  "heal", // healing events
  "hurt", // taking dmg events
  "pre_damage_mods", //yunjin and shenhe still uses this
  "icd", // ele and dmg app icd
  "calc", // detailed damage calc
  "snapshot", // detailed snapshot calc
  "character",
  "weapon",
  "enemy",
  "artifact",
  //debug stuff
  "debug",
  "sim",
  "hitlag",
  //these options have been deprecated
  "hook",
  "procs", //don't think this was ever used
  "task",
  "reaction",
  "snapshot_mods",
  "queue",
];

export const SimplePreset = ["action", "damage", "energy", "warning"];

//include cooldowns,
export const AdvancedPreset = [
  ...SimplePreset,
  "status",
  "cooldown",
  "element",
  "shield",
  "construct",
  "user",
];

export const VerbosePreset = [
  ...AdvancedPreset,
  "player",
  "heal",
  "hurt",
  "pre_damage_mod",
  "icd", // ele and dmg app icd
  "calc", // detailed damage calc
  "snapshot", // detailed snapshot calc
  "character",
  "weapon",
  "enemy",
  "artifact",
];

export const DebugPreset = [...VerbosePreset, "debug", "sim", "hitlag"];
