import { ItemPredicate, Omnibar } from "@blueprintjs/select";
import { ICharacter, characterSelectProps, elementRender } from "./characters";
import { useTranslation } from 'react-i18next'

const CharacterOmnibar = Omnibar.ofType<ICharacter>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: ICharacter) => void;
  disabled?: string[];
};

export function CharacterSelect(props: Props) {
  let { i18n } = useTranslation()

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
        `${item.key} ${item.name} ${elementRender[i18n.language][item.element]}`.indexOf(normalizedQuery) >= 0
      );
    }
  };

  return (
    <CharacterOmnibar
      resetOnSelect
      {...characterSelectProps[i18n.language]}
      itemPredicate={filter}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
