import { Button, Callout, Checkbox, Classes, Dialog, Icon, Intent } from "@blueprintjs/core";
import classNames from "classnames";
import { memo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

const SendTo = ({ config }: { config?: string }) => {
  const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const history = useHistory();

  const [isOpen, setOpen] = useState(false);
  const [keepTeam, setKeep] = useState<boolean>(() => {
    return localStorage.getItem(LOCALSTORAGE_KEY) === "true";
  });

  const toggleKeepTeam = () => {
    localStorage.setItem(LOCALSTORAGE_KEY, String(!keepTeam));
    setKeep(!keepTeam);
  };

  const toSimulator = () => {
    if (config == null) {
      return;
    }
    dispatch(appActions.setCfg({ cfg: config, keepTeam: keepTeam }));
    history.push("/simulator");
  };

  return (
    <>
      <Button
        icon={<Icon icon="send-to" className="!mr-0" />}
        onClick={() => setOpen(true)}
        disabled={config == null}
      >
        <div className="hidden ml-[7px] sm:flex">{t<string>("viewer.send_to_simulator")}</div>
      </Button>
      <Dialog
        isOpen={isOpen}
        onClose={() => setOpen(false)}
        title={t<string>("viewer.load_this_configuration")}
        icon="bring-data"
      >
        <div className={Classes.DIALOG_BODY}>
          <Callout intent="warning" className="">
            {t<string>("viewer.this_will_overwrite")}
          </Callout>
          <Checkbox
            label={t<string>("viewer.copy_list_only")}
            className="my-3 mx-1"
            checked={keepTeam}
            onClick={toggleKeepTeam}
          />
        </div>
        <div className={classNames(Classes.DIALOG_FOOTER, Classes.DIALOG_FOOTER_ACTIONS)}>
          <Button onClick={toSimulator} intent={Intent.PRIMARY} text={t<string>("viewer.continue")} />
          <Button onClick={() => setOpen(false)} text={t<string>("viewer.cancel")} />
        </div>
      </Dialog>
    </>
  );
};

export default memo(SendTo);