import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";

export const artifactKeyToName = {
  English: {
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
  },
  Chinese: {
    archaicpetra: "悠古的磐岩",
    berserker: "战狂",
    blizzardstrayer: "冰风迷途的勇士",
    bloodstainedchivalry: "染血的骑士道",
    braveheart: "勇士之心",
    crimsonwitchofflames: "炽烈的炎之魔女",
    defenderswill: "守护之心",
    emblemofseveredfate: "绝缘之旗印",
    gambler: "赌徒",
    gladiatorsfinale: "角斗士的终幕礼",
    heartofdepth: "沉沦之心",
    huskofopulentdreams: "华馆梦醒形骸记",
    instructor: "教官",
    lavawalker: "渡过烈火的贤人",
    maidenbeloved: "被怜爱的少女",
    martialartist: "武人",
    noblesseoblige: "昔日宗室之仪",
    oceanhuedclam: "海染砗磲",
    paleflame: "苍白之火",
    prayersfordestiny: "祭水之人",
    prayersforillumination: "祭火之人",
    prayersforwisdom: "祭雷之人",
    prayerstospringtime: "祭冰之人",
    resolutionofsojourner: "行者之心",
    retracingbolide: "逆飞的流星",
    scholar: "学士",
    shimenawasreminiscence: "追忆之注连",
    tenacityofthemillelith: "千岩牢固",
    theexile: "流放者",
    thunderingfury: "如雷的盛怒",
    thundersoother: "平息鸣雷的尊者",
    tinymiracle: "奇迹",
    viridescentvenerer: "翠绿之影",
    wandererstroupe: "流浪大地的乐团",
  },
};

export interface IArtifact {
  key: string;
  name: string;
}

export const items = {
  English: Object.keys(artifactKeyToName.English).map((k) => ({
    key: k,
    name: artifactKeyToName.English[k],
  })),
  Chinese: Object.keys(artifactKeyToName.Chinese).map((k) => ({
    key: k,
    name: artifactKeyToName.Chinese[k],
  })),
};

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
  English: {
    itemPredicate: filter,
    itemRenderer: render,
    items: items.English,
  },
  Chinese: {
    itemPredicate: filter,
    itemRenderer: render,
    items: items.Chinese,
  },
};
