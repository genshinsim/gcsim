import { memo } from "react";
import { Item } from "./Item";

type Props = {
  standard?: string;
};

export const Standard = memo(({ standard }: Props) => {
  if (standard == null) {
    return null;
  }

  return (
    <Item title="standard" value={standard} intent="success" bold />
  );
});