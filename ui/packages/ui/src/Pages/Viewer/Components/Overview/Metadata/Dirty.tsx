import { memo } from "react";
import { Item } from "./Item";
import { useTranslation } from "react-i18next";

type Props = {
  modified?: boolean;
};

export const Dirty = memo(({ modified }: Props) => {
  const { t } = useTranslation();
  if (!modified) {
    return null;
  }
  return <Item value={t<string>("result.metadata_dirty")} intent="danger" bright bold />;
});