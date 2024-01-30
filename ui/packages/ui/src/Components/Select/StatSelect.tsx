import { Omnibar } from "@blueprintjs/select";
import { IStat } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { statSelectProps } from "./stats";

const StatOmnibar = Omnibar.ofType<IStat>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (stat: IStat) => void;
};

export function StatSelect(props: Props) {
  const { t } = useTranslation();
  return (
    <StatOmnibar
      resetOnSelect
      {...statSelectProps}
      initialContent={undefined}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
      className="!absolute !left-0 !right-0 !mx-auto !w-80"
    />
  );
}
