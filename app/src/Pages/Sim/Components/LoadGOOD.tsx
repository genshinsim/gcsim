import { Button, Classes, Dialog } from "@blueprintjs/core";
import React from "react";
import { useAppDispatch } from "~src/store";
import { userDataActions } from "../userDataSlice";
import { IGOODImport, parseFromGO } from "./Import";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const lsKey = "GOOD-import";

export function LoadGOOD(props: Props) {
  const [str, setStr] = React.useState<string>("");
  const [data, setData] = React.useState<IGOODImport>();
  const [isSuccess, setIsSuccess] = React.useState(false);
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    const val = localStorage.getItem(lsKey);
    if (val !== null && val !== "") {
      setStr(val);
      // console.log("spahget", parseFromGO(val));
      setData(parseFromGO(val));
    }
  }, []);
  const handleLoad = () => {
    if (data !== undefined) {
      // setData(parseFromGO());
      dispatch(userDataActions.loadFromGOOD({ data: data.characters }));
      setIsSuccess(true);
    }
  };
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setStr(e.target.value);
    localStorage.setItem(lsKey, e.target.value);
    setData(parseFromGO(e.target.value));
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
            <text> here</text>
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
          {data && data.err && (
            <div className="pt-1.5 text-red-500">{data.err}</div>
          )}
          {data && isSuccess && data.err === "" && (
            <div className="pt-1.5 text-green-500">Successfully imported!</div>
          )}

          <Button onClick={handleLoad}>Load</Button>
          <Button onClick={props.onClose}>Close</Button>
        </div>
      </div>
    </Dialog>
  );
}
