import { MenuItem } from '@blueprintjs/core';
import { ItemPredicate, ItemRenderer } from '@blueprintjs/select';
import i18n from 'i18next';
import { Character } from '~src/Types/sim';

export interface ICharacter {
  key: string;
  element: string;
  weapon_class: string;
  notes?: string;
}

export const characterKeyToICharacter: { [key: string]: ICharacter } = {
  // aether: {
  //   key: "aether",
  //   element: "none",
  //   weapon_type: "sword",
  // },
  // lumine: {
  //   key: "lumine",
  //   element: "none",
  //   weapon_type: "sword",
  // },
  traveler: {
    key: 'aether',
    element: 'none',
    weapon_class: 'sword',
  },
  albedo: { key: 'albedo', element: 'geo', weapon_class: 'sword' },
  aloy: { key: 'aloy', element: 'cryo', weapon_class: 'bow' },
  amber: { key: 'amber', element: 'pyro', weapon_class: 'bow' },
  barbara: {
    key: 'barbara',
    element: 'hydro',
    weapon_class: 'catalyst',
  },
  beidou: {
    key: 'beidou',
    element: 'electro',
    weapon_class: 'claymore',
  },
  bennett: { key: 'bennett', element: 'pyro', weapon_class: 'sword' },
  chongyun: {
    key: 'chongyun',
    element: 'cryo',
    weapon_class: 'claymore',
  },
  diluc: { key: 'diluc', element: 'pyro', weapon_class: 'claymore' },
  diona: { key: 'diona', element: 'cryo', weapon_class: 'bow' },
  dori: {
    key: 'dori',
    element: 'electro',
    weapon_class: 'claymore',
  },
  eula: { key: 'eula', element: 'cryo', weapon_class: 'claymore' },
  fischl: { key: 'fischl', element: 'electro', weapon_class: 'bow' },
  ganyu: { key: 'ganyu', element: 'cryo', weapon_class: 'bow' },
  hutao: { key: 'hutao', element: 'pyro', weapon_class: 'polearm' },
  jean: { key: 'jean', element: 'anemo', weapon_class: 'sword' },
  kazuha: { key: 'kazuha', element: 'anemo', weapon_class: 'sword' },
  kaeya: { key: 'kaeya', element: 'cryo', weapon_class: 'sword' },
  ayaka: { key: 'ayaka', element: 'cryo', weapon_class: 'sword' },
  keqing: { key: 'keqing', element: 'electro', weapon_class: 'sword' },
  klee: { key: 'klee', element: 'pyro', weapon_class: 'catalyst' },
  sara: { key: 'sara', element: 'electro', weapon_class: 'bow' },
  lisa: { key: 'lisa', element: 'electro', weapon_class: 'catalyst' },
  mona: { key: 'mona', element: 'hydro', weapon_class: 'catalyst' },
  ningguang: {
    key: 'ningguang',
    element: 'geo',
    weapon_class: 'catalyst',
  },
  noelle: { key: 'noelle', element: 'geo', weapon_class: 'claymore' },
  qiqi: { key: 'qiqi', element: 'cryo', weapon_class: 'sword' },
  raiden: { key: 'raiden', element: 'electro', weapon_class: 'polearm' },
  razor: { key: 'razor', element: 'electro', weapon_class: 'claymore' },
  rosaria: { key: 'rosaria', element: 'cryo', weapon_class: 'polearm' },
  kokomi: { key: 'kokomi', element: 'hydro', weapon_class: 'catalyst' },
  sayu: { key: 'sayu', element: 'anemo', weapon_class: 'claymore' },
  sucrose: {
    key: 'sucrose',
    element: 'anemo',
    weapon_class: 'catalyst',
  },
  tartaglia: { key: 'tartaglia', element: 'hydro', weapon_class: 'bow' },
  thoma: { key: 'thoma', element: 'pyro', weapon_class: 'polearm' },
  venti: { key: 'venti', element: 'anemo', weapon_class: 'bow' },
  xiangling: {
    key: 'xiangling',
    element: 'pyro',
    weapon_class: 'polearm',
  },
  xiao: { key: 'xiao', element: 'anemo', weapon_class: 'polearm' },
  xingqiu: { key: 'xingqiu', element: 'hydro', weapon_class: 'sword' },
  xinyan: { key: 'xinyan', element: 'pyro', weapon_class: 'claymore' },
  yanfei: { key: 'yanfei', element: 'pyro', weapon_class: 'catalyst' },
  yoimiya: { key: 'yoimiya', element: 'pyro', weapon_class: 'bow' },
  zhongli: { key: 'zhongli', element: 'geo', weapon_class: 'polearm' },
  gorou: { key: 'gorou', element: 'geo', weapon_class: 'bow' },
  itto: { key: 'itto', element: 'geo', weapon_class: 'claymore' },
  shenhe: { key: 'shenhe', element: 'cryo', weapon_class: 'polearm' },
  yunjin: { key: 'yunjin', element: 'geo', weapon_class: 'polearm' },
  yaemiko: {
    key: 'yaemiko',
    element: 'electro',
    weapon_class: 'catalyst',
  },
  ayato: { key: 'ayato', element: 'hydro', weapon_class: 'sword' },
  yelan: { key: 'yelan', element: 'hydro', weapon_class: 'bow' },
  kuki: { key: 'kuki', element: 'electro', weapon_class: 'sword' },
  heizou: { key: 'heizou', element: 'anemo', weapon_class: 'catalyst' },
  collei: { key: 'collei', element: 'dendro', weapon_class: 'bow' },
  tighnari: { key: 'tighnari', element: 'dendro', weapon_class: 'bow' },
  travelerelectro: {
    key: 'travelerelectro',
    element: 'electro',
    weapon_class: 'Sword',
  },
  traveleranemo: {
    key: 'traveleranemo',
    element: 'anemo',
    weapon_class: 'Sword',
  },
  travelergeo: {
    key: 'travelergeo',
    element: 'geo',
    weapon_class: 'Sword',
  },
  travelerdendro: {
    key: 'travelerdendro',
    element: 'dendro',
    weapon_class: 'Sword',
  },
  nahida: {
    key: 'nahida',
    element: 'dendro',
    weapon_class: 'catalyst',
  },
  cyno: {
    key: 'cyno',
    element: 'electro',
    weapon_class: 'polearm',
  },
  nilou: {
    key: 'nilou',
    element: 'hydro',
    weapon_class: 'sword',
  },
  alhaitham: {
    key: 'alhaitham',
    element: 'dendro',
    weapon_class: 'sword',
  },
  layla: {
    key: 'layla',
    element: 'cryo',
    weapon_class: 'sword',
  },
  faruzan: {
    key: 'faruzan',
    element: 'anemo',
    weapon_class: 'bow',
  },
  wanderer: {
    key: 'wanderer',
    element: 'anemo',
    weapon_class: 'catalyst',
  },

  dehya: { 
    key: 'dehya', 
    element: 'pyro',
    weapon_type: 'claymore' ,
    },
  
  yaoyao: {
    key: 'yaoyao',
    element: 'dendro',
    weapon_class: 'polearm',
  },
  mika: {
    key: 'mika',
    element: 'cryo',
    weapon_type: 'polearm',
  }
};

export const items: ICharacter[] = Object.keys(characterKeyToICharacter).map(
  (k) => characterKeyToICharacter[k]
);

export const isTraveler = (key: string): boolean =>
  key == 'aether' || key == 'lumine' || key == 'traveler';

export const newChar = (info: ICharacter): Character => {
  let key = info.key;
  if (isTraveler(key) && info.element != 'none')
    key = 'traveler' + info.element;
  //default weapons
  return {
    name: key,
    level: 80,
    max_level: 90,
    element: info.element,
    cons: 0,
    weapon: {
      name: 'dullblade',
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

export const render: ItemRenderer<ICharacter> = (
  item: ICharacter,
  { handleClick, modifiers, query }
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={`${i18n.t('elements.' + item.element)}`}
      key={`${
        item.notes
          ? i18n.t('game:character_names.' + item.key) + ` (${item.notes})`
          : i18n.t('game:character_names.' + item.key)
      }`}
      onClick={handleClick}
      text={highlightText(
        item.notes
          ? i18n.t('game:character_names.' + item.key) + ` (${item.notes})`
          : i18n.t('game:character_names.' + item.key),
        query
      )}
    />
  );
};

// export const render: { [key: string]: ItemRenderer<Character> } = {
//   English: (item: Character, { handleClick, modifiers, query }) => {
//     if (!modifiers.matchesPredicate) {
//       return null;
//     }
//     return (
//       <MenuItem
//         active={modifiers.active}
//         disabled={modifiers.disabled}
//         label={`${
//           item.date_added
//             ? elementRender.English[item.element].concat(
//                 `, Imported: ${item.date_added}`
//               )
//             : elementRender.English[item.element]
//         }`}
//         key={`${
//           item.date_added ? item.name.concat(item.date_added) : item.name
//         }`}
//         onClick={handleClick}
//         text={highlightText(item.name, query)}
//       />
//     );
//   },
//   Chinese: (item: Character, { handleClick, modifiers, query }) => {
//     if (!modifiers.matchesPredicate) {
//       return null;
//     }
//     return (
//       <MenuItem
//         active={modifiers.active}
//         disabled={modifiers.disabled}
//         label={`${
//           item.date_added
//             ? elementRender.Chinese[item.element].concat(
//                 `, Imported: ${item.date_added}`
//               )
//             : elementRender.Chinese[item.element]
//         }`}
//         key={`${
//           item.date_added ? item.name.concat(item.date_added) : item.name
//         }`}
//         onClick={handleClick}
//         text={highlightText(item.name, query)}
//       />
//     );
//   },
// };

// export const filter: ItemPredicate<ICharacter> = (
//   query,
//   item,
//   _index,
//   exactMatch
// ) => {
//   const normalizedQuery = query.toLowerCase();

//   if (exactMatch) {
//     return item.key === normalizedQuery;
//   } else {
//     return (
//       `${item.key} ${item.name} ${item.element}`.indexOf(normalizedQuery) >= 0
//     );
//   }
// };

function escapeRegExpChars(text: string) {
  return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, '\\$1');
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
  const regexp = new RegExp(words.join('|'), 'gi');
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

export const characterSelectProps: {
  itemRenderer: ItemRenderer<ICharacter>;
  items: ICharacter[];
} = {
  itemRenderer: render,
  items: items,
};
