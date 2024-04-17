import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  date?: string;
};

export const DateItem = memo(({ date }: Props) => {
  const { t } = useTranslation();
  const { i18n } = useTranslation();

  if (date == null) {
    return null;
  }

  const d = new Date(Date.parse(date));

  return <Item title={t<string>("result.metadata_created")} value={d.toLocaleDateString(i18n.language)} />;
});