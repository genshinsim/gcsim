import { ItemPredicate, Omnibar } from "@blueprintjs/select";
import { ICharacter, characterSelectProps } from "./characters";

const CharacterOmnibar = Omnibar.ofType<ICharacter>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: ICharacter) => void;
  disabled?: string[];
};

export function CharacterSelect(props: Props) {
  let disabled: string[] = [];
  if (props.disabled) {
    disabled = props.disabled;
  }

  const filter: ItemPredicate<ICharacter> = (
    query,
    item,
    _index,
    exactMatch
  ) => {
    //ignore filtered items
    if (disabled.findIndex((v) => v === item.key) > -1) {
      return false;
    }

    const normalizedQuery = query.toLowerCase();

    if (exactMatch) {
      return item.key === normalizedQuery;
    } else {
      return (
        `${item.key} ${item.name} ${item.element}`.indexOf(normalizedQuery) >= 0
      );
    }
  };

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
