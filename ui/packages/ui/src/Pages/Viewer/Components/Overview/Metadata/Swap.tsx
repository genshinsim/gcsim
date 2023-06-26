import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  swap?: number;
};

export const Swap = memo(({ swap }: Props ) => {
  const { i18n } = useTranslation();

  if (swap == null) {
    return null;
  }
  return (
    <Item
      title="swap delay"
      value={swap.toLocaleString(i18n.language) + "f"}
      valueCase="lowercase"
    />
  );
});