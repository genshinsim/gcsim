import { Button, Classes, Dialog } from "@blueprintjs/core";
import React from "react";

type Props = {
  open: boolean;
  onClose: () => void;
};
export function LoadGOOD(props: Props) {
  const [str, setStr] = React.useState<string>("");
  const handleLoad = () => {};
  return (
    <Dialog
      isOpen={props.open}
      onClose={props.onClose}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title="Import from GOOD"
    >
      <div className={Classes.DIALOG_BODY}></div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={handleLoad}>Load</Button>
          <Button onClick={props.onClose}>Close</Button>
        </div>
      </div>
    </Dialog>
  );
}
