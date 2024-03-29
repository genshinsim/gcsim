import { memo } from "react";
import { Item } from "./Item";
import { useTranslation } from "react-i18next";

type Props = {
  standard?: string;
};

export const Standard = memo(({ standard }: Props) => {
  const { t } = useTranslation();
  if (standard == null) {
    return null;
  }

  return (
    <Item title={t<string>("result.metadata_standard")} value={standard} intent="success" bold />
  );
});