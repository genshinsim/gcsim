import { Warnings } from "@gcsim/types";
import { memo } from "react";
import { Item } from "./Item";

type Props = {
  warnings?: Warnings;
};

export const WarningItem = memo(({ warnings }: Props ) => {
  if (warnings == null) {
    return null;
  }
  const count = Object.entries(warnings).filter(([, v]) => (v as boolean)).length;
  if (count == 0) {
    return null;
  }

  return (
    <Item
      title="warnings"
      value={count.toLocaleString()}
      intent="warning"
      bold
    />
  );
});