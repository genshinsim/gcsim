import { Omnibar } from "@blueprintjs/select";
import { IAction, IEnemy } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { enemySelectProps } from "./enemies";

const EnemyOmnibar = Omnibar.ofType<IEnemy>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (action: IAction) => void;
};

export function EnemySelect(props: Props) {
  const { t } = useTranslation();
  return (
    <EnemyOmnibar
      resetOnSelect
      {...enemySelectProps}
      initialContent={undefined}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
      className="!absolute !left-0 !right-0 !mx-auto !w-80"
    />
  );
}
