import { Button, Icon, Intent, Toaster } from "@blueprintjs/core";
import { memo, RefObject } from "react";
import { useTranslation } from "react-i18next";

type Props = {
  copyToast: RefObject<Toaster>;
  config?: string;
  className?: string;
}

const CopyTo = ({ copyToast, config, className }: Props) => {
  const { t } = useTranslation();

  const action = () => {
    navigator.clipboard.writeText(config ?? "").then(() => {
      copyToast.current?.show({
        message: t<string>("viewer.copied_to_clipboard"),
        intent: Intent.SUCCESS,
        timeout: 2000,
      });
    });
  };

  return (
    <>
      <Button
          icon={<Icon icon="clipboard" className="!mr-0" />}
          onClick={action}
          disabled={config == null}>
        <div className={className}>{t<string>("viewer.copy")}</div>
      </Button>
    </>
  );
};

export default memo(CopyTo);