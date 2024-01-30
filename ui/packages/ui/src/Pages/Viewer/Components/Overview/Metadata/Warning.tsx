import { Warnings } from "@gcsim/types";
import { memo } from "react";
import { Item } from "./Item";
import { useTranslation } from "react-i18next";

type Props = {
  warnings?: Warnings;
};

export const WarningItem = memo(({ warnings }: Props ) => {
  const { t } = useTranslation();
  if (warnings == null) {
    return null;
  }
  const count = Object.entries(warnings).filter(([, v]) => (v as boolean)).length;
  if (count == 0) {
    return null;
  }

  return (
    <Item
      title={t<string>("result.metadata_warnings")}
      value={count.toLocaleString()}
      intent="warning"
      bold
    />
  );
});