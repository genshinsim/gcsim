import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  itr?: number;
};

export const Iterations = ({ itr }: Props) => {
  const { i18n, t } = useTranslation();

  return (
    <Item
      title={t<string>("result.iterations")}
      value={(itr ?? 0).toLocaleString(i18n.language, { notation: "compact" })}
    />
  );
};