import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";
import { Character } from "~src/types";

export const characterKeyToICharacter: { [key: string]: ICharacter } = {
  aether: {
    key: "aether",
    name: "Aether",
    element: "none",
    weapon_type: "sword",
  },
  lumine: {
    key: "lumine",
    name: "Lumine",
    element: "none",
    weapon_type: "sword",
  },
  albedo: {
    key: "albedo",
    name: "Albedo",
    element: "geo",
    weapon_type: "sword",
  },
  aloy: { key: "aloy", name: "Aloy", element: "cryo", weapon_type: "bow" },
  amber: { key: "amber", name: "Amber", element: "pyro", weapon_type: "bow" },
  barbara: {
    key: "barbara",
    name: "Barbara",
    element: "hydro",
    weapon_type: "catalyst",
  },
  beidou: {
    key: "beidou",
    name: "Beidou",
    element: "electro",
    weapon_type: "claymore",
  },
  bennett: {
    key: "bennett",
    name: "Bennett",
    element: "pyro",
    weapon_type: "sword",
  },
  chongyun: {
    key: "chongyun",
    name: "Chongyun",
    element: "cryo",
    weapon_type: "claymore",
  },
  diluc: {
    key: "diluc",
    name: "Diluc",
    element: "pyro",
    weapon_type: "claymore",
  },
  diona: { key: "diona", name: "Diona", element: "cryo", weapon_type: "bow" },
  eula: { key: "eula", name: "Eula", element: "cryo", weapon_type: "claymore" },
  fischl: {
    key: "fischl",
    name: "Fischl",
    element: "electro",
    weapon_type: "bow",
  },
  ganyu: { key: "ganyu", name: "Ganyu", element: "cryo", weapon_type: "bow" },
  hutao: {
    key: "hutao",
    name: "Hu Tao",
    element: "pyro",
    weapon_type: "polearm",
  },
  jean: { key: "jean", name: "Jean", element: "anemo", weapon_type: "sword" },
  kazuha: {
    key: "kazuha",
    name: "Kaedehara Kazuha",
    element: "anemo",
    weapon_type: "sword",
  },
  kaeya: { key: "kaeya", name: "Kaeya", element: "cryo", weapon_type: "sword" },
  ayaka: {
    key: "ayaka",
    name: "Kamisato Ayaka",
    element: "cryo",
    weapon_type: "sword",
  },
  keqing: {
    key: "keqing",
    name: "Keqing",
    element: "electro",
    weapon_type: "sword",
  },
  klee: { key: "klee", name: "Klee", element: "pyro", weapon_type: "catalyst" },
  sara: {
    key: "sara",
    name: "Kujou Sara",
    element: "electro",
    weapon_type: "bow",
  },
  lisa: {
    key: "lisa",
    name: "Lisa",
    element: "electro",
    weapon_type: "catalyst",
  },
  mona: {
    key: "mona",
    name: "Mona",
    element: "hydro",
    weapon_type: "catalyst",
  },
  ningguang: {
    key: "ningguang",
    name: "Ningguang",
    element: "geo",
    weapon_type: "catalyst",
  },
  noelle: {
    key: "noelle",
    name: "Noelle",
    element: "geo",
    weapon_type: "claymore",
  },
  qiqi: { key: "qiqi", name: "Qiqi", element: "cryo", weapon_type: "sword" },
  raiden: {
    key: "raiden",
    name: "Raiden Shogun",
    element: "electro",
    weapon_type: "polearm",
  },
  razor: {
    key: "razor",
    name: "Razor",
    element: "electro",
    weapon_type: "claymore",
  },
  rosaria: {
    key: "rosaria",
    name: "Rosaria",
    element: "cryo",
    weapon_type: "polearm",
  },
  kokomi: {
    key: "kokomi",
    name: "Sangonomiya Kokomi",
    element: "hydro",
    weapon_type: "catalyst",
  },
  sayu: {
    key: "sayu",
    name: "Sayu",
    element: "anemo",
    weapon_type: "claymore",
  },
  sucrose: {
    key: "sucrose",
    name: "Sucrose",
    element: "anemo",
    weapon_type: "catalyst",
  },
  tartaglia: {
    key: "tartaglia",
    name: "Tartaglia",
    element: "hydro",
    weapon_type: "bow",
  },
  thoma: {
    key: "thoma",
    name: "Thoma",
    element: "pyro",
    weapon_type: "polearm",
  },
  venti: { key: "venti", name: "Venti", element: "anemo", weapon_type: "bow" },
  xiangling: {
    key: "xiangling",
    name: "Xiangling",
    element: "pyro",
    weapon_type: "polearm",
  },
  xiao: { key: "xiao", name: "Xiao", element: "anemo", weapon_type: "polearm" },
  xingqiu: {
    key: "xingqiu",
    name: "Xingqiu",
    element: "hydro",
    weapon_type: "sword",
  },
  xinyan: {
    key: "xinyan",
    name: "Xinyan",
    element: "pyro",
    weapon_type: "claymore",
  },
  yaemiko: {
    key: "yaemiko",
    name: "Yae Miko",
    element: "electro",
    weapon_type: "catalyst",
  },
  yanfei: {
    key: "yanfei",
    name: "Yanfei",
    element: "pyro",
    weapon_type: "catalyst",
  },
  yoimiya: {
    key: "yoimiya",
    name: "Yoimiya",
    element: "pyro",
    weapon_type: "bow",
  },
  zhongli: {
    key: "zhongli",
    name: "Zhongli",
    element: "geo",
    weapon_type: "polearm",
  },
  gorou: { key: "gorou", name: "Gorou", element: "geo", weapon_type: "bow" },
  itto: {
    key: "itto",
    name: "Arataki Itto",
    element: "geo",
    weapon_type: "claymore",
  },
  shenhe: {
    key: "shenhe",
    name: "Shenhe",
    element: "cryo",
    weapon_type: "polearm",
  },
  yunjin: {
    key: "yunjin",
    name: "Yun Jin",
    element: "geo",
    weapon_type: "polearm",
  },
};

export interface ICharacter {
  key: string;
  name: string;
  element: string;
  weapon_type: string;
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

export const items: Character[] = Object.keys(characterKeyToICharacter).map(
  (key) => {
    return newChar(key);
  }
);

export const render: ItemRenderer<Character> = (
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
      label={`${
        item.date_added
          ? item.element.concat(`, Imported: ${item.date_added}`)
          : item.element
      }`}
      key={`${item.date_added ? item.name.concat(item.date_added) : item.name}`}
      onClick={handleClick}
      text={highlightText(item.name, query)}
    />
  );
};

export const filter: ItemPredicate<ICharacter> = (
  query,
  item,
  _index,
  exactMatch
) => {
  const normalizedQuery = query.toLowerCase();

  if (exactMatch) {
    return item.key === normalizedQuery;
  } else {
    return (
      `${item.key} ${item.name} ${item.element}`.indexOf(normalizedQuery) >= 0
    );
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

export const characterSelectProps = {
  itemRenderer: render,
  items: items,
};
