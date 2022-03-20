import { ItemPredicate, Omnibar } from "@blueprintjs/select";
import { characterSelectProps, ICharacter } from "./characters";
const CharacterOmnibar = Omnibar.ofType<ICharacter>();
import { useTranslation } from "react-i18next";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: ICharacter) => void;
  additionalOptions?: ICharacter[];
  disabled?: string[];
};

export function CharacterSelect(props: Props) {
  let { t } = useTranslation();

  let disabled: string[] = [];
  if (props.disabled) {
    disabled = props.disabled;
  }

  console.log(characterSelectProps);
  let items = [...characterSelectProps.items];
  console.log("before additional", items);
  if (props.additionalOptions) {
    items = items.concat(props.additionalOptions);
  }
  console.log("after additional", items);

  const filter: ItemPredicate<ICharacter> = (
    query,
    item,
    _index,
    exactMatch
  ) => {
    //ignore filtered items
    if (
      disabled.findIndex((v) => v === t("character_names." + item.key)) > -1
    ) {
      return false;
    }

    const normalizedQuery = query.toLowerCase();
    const transChar = t("character_names." + item.key)
      .replace(" ", "")
      .toLowerCase();
    if (exactMatch) {
      return t("character_names." + item.key) === normalizedQuery;
    } else {
      return (
        `${transChar} ${item.key} ${item.notes} ${t(
          "elements." + item.element
        )}`.indexOf(normalizedQuery) >= 0
      );
    }
  };
  return (
    <CharacterOmnibar
      resetOnSelect
      items={items}
      itemRenderer={characterSelectProps.itemRenderer}
      itemPredicate={filter}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
