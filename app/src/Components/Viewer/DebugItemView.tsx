import { Dialog } from "@blueprintjs/core";
import React from "react";
import { DebugItem } from "./parse";

export function DebugItemView({ item }: { item: DebugItem }) {
  const [open, setOpen] = React.useState<boolean>(false);
  const handleClick = () => {
    setOpen(true);
  };
  return (
    <div
      className="flex flex-row gap-2 items-center pl-1 pr-1 pt-px pb-px rounded-md m-1 cursor-pointer"
      style={{ backgroundColor: item.color }}
      onClick={handleClick}
    >
      <span className="material-icons text-sm">{item.icon}</span>
      <div className="flex-grow">{item.msg}</div>
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
