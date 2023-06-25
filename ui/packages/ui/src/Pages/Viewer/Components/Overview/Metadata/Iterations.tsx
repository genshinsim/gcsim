import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  itr?: number;
};

export const Iterations = ({ itr }: Props) => {
  const {i18n} = useTranslation();

  return (
    <Item
      title="iterations"
      value={(itr ?? 0).toLocaleString(i18n.language, { notation: "compact" })}
    />
  );
};