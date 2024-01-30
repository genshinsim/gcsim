import { Omnibar } from "@blueprintjs/select";
import { IAction } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { actionSelectProps } from "./actions";

const ActionOmnibar = Omnibar.ofType<IAction>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (action: IAction) => void;
};

export function ActionSelect(props: Props) {
  const { t } = useTranslation();
  return (
    <ActionOmnibar
      resetOnSelect
      {...actionSelectProps}
      initialContent={undefined}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
      className="!absolute !left-0 !right-0 !mx-auto !w-80"
    />
  );
}
