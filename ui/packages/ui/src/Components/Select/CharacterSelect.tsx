import { Omnibar } from "@blueprintjs/select";
import { ICharacter } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { characterSelectProps } from "./characters";

const CharacterOmnibar = Omnibar.ofType<ICharacter>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (character: ICharacter) => void;
};

export function CharacterSelect(props: Props) {
  const { t } = useTranslation();
  return (
    <CharacterOmnibar
      resetOnSelect
      {...characterSelectProps}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
      className="absolute left-0 right-0 mx-auto w-80"
    />
  );
}
