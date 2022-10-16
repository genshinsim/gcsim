import { Button, ButtonGroup, Intent, Tab, Tabs, Toaster, Position, Icon } from "@blueprintjs/core";
import classNames from "classnames";
import { Dispatch, SetStateAction } from "react";
import { useTranslation } from "react-i18next";

// TODO: shared toaster in Viewer?
const NavToaster = Toaster.create({
  position: Position.TOP
});

type NavProps = {
  tabState: [string, Dispatch<SetStateAction<string>>];
  config: string | undefined;
};

export default ({ tabState, config }: NavProps ) => {
  const { t } = useTranslation();
  const [tabId, setTabId] = tabState;

  const tabClass = classNames("focus:outline-none");

  const copyToClipboard = () => {
    navigator.clipboard.writeText(config ?? "").then(() => {
      NavToaster.show({
        message: t("viewer.copied_to_clipboard"),
        intent: Intent.SUCCESS
      });
    });
  };

  return (
    <Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
      <Tab id="results" title={t("viewer.results")} className={tabClass} />
      <Tab id="config" title={t("viewer.config")} className={tabClass} />
      <Tab id="analyze" title={t("viewer.analyze")} className={tabClass} />
      <Tab id="debug" title={t("viewer.debug")} className={tabClass} />
      <Tabs.Expander />
      <ButtonGroup>
        <Button
            icon={<Icon icon="clipboard" className="!mr-0" />}
            onClick={copyToClipboard}>
          <div className="hidden ml-[7px] sm:flex">{t("viewer.copy")}</div>
        </Button>
        <Button
            className="!hidden sm:!flex"
            icon={<Icon icon="send-to" className="!mr-0" />}>
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