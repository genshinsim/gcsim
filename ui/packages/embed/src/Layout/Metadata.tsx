import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { DevBuild } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/DevBuild";
import { Dirty } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Dirty";
import { Item } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Item";
import { Iterations } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Iterations";
import { ModeItem } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Mode";
import { Standard } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Standard";
import { WarningItem } from "@gcsim/ui/src/Pages/Viewer/Components/Overview/Metadata/Warning";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults;
};

export const Metadata = ({ data }: Props) => {
  if (data.schema_version == null) {
    return (
      <Card className="flex flex-row flex-wrap !p-2 gap-2 justify-center">
        <Item value="legacy sim" intent="danger" bold bright />
      </Card>
    );
  }
  let dps: number | undefined = data?.statistics?.dps?.mean;
  let count: number = Object.keys(data?.statistics?.target_dps ?? {}).length;
  if (count > 0 && dps != undefined) {
    dps = dps / (count * 1.0);
  } else {
    dps = undefined;
  }

  return (
    <Card className="flex flex-row flex-wrap !p-2 gap-2 justify-center">
      <Error signKey={data.key_type} modified={data.modified} />
      {!data.modified && (data.key_type == null || data.key_type == "prod") ? (
        <DPS dps={dps} />
      ) : null}
      <WarningItem warnings={data?.statistics?.warnings} />
      <Standard standard={data?.standard} />
      <Iterations itr={data?.statistics?.iterations} />
      <ModeItem mode={data?.mode} />
    </Card>
  );
};

type DPSProps = {
  dps?: number;
};

export const DPS = ({ dps }: DPSProps) => {
  const { i18n } = useTranslation();

  if (dps == undefined) {
    <Item title="dps/target" value={"n/a"} />;
  }

  return (
    <Item
      title="dps/target"
      value={(dps ?? 0).toLocaleString(i18n.language, { notation: "compact" })}
    />
  );
};

type ErrorProps = {
  signKey?: string;
  modified?: boolean;
};

export const Error = ({ signKey, modified }: ErrorProps) => {
  if (signKey == null || signKey == "prod") {
    return <Dirty modified={modified} />;
  }

  return <DevBuild signKey={signKey} />;
};
