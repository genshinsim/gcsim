import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  Position,
  Toaster,
} from "@blueprintjs/core";
import React from "react";
import { useTranslation } from "react-i18next";
import { useAppDispatch } from "~src/store";
import { userDataActions } from "../../Pages/Sim/userDataSlice";
import { IGOODImport, parseFromGOOD } from "./parseFromGOOD";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

const lsKey = "GOOD-import";

export function ImportFromGOODDialog(props: Props) {
  const [data, setData] = React.useState<IGOODImport>();
  const dispatch = useAppDispatch();
  let { t } = useTranslation();

  const handleLoad = () => {
    if (data !== undefined) {
      dispatch(userDataActions.loadFromGOOD({ data: data.characters }));
      props.onClose();
      AppToaster.show({
        message: t("importer.import_success"),
        intent: "success",
      });
    }
  };
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    localStorage.setItem(lsKey, e.target.value);
    setData(parseFromGOOD(e.target.value));
  };
  return (
    <Dialog
      className="w-screen"
      isOpen={props.isOpen}
      onClose={props.onClose}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title="Import from Genshin Optimizer/GOOD"
      style={{ width: "85%" }}
    >
      <div className={Classes.DIALOG_BODY}>
        <p>
          Paste import data in GOOD format in the textbox below. (If you are
          coming from Genshin Optimizer, you can export your database in GOOD
          format{" "}
          <a
            href="https://frzyc.github.io/genshin-optimizer/#/database"
            target="_blank"
          >
            here
          </a>
        </p>
        <Callout intent="warning" title="Warning">
          Importing will replace any existing GOOD import you already have. This
          action cannot be reversed.
        </Callout>
        <textarea
          value={localStorage.getItem(lsKey) ?? ""}
          onChange={handleChange}
          className="w-full p-2 bg-gray-600 rounded-md mt-2"
          rows={7}
        />
        <p className="font-bold">
          Once your character data has been imported, you can add your imported
          character via Add Character button and search for the character's
          name.
        </p>
        {data ? (
          data.err === "" ? (
            <Callout intent="success" className="mt-2 p-2">
              Data parsed successfully
            </Callout>
          ) : (
            <Callout intent="warning" className="mt-2 p-2">
              {data!.err}
            </Callout>
          )
        ) : null}
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <ButtonGroup>
            <Button
              onClick={handleLoad}
              disabled={!data || data.err !== ""}
              intent="primary"
            >
              Load
            </Button>
            <Button onClick={props.onClose} intent="danger">
              Cancel
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </Dialog>
  );
}
