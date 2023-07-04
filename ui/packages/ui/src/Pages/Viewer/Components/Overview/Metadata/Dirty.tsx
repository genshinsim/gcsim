import { memo } from "react";
import { Item } from "./Item";

type Props = {
  modified?: boolean;
};

export const Dirty = memo(({ modified }: Props) => {
  if (!modified) {
    return null;
  }
  return <Item value="dirty" intent="danger" bright bold />;
});