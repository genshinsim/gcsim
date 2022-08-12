import { Button, Classes, Dialog } from "@blueprintjs/core";
import ReactPlayer from "react-player/file";
import { Trans, useTranslation } from "react-i18next";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  url: string;
};

export function VideoPlayer(props: Props) {
  useTranslation()

  return (
    <Dialog
      style={{ width: "680px" }}
      isOpen={props.isOpen}
      onClose={props.onClose}
      title={props.title}
      icon="help"
    >
      <div className={Classes.DIALOG_BODY}>
        <ReactPlayer controls url={props.url} />
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={props.onClose}><Trans>components.close</Trans></Button>
        </div>
      </div>
    </Dialog>
  );
}
