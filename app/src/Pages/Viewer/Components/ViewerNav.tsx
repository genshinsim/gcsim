import { Button, ButtonGroup, Intent, Tab, Tabs, Toaster, Icon } from "@blueprintjs/core";
import { Dispatch, RefObject, SetStateAction } from "react";
import { useTranslation } from "react-i18next";

type NavProps = {
  tabState: [string, Dispatch<SetStateAction<string>>];
  copyToast: RefObject<Toaster>;
  config?: string;
};

export default ({ tabState, copyToast, config }: NavProps ) => {
  const { t } = useTranslation();
  const [tabId, setTabId] = tabState;

  const copyToClipboard = () => {
    navigator.clipboard.writeText(config ?? "").then(() => {
      copyToast.current?.show({
        message: t("viewer.copied_to_clipboard"),
        intent: Intent.SUCCESS,
        timeout: 2000
      });
    });
  };

  return (
    <Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
      <Tab id="results" title={t("viewer.results")} className="focus:outline-none" />
      <Tab id="config" title={t("viewer.config")} className="focus:outline-none" />
      <Tab id="analyze" title={t("viewer.analyze")} className="focus:outline-none" />
      <Tab id="debug" title={t("viewer.debug")} className="focus:outline-none" />
      <Tabs.Expander />
      <ButtonGroup>
        <Button
            icon={<Icon icon="clipboard" className="!mr-0" />}
            onClick={copyToClipboard}
            disabled={config == null}>
          <div className="hidden ml-[7px] sm:flex">{t("viewer.copy")}</div>
        </Button>
        <Button
            className="!hidden sm:!flex"
            icon={<Icon icon="send-to" className="!mr-0" />}
            disabled={config == null}>
          <div className="hidden ml-[7px] sm:flex">{t("viewer.send_to_simulator")}</div>
        </Button>
        <Button
            icon={<Icon icon="link" className="!mr-0" />}
            intent={Intent.PRIMARY}
            disabled={true}>
        <div className="hidden ml-[7px] sm:flex">{t("viewer.share")}</div>
        </Button>
      </ButtonGroup>
    </Tabs>
  );
};