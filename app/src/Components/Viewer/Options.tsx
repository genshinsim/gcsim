import { Button, Classes, Dialog } from "@blueprintjs/core";
import { eventColor } from "./parse";
import { Trans, useTranslation } from "react-i18next";

export interface OptionsProp {
  isOpen: boolean;
  handleClose: () => void;
  handleToggle: (opt: string) => void;
  handleClear: () => void;
  handleResetDefault: () => void;
  selected: string[];
  options: string[];
}

export function Options(props: OptionsProp) {
  useTranslation()

  const cols = props.options.map((o, index) => {
    return (
      <div className="flex flex-row gap-1 p-1 items-center" key={index}>
        <label className="cursor-pointer">
          <input
            type="checkbox"
            checked={props.selected.indexOf(o) > -1}
            className="checkbox cursor-pointer"
            onChange={() => props.handleToggle(o)}
          />
          <span className="font-medium text-sm pl-1" style={{ color: eventColor(o) }}>
          {o}
        </span>
        </label>
      </div>
    );
  });

  return (
    <Dialog
      canEscapeKeyClose
      canOutsideClickClose
      autoFocus
      enforceFocus
      shouldReturnFocusOnClose
      isOpen={props.isOpen}
      onClose={props.handleClose}
    >
      <div className="p-2">
        <div className={Classes.DIALOG_BODY}>
          <div className="text-md font-medium"><Trans>viewer.log_options</Trans></div>
          <div className="grid grid-cols-3">{cols}</div>
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button intent="primary" onClick={props.handleResetDefault}>
              <Trans>viewer.defaults</Trans>
            </Button>
            <Button intent="danger" onClick={props.handleClear}>
              <Trans>viewer.clear</Trans>
            </Button>
            <Button intent="none" onClick={props.handleClose}>
              <Trans>viewer.close</Trans>
            </Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}
