import { Button, Classes, Dialog } from "@blueprintjs/core";
import { eventColor } from "./parse";
import { Trans, useTranslation } from "react-i18next";

export interface OptionsProp {
  isOpen: boolean;
  handleClose: () => void;
  handleToggle: (opt: string) => void;
  handleClear: () => void;
  handleResetDefault: () => void;
  handleSetPresets: (opt: "simple" | "advanced" | "verbose" | "debug") => void;
  selected: string[];
  options: string[];
}

export function Options(props: OptionsProp) {
  const { t } = useTranslation();

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
          <div className="text-md font-medium">
            <Trans>viewer.log_options</Trans>
          </div>
          <div className="grid grid-cols-3">{cols}</div>
          <div>{/* <ButtonGroup></ButtonGroup> */}</div>
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={() => props.handleSetPresets("simple")}>
              {t<string>("viewer.simple")}
            </Button>
            <Button onClick={() => props.handleSetPresets("advanced")}>
              {t<string>("viewer.advanced")}
            </Button>
            <Button onClick={() => props.handleSetPresets("verbose")}>
              {t<string>("viewer.verbose")}
            </Button>
            <Button onClick={() => props.handleSetPresets("debug")}>
              {t<string>("viewer.debug")}
              </Button>
            <Button intent="danger" onClick={props.handleClear}>
              {t<string>("viewer.clear")}
            </Button>
            <Button intent="none" onClick={props.handleClose}>
              {t<string>("viewer.close")}
            </Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}
