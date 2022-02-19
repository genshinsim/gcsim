import { SimResults } from "./DataType";
import React, { KeyboardEvent, ClipboardEvent, MouseEvent } from "react";
import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  Position,
  Toast,
  Toaster,
} from "@blueprintjs/core";
import { useAppDispatch } from "~src/store";
import { simActions } from "~src/Pages/Sim";
import { useLocation } from "wouter";

export const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

export function Config({ data }: { data: SimResults }) {
  const [open, setOpen] = React.useState<boolean>(false);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  function copyToClipboard(e: MouseEvent) {
    navigator.clipboard.writeText(data.config_file).then(() => {
      AppToaster.show({ message: "Copied to clipboard", intent: "success" });
    });
    // TODO: Need to add a blueprintjs Toaster for ephemeral confirmation box
  }

  const openInSim = () => {
    setOpen(false);
    dispatch(simActions.setAdvCfg(data.config_file));
    setLocation("/advanced");
  };

  return (
    <div className="flex flex-col">
      {/* <button className="m-2 p-2 rounded-md bg-gray-600" onClick={ copyToClipboard }>Copy Config to Clipboard
        </button> */}
      <div className="ml-2 mr-auto">
        <ButtonGroup>
          <Button onClick={copyToClipboard} icon="clipboard">
            Copy
          </Button>
          <Button onClick={() => setOpen(true)} icon="send-to">
            Send To Simulator
          </Button>
        </ButtonGroup>
      </div>
      <div className="m-2 p-2 rounded-md bg-gray-600">
        <pre className="whitespace-pre-wrap">{data.config_file}</pre>
      </div>
      <Dialog isOpen={open} onClose={() => setOpen(false)}>
        <div className={Classes.DIALOG_BODY}>
          Load this configuration in <span className="font-bold">Advanced</span>{" "}
          mode.
          <Callout intent="warning" className="mt-2">
            This will overwrite any existing configuration you may have. Are you
            sure you wish to continue?
          </Callout>
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={openInSim} intent="primary">
              Continue
            </Button>
            <Button onClick={() => setOpen(false)}>Cancel</Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
