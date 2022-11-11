import { Button, Classes, Dialog, Icon, InputGroup, Intent, Label, Toaster } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import axios from "axios";
import classNames from "classnames";
import { RefObject, useState } from "react";
import { useTranslation } from "react-i18next";

type ShareProps = {
  running: boolean;
  copyToast: RefObject<Toaster>;
  data: SimResults | null;
  className?: string;
}

// TODO: separate share handling away from the button for caching across pages
export default ({ running, copyToast, data, className }: ShareProps) => {
  const { t } = useTranslation();
  const [isOpen, setOpen] = useState(false);
  const [shareLink, setShareLink] = useState<string | null>(null);

  const handleShare = () => {
    if (data === null) {
      return;
    }

    axios
      .post("/api/share", data)
      .then((resp) => {
        setShareLink(
          window.location.protocol +
            "//" +
            window.location.host +
            "/viewer/share/" +
            resp.data
        );
      })
      .catch((err) => {
        console.log(err);
      });
  };

  const copy = () => {
    navigator.clipboard.writeText(shareLink ?? "").then(() => {
      copyToast.current?.show({
        message: "Link copied to clipboard!",
        intent: Intent.SUCCESS,
        timeout: 2000,
      });
    });
  };

  return (
    <>
      <Button
        icon={<Icon icon="link" className="!mr-0" />}
        intent={Intent.PRIMARY}
        disabled={running || data == null}
        onClick={() => {
          handleShare();
          setOpen(true);
        }}
      >
        <div className={className}>{t<string>("viewer.share")}</div>
      </Button>
      <Dialog
        isOpen={isOpen}
        onClose={() => setOpen(false)}
        title={t<string>("viewer.create_a_shareable")}
        icon="link"
        className="!pb-0"
      >
        <div className={classNames(Classes.DIALOG_BODY, "flex flex-col justify-center gap-2")}>
          <Label>
            Share Link
            <InputGroup
              readOnly={true}
              fill={true}
              onFocus={(e) => {
                e.target.select();
                copy();
              }}
              value={shareLink ?? ""}
              className={classNames({ "bp4-skeleton": shareLink == null })}
              large={true}
              rightElement={<Button icon="duplicate" onClick={() => copy()} />}
            />
          </Label>
        </div>
      </Dialog>
    </>
  );
};
