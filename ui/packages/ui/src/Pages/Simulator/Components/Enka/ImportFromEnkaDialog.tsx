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
import { userDataActions } from "../../../../Stores/userDataSlice";
import { useAppDispatch } from "../../../../Stores/store";
import { IGOOD } from "../GOOD/GOODTypes";
import { parseFromGOOD } from "../GOOD/parseFromGOOD";
import FetchCharsFromEnka from "./FetchCharsFromEnka";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

const AppToaster = Toaster.create({
  position: Position.BOTTOM,
});

const lsKey = "Enka-UID";

export function ImportFromEnkaDialog(props: Props) {
  const [message, setMessage] = React.useState<string>("");
  const [errors, setErrors] = React.useState<string[]>([]);
  const [uid, setUid] = React.useState<string>("");
  const dispatch = useAppDispatch();

  async function handleClick() {
    localStorage.setItem(lsKey, uid);
    if (uid && validateUid(uid)) {
      try {
        const GOODchars = await FetchCharsFromEnka(uid);
        setErrors(GOODchars.errors ? GOODchars.errors : [])
        console.log(GOODchars);
        const chars = parseFromGOOD(JSON.stringify(GOODchars));

        dispatch(userDataActions.loadFromGOOD({ data: chars.characters }));
        setMessage("success");
      } catch (e) {
        setMessage(`Error importing chars: ${e}`);
      }
    } else {
      setMessage("Invalid UID");
    }
  }

  return (
    <Dialog
      className="w-screen"
      isOpen={props.isOpen}
      onClose={() => {
        props.onClose();
        setMessage("");
      }}
      canEscapeKeyClose
      canOutsideClickClose
      icon="import"
      title="Import from Enka"
      style={{ width: "85%" }}
    >
      <div className={Classes.DIALOG_BODY}>
        <p>
          Ensure your UID has no problems on{" "} 
          <a href="https://enka.network/" target="_blank" rel="noreferrer">
            Enka
          </a>
          .
        </p>
        <Callout intent="warning" title="Warning">
          Importing will replace any existing GOOD/Enka import you already have. This action cannot
          be reversed.
        </Callout>
        <input
          value={uid}
          onChange={(e) => {
            setUid(e.target.value.trim());
          }}
          className="w-full p-2 bg-gray-600 rounded-md mt-2"
          placeholder="Paste UID here"
        />

        {message === "success" ? (
          <>
          <Callout intent="success" className="mt-2 p-2">
            Data retrieved successfully
          </Callout>
          {
            errors.length > 0 ?
            <Callout intent="warning" className="mt-2 p-2">
              Encountered the following issue(s) importing data:
              {errors.map((e,i) => {
                return <div key={i} className="ml-2">{e}</div>
              })}
            </Callout>
            : null
          }
          </>
        ) : (
          <div>
            {message && (
              <Callout intent="warning" className="mt-2 p-2">
                {message}
              </Callout>
            )}
          </div>
        )}

        <p className="font-bold pt-2">
          Once your character data has been imported, you can add your imported character via Add
          Character button and search for the character&apos;s name.
        </p>
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <ButtonGroup>
            <Button onClick={handleClick} intent="primary">
              Fetch
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </Dialog>
  );
}

function hasAlphabet(input: string) {
  return !/^\d+$/.test(input);
}

function validateUid(uid: string) {
  if (uid.length !== 9 || hasAlphabet(uid)) {
    AppToaster.show({
      message: "Invalid UID",
      intent: "danger",
    });
    return false;
  }
  return true;
}
