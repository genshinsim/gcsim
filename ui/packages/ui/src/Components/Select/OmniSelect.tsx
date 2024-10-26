import {MenuItem} from '@blueprintjs/core';
import {ItemPredicate, ItemRenderer, Omnibar} from '@blueprintjs/select';
import i18n from 'i18next';
import React from 'react';
import {CharMap} from '../../Data';

type CharSource = 'user' | 'default';
export interface Item {
  key: string; // this is UI key
  char_key: string; // this is the unique gcsim character key
  char_source: CharSource; // this denotes the source for the char data
  text: string; // main text showing up on search
  label: string; // secondary notes
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
    if (disabled.findIndex((v) => v === item.char_key) > -1) {
      return false;
    }
    const normalizedQuery = query.toLowerCase();
    //the key search should maybe return the result even if not typed in
    //the translated language but wont work for element etc
    return (
      `${item.label} ${item.key} ${item.text}`
        .toLowerCase()
        .indexOf(normalizedQuery) >= 0
    );
  };
  return (
    <CharacterOmnibar
      overlayProps={{usePortal: false}}
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
  {handleClick, modifiers, query},
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={item.label}
      key={item.char_source + '-' + item.key} // it's possible for item.key to be duplicated across 2 diff sources
      onClick={handleClick}
      text={highlightText(item.text, query)}
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
    return {
      key: k,
      char_key: k,
      char_source: 'default',
      text: i18n.t('game:character_names.' + k),
      label: '',
    };
  });
}
