import { memo } from "react";
import { Item } from "./Item";

type Props = {
  mode?: number;
};

export const ModeItem = memo(({ mode }: Props) => {
  if (mode == null) {
    return null;
  }
  
  const modeName = mode == 2 ? "ttk" : "duration";
  return <Item title="mode" value={modeName} />;
});