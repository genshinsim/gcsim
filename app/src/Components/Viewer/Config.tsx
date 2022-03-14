import { SimResults } from "./DataType";
import React, { MouseEvent } from "react";
import { Button, ButtonGroup, Callout, Classes, Dialog, Position, Toaster, } from "@blueprintjs/core";
import { useAppDispatch } from "~src/store";
import { simActions } from "~src/Pages/Sim";
import { useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";

export const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

export function Config({ data }: { data: SimResults }) {
  let { t } = useTranslation()

  const [open, setOpen] = React.useState<boolean>(false);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  function copyToClipboard(e: MouseEvent) {
    navigator.clipboard.writeText(data.config_file).then(() => {
      AppToaster.show({ message: t("viewer.copied_to_clipboard"), intent: "success" });
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
            <Trans>viewer.copy</Trans>
          </Button>
          <Button onClick={() => setOpen(true)} icon="send-to">
            <Trans>viewer.send_to_simulator</Trans>
          </Button>
        </ButtonGroup>
      </div>
      <div className="m-2 p-2 rounded-md bg-gray-600">
        <pre className="whitespace-pre-wrap">{data.config_file}</pre>
      </div>
      <Dialog isOpen={open} onClose={() => setOpen(false)}>
        <div className={Classes.DIALOG_BODY}>
          <Trans>viewer.load_this_configuration</Trans>
          <Callout intent="warning" className="mt-2">
            <Trans>viewer.this_will_overwrite</Trans>
          </Callout>
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={openInSim} intent="primary">
              <Trans>viewer.continue</Trans>
            </Button>
            <Button onClick={() => setOpen(false)}><Trans>viewer.cancel</Trans></Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
