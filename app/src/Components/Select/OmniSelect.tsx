import { ItemPredicate, ItemRenderer, Omnibar } from '@blueprintjs/select';
import { MenuItem } from '@blueprintjs/core';
import { CharMap, TransformTravelerKeyToName, TravelerCheck } from '~src/Data';
import i18n from 'i18next';

export interface Item {
  key: string;
  text: string;
  label: string;
  notes?: string;
}

type Props = {
  isOpen: boolean;
  items: Item[];
  onClose: () => void;
  onSelect: (item: Item) => void;
  disabled?: string[];
};

const CharacterOmnibar = Omnibar.ofType<Item>();

export function OmniSelect(props: Props) {
  let disabled: string[] = [];
  if (props.disabled) {
    disabled = props.disabled;
  }

  const filter: ItemPredicate<Item> = (query, item, _index) => {
    //ignore filtered items
    if (disabled.findIndex((v) => v === item.key) > -1) {
      return false;
    }
    const normalizedQuery = query.toLowerCase();
    //the key search should maybe return the result even if not typed in
    //the translated language but wont work for element etc
    return (
      `${item.label} ${item.key} ${item.text} ${item.notes}`
        .toLowerCase()
        .indexOf(normalizedQuery) >= 0
    );
  };
  return (
    <CharacterOmnibar
      overlayProps={{ usePortal: false }}
      resetOnSelect
      items={props.items}
      itemRenderer={OmniSelectRenderer}
      itemPredicate={filter}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}

const OmniSelectRenderer: ItemRenderer<Item> = (
  item: Item,
  { handleClick, modifiers, query }
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  const text = item.notes ? `${item.text} (${item.notes})` : item.text;
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={item.label}
      key={text}
      onClick={handleClick}
      text={highlightText(text, query)}
    />
  );
};

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

export function GenerateDefaultCharacters(): Item[] {
  return Object.keys(CharMap).map((k) => {
    const ele = i18n.t(`elements.${CharMap[k].element}`);
    let extra = '';
    if (TravelerCheck(k)) {
      extra = ` (${ele})`;
    }
    return {
      key: k,
      text:
        i18n.t('game:character_names.' + TransformTravelerKeyToName(k)) + extra,
      label: i18n.t(`elements.${CharMap[k].element}`),
    };
  });
}
