import { model } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { Card } from "../../common/ui/card";
import {
  DevBuild,
  Dirty,
  Item,
  Iterations,
  Mode,
  WarningItem,
} from "../../Metadata";

type Props = {
  data: model.SimulationResult;
};

export const Metadata = ({ data }: Props) => {
  if (data.schema_version == null) {
    return (
      <Card className="flex flex-row flex-wrap !p-2 gap-2 justify-center">
        <Item value="legacy sim" intent="danger" bold bright />
      </Card>
    );
  }
  //@ts-ignore: auto generate proto is wrong here. the key is lower case dps, not upper DPS
  let dps: number | undefined = data?.statistics?.dps?.mean;
  let count: number = Object.keys(data?.statistics?.target_dps ?? {}).length;
  if (count > 0 && dps != undefined) {
    dps = dps / (count * 1.0);
  } else {
    dps = undefined;
  }

  return (
    <div className="flex flex-row flex-wrap !p-1.5 gap-2 justify-center bg-slate-700 m-1 -mt-0  rounded-sm border border-gray-600">
      <Error signKey={data.key_type} modified={data.modified} />
      {!data.modified && (data.key_type == null || data.key_type == "prod") ? (
        <DPS dps={dps} />
      ) : null}
      <WarningItem warnings={data?.statistics?.warnings ?? undefined} />
      {/* <Standard standard={data?.standard} /> */}
      <Iterations itr={data?.statistics?.iterations ?? undefined} />
      <Mode mode={data?.mode} />
    </div>
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
      value={(dps ?? 0).toLocaleString(i18n.language, {
        notation: "compact",
        minimumSignificantDigits: 3,
        maximumSignificantDigits: 3,
      })}
    />
  );
};

type ErrorProps = {
  signKey?: string | null | undefined;
  modified?: boolean | null | undefined;
};

export const Error = ({ signKey, modified }: ErrorProps) => {
  if (signKey == null || signKey == undefined || signKey == "prod") {
    return <Dirty modified={modified ?? false} />;
  }

  return <DevBuild signKey={signKey} />;
};
