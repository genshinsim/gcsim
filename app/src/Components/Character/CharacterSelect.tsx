import { ItemPredicate, Omnibar } from "@blueprintjs/select";
import { ICharacter, characterSelectProps } from "./characters";
import { simSlice } from "~src/Pages/Sim";
import { RootState, useAppSelector } from "~src/store";
import { Character } from "~src/types";
const CharacterOmnibar = Omnibar.ofType<Character>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: Character) => void;
  disabled?: string[];
};

export function CharacterSelect(props: Props) {
  let disabled: string[] = [];
  if (props.disabled) {
    disabled = props.disabled;
  }

  const { goChars } = useAppSelector((state: RootState) => {
    return {
      goChars: state.sim.GOChars,
    };
  });

  //Forbidden code, this should never reach production
  if (characterSelectProps.items.length < 70) {
    characterSelectProps.items = characterSelectProps.items.concat(goChars);
  }
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
        `${item.name} ${item.date_added} ${item.element}`.indexOf(
          normalizedQuery
        ) >= 0
      );
    }
  };
  console.log("items ", characterSelectProps.items);
  return (
    <CharacterOmnibar
      resetOnSelect
      {...characterSelectProps}
      itemPredicate={filter}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
