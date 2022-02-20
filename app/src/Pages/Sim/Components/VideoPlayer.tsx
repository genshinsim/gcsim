import { Button, Callout, Classes, Dialog } from "@blueprintjs/core";
import ReactPlayer from "react-player/file";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  url: string;
};

export function VideoPlayer(props: Props) {
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
          <Button onClick={props.onClose}>Close</Button>
        </div>
      </div>
    </Dialog>
  );
}
