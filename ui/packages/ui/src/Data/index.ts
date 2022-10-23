import characterData from "./character_data.json";

export interface CharDataMap {
  [key: string]: {
    rarity: number;
    element: string;
    weapon_class: string;
  };
}

export const CharMap: CharDataMap = characterData;

export function TransformTravelerKeyToName(key: string): string {
  if (key.startsWith("aether")) {
    return "aether";
  }
  if (key.startsWith("lumine")) {
    return "lumine";
  }

  return key;
}

export function TravelerCheck(key: string) {
  return key.startsWith("aether") || key.startsWith("lumine") || key.startsWith("traveler");
}

import ArtifactMainStatsData from "./artifact_main_gen.json";

export { ArtifactMainStatsData };
