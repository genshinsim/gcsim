import { Button, Dialog, Icon } from "@blueprintjs/core";
import React from "react";
import { DebugItem } from "./parse";

export function DebugItemView({
  item,
  showBuffDuration,
}: {
  item: DebugItem;
  showBuffDuration: (e: DebugItem) => void;
}) {
  const [open, setOpen] = React.useState<boolean>(false);
  const handleClick = () => {
    setOpen(true);
  };
  return (
    <div
      className="flex flex-row gap-2 items-center pl-1 pr-1 pt-px pb-px rounded-md m-1 "
      style={{ backgroundColor: item.color }}
    >
      <span
        className="material-icons text-sm cursor-pointer"
        onClick={() => showBuffDuration(item)}
      >
        {item.icon}
      </span>
      <div className="flex-grow cursor-pointer" onClick={handleClick}>
        {item.msg}
      </div>
      <div>{item.target}</div>
      <Dialog
        canEscapeKeyClose
        canOutsideClickClose
        autoFocus
        enforceFocus
        shouldReturnFocusOnClose
        isOpen={open}
        onClose={() => {
          setOpen(false);
        }}
      >
        <pre className="m-2 whitespace-pre-wrap">{item.raw}</pre>
      </Dialog>
    </div>
  );
}
