import charPipelineData from "./char_data.generated.json";
import artifactPipelineData from "./artifact_data.generated.json";
import weaponPipelineData from "./weapon_data.generated.json";

export interface CharDataMap {
  [key: string]: {
    rarity: number;
    element: string;
    weapon_class: string;
  };
}

export function protoEleToDisplayString(ele: string): string {
  switch (ele) {
    case "Electric":
      return "electro";
    case "Fire":
      return "pyro";
    case "Ice":
      return "cryo";
    case "Water":
      return "hydro";
    case "Grass":
      return "dendro";
    case "Wind":
      return "anemo";
    case "Rock":
      return "geo";
    default:
      return "unknown";
  }
}

export function protoWeapTypeToDisplayString(w: string): string {
  switch (w) {
    case "WEAPON_SWORD_ONE_HAND":
      return "sword";
    case "WEAPON_CLAYMORE":
      return "claymore";
    case "WEAPON_POLE":
      return "polearm";
    case "WEAPON_BOW":
      return "bow";
    case "WEAPON_CATALYST":
      return "catalyst";
    default:
      return "unknown";
  }
}

export const valid_characters: string[] = Object.keys(charPipelineData.data);
export const valid_artifacts: string[] = Object.keys(artifactPipelineData.data);
export const valid_weapons: string[] = Object.keys(weaponPipelineData.data);
// TODO: maybe move these somewhere else?
export const valid_actions: string[] = ["attack", "charge", "aim", "skill", "burst", "low_plunge", "high_plunge", "dash", "jump", "walk", "swap"];
export const valid_stats: string[] = ["hp", "hp%", "atk", "atk%", "def", "def%", "cr", "cd", "er", "heal", "em", "phys%", "pyro%", "electro%", "hydro%", "dendro%", "anemo%", "geo%", "cryo%"];

let charData: CharDataMap = {};

for (const [k, v] of Object.entries(charPipelineData.data)) {
  charData[k] = {
    rarity: v.rarity === "QUALITY_ORANGE" ? 5 : 4,
    element: protoEleToDisplayString(v.element),
    weapon_class: protoWeapTypeToDisplayString(v.weapon_class),
  };
}

export const CharMap: CharDataMap = charData;

import ArtifactMainStatsData from "./artifact_main_gen.json";

export { ArtifactMainStatsData };
