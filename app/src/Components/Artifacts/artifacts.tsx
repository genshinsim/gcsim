import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";

export const artifactKeyToName: { [key: string]: string } = {
  archaicpetra: "Archaic Petra",
  berserker: "Berserker",
  blizzardstrayer: "Blizzard Strayer",
  bloodstainedchivalry: "Bloodstained Chivalry",
  braveheart: "Brave Heart",
  crimsonwitchofflames: "Crimson Witch of Flames",
  defenderswill: "Defender's Will",
  emblemofseveredfate: "Emblem of Severed Fate",
  gambler: "Gambler",
  gladiatorsfinale: "Gladiator's Finale",
  heartofdepth: "Heart of Depth",
  huskofopulentdreams: "Husk of Opulent Dreams",
  instructor: "Instructor",
  lavawalker: "Lavawalker",
  maidenbeloved: "Maiden Beloved",
  martialartist: "Martial Artist",
  noblesseoblige: "Noblesse Oblige",
  oceanhuedclam: "Ocean-Hued Clam",
  paleflame: "Pale Flame",
  prayersfordestiny: "Prayers for Destiny",
  prayersforillumination: "Prayers for Illumination",
  prayersforwisdom: "Prayers for Wisdom",
  prayerstospringtime: "Prayers to Springtime",
  resolutionofsojourner: "Resolution of Sojourner",
  retracingbolide: "Retracing Bolide",
  scholar: "Scholar",
  shimenawasreminiscence: "Shimenawa's Reminiscence",
  tenacityofthemillelith: "Tenacity of the Millelith",
  theexile: "The Exile",
  thunderingfury: "Thundering Fury",
  thundersoother: "Thundersoother",
  tinymiracle: "Tiny Miracle",
  viridescentvenerer: "Viridescent Venerer",
  wandererstroupe: "Wanderer's Troupe",
};

export interface IArtifact {
  key: string;
  name: string;
}

export const items: IArtifact[] = Object.keys(artifactKeyToName).map((k) => ({
  key: k,
  name: artifactKeyToName[k],
}));

export const render: ItemRenderer<IArtifact> = (
  item,
  { handleClick, modifiers, query }
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={item.name}
      key={item.key}
      onClick={handleClick}
      text={highlightText(item.name, query)}
    />
  );
};

export const filter: ItemPredicate<IArtifact> = (
  query,
  item,
  _index,
  exactMatch
) => {
  const normalizedQuery = query.toLowerCase();

  if (exactMatch) {
    return item.key === normalizedQuery;
  } else {
    return `${item.key} ${item.name}`.indexOf(normalizedQuery) >= 0;
  }
};

function escapeRegExpChars(text: string) {
  return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1");
}

function highlightText(text: string, query: string) {
  let lastIndex = 0;
  const words = query
    .split(/\s+/)
    .filter((word) => word.length > 0)
    .map(escapeRegExpChars);
  if (words.length === 0) {
    return [text];
  }
  const regexp = new RegExp(words.join("|"), "gi");
  const tokens: React.ReactNode[] = [];
  while (true) {
    const match = regexp.exec(text);
    if (!match) {
      break;
    }
    const length = match[0].length;
    const before = text.slice(lastIndex, regexp.lastIndex - length);
    if (before.length > 0) {
      tokens.push(before);
    }
    lastIndex = regexp.lastIndex;
    tokens.push(<strong key={lastIndex}>{match[0]}</strong>);
  }
  const rest = text.slice(lastIndex);
  if (rest.length > 0) {
    tokens.push(rest);
  }
  return tokens;
}

export const artifactSelectProps = {
  itemPredicate: filter,
  itemRenderer: render,
  items: items,
};
