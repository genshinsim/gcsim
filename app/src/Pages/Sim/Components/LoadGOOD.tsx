import { Button, Classes, Dialog } from "@blueprintjs/core";
import React from "react";
import { IGOODImport, staticPath, parseFromGO } from "./char";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const lsKey = "GOOD-import";

export function LoadGOOD(props: Props) {
  const [str, setStr] = React.useState<string>("");
  React.useEffect(() => {
    const val = localStorage.getItem(lsKey);
    if (val !== null && val !== "") {
      setStr(val);
    }
  }, []);
  const handleLoad = () => {};
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setStr(e.target.value);
    localStorage.setItem(lsKey, e.target.value);
    console.log(parseFromGO(e.target.value));
  };
  return (
    <Dialog
      isOpen={props.isOpen}
      onClose={props.onClose}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title="Import from Genshin Optimizer/GOOD"
    >
      <div className={Classes.DIALOG_BODY}>
        <p>
          Paste import data in GOOD format in the textbox below. (If you are
          coming from Genshin Optimizer, you can export your database in GOOD
          format
          <a
            href="https://frzyc.github.io/genshin-optimizer/#/database"
            target="_blank"
          >
            here
          </a>
        </p>
        <textarea
          value={str}
          onChange={handleChange}
          className="w-full p-2 bg-gray-600 rounded-md mt-2"
          rows={7}
        />
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={handleLoad}>Load</Button>
          <Button onClick={props.onClose}>Close</Button>
        </div>
      </div>
    </Dialog>
  );
}
