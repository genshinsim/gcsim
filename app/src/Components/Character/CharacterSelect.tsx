import { ItemPredicate, Omnibar } from "@blueprintjs/select";
import { characterSelectProps,elementRender } from "./characters";
import { RootState, useAppSelector } from "~src/store";
import { Character } from "~src/types";
const CharacterOmnibar = Omnibar.ofType<Character>();
import { useTranslation } from 'react-i18next'

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: Character) => void;
  disabled?: string[];
};

export function CharacterSelect(props: Props) {
  let { i18n } = useTranslation()

  let disabled: string[] = [];
  if (props.disabled) {
    disabled = props.disabled;
  }

  const { goChars } = useAppSelector((state: RootState) => {
    return {
      goChars: state.sim.GOChars,
    };
  });
  const loadedCharacterSelectProps = characterSelectProps[i18n.language]
  const items = loadedCharacterSelectProps.items.concat(goChars);

  const filter: ItemPredicate<Character> = (
    query,
    item,
    _index,
    exactMatch
  ) => {
    //ignore filtered items
    if (disabled.findIndex((v) => v === item.name) > -1) {
      return false;
    }

    const normalizedQuery = query.toLowerCase();

    if (exactMatch) {
      return item.name === normalizedQuery;
    } else {
      return (
        `${item.name} ${item.date_added} ${elementRender[i18n.language][item.element]}`.indexOf(
          normalizedQuery
        ) >= 0
      );
    }
  };
  return (
    <CharacterOmnibar
      resetOnSelect
      items={items}
      itemRenderer={loadedCharacterSelectProps.itemRenderer}
      itemPredicate={filter}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
