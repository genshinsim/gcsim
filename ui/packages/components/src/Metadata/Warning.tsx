import { model } from "@gcsim/types";
import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  warnings?: model.Warnings;
};

export const WarningItem = memo(({ warnings }: Props) => {
  const { t } = useTranslation();
  if (warnings == null) {
    return null;
  }
  const count = Object.entries(warnings).filter(([, v]) => v as boolean).length;
  if (count == 0) {
    return null;
  }

  return (
    <Item
      title={t("result.metadata_warnings")}
      value={count.toLocaleString()}
      intent="warning"
      bold
    />
  );
});
