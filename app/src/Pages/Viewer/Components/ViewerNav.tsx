import { Button, ButtonGroup, Intent, Tab, Tabs, Toaster, Icon, Dialog, Classes, Position, Callout, Checkbox } from "@blueprintjs/core";
import classNames from "classnames";
import { Dispatch, SetStateAction, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";
import { updateCfg } from "~src/Pages/Sim";
import { useAppDispatch } from "~src/store";

const btnClass = classNames("hidden ml-[7px] sm:flex");

type NavProps = {
  tabState: [string, Dispatch<SetStateAction<string>>];
  config?: string;
};

export default ({ tabState, config }: NavProps ) => {
  const { t } = useTranslation();
  const [tabId, setTabId] = tabState;

  return (
    <Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
      <Tab id="results" title={t("viewer.results")} className="focus:outline-none" />
      <Tab id="config" title={t("viewer.config")} className="focus:outline-none" />
      <Tab id="analyze" title={t("viewer.analyze")} className="focus:outline-none" />
      <Tab id="debug" title={t("viewer.debug")} className="focus:outline-none" />
      <Tabs.Expander />
      <ButtonGroup>
        <CopyToClipboard config={config} />
        <SendToSim config={config} />
        <Share />
      </ButtonGroup>
    </Tabs>
  );
};

const CopyToClipboard = ({ config }: { config?: string }) => {
  const copyToast = useRef<Toaster>(null);
  const { t } = useTranslation();
  
  const action = () => {
    navigator.clipboard.writeText(config ?? "").then(() => {
      copyToast.current?.show({
        message: t("viewer.copied_to_clipboard"),
        intent: Intent.SUCCESS,
        timeout: 2000
      });
    });
  };

  return (
    <>
      <Button
          icon={<Icon icon="clipboard" className="!mr-0" />}
          onClick={action}
          disabled={config == null}>
        <div className={btnClass}>{t("viewer.copy")}</div>
      </Button>
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
    </>
  );
};

const SendToSim = ({ config }: { config?: string }) => {
  const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";
  const { t } = useTranslation();
  const [, setLocation] = useLocation();
  const dispatch = useAppDispatch();

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
    dispatch(updateCfg(config, keepTeam));
    setLocation("/simulator");
  };

  return (
    <>
      <Button
          className="!hidden sm:!flex"
          icon={<Icon icon="send-to" className="!mr-0" />}
          onClick={() => setOpen(true)}
          disabled={config == null}>
        <div className="hidden ml-[7px] sm:flex">{t("viewer.send_to_simulator")}</div>
      </Button>
      <Dialog
          isOpen={isOpen}
          onClose={() => setOpen(false)}
          title={t("viewer.load_this_configuration")}
          icon="bring-data">
        <div className={Classes.DIALOG_BODY}>
          <Callout intent="warning" className="">
            {t("viewer.this_will_overwrite")}
          </Callout>
          <Checkbox
              label="Copy action list only (ignore character stats)"
              className="my-3 mx-1"
              checked={keepTeam}
              onClick={toggleKeepTeam} />
        </div>
        <div className={classNames(Classes.DIALOG_FOOTER, Classes.DIALOG_FOOTER_ACTIONS)}>
          <Button onClick={toSimulator} intent={Intent.PRIMARY} text={t("viewer.continue")} />
          <Button onClick={() => setOpen(false)} text={t("viewer.cancel")}/>
        </div>
      </Dialog>
    </>
  );
};

const Share = () => {
  const { t } = useTranslation();

  return (
    <Button
        icon={<Icon icon="link" className="!mr-0" />}
        intent={Intent.PRIMARY}
        disabled={true}>
      <div className={btnClass}>{t("viewer.share")}</div>
    </Button>
  );
};