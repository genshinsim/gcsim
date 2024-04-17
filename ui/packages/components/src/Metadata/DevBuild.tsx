import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
  signKey?: string;
  className?: string;
};

export const DevBuild = ({ signKey }: Props) => {
  const { t } = useTranslation();
  if (signKey == null || signKey == "prod") {
    return null;
  }

  if (signKey != "dev") {
    return (
      <Item
        value={t("result.metadata_dev_unofficial") ?? ""}
        intent="danger"
        bright
        bold
      />
    );
  }
  return (
    <Item
      value={t("result.metadata_dev_build") ?? ""}
      intent="danger"
      bright
      bold
    />
  );
};
