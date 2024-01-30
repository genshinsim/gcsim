import { memo } from "react";
import { Item } from "./Item";
import { useTranslation } from "react-i18next";

type Props = {
  mode?: number;
};

export const ModeItem = memo(({ mode }: Props) => {
  const { t } = useTranslation();
  if (mode == null) {
    return null;
  }
  
  const modeName = mode == 2 ? t<string>("db.ttk") : t<string>("db.duration");
  return <Item title={t<string>("db.simMode")} value={modeName} valueCase="lowerCase" />;
});