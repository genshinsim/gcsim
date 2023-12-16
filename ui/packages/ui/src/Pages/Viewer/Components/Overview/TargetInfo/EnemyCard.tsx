import { Card } from "@blueprintjs/core";
import { Enemy } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { IconAnemo, IconCryo, IconDendro, IconElectro, IconGeo, IconHydro, IconPhysical, IconPyro } from "../../../../../Components/Icons";
import { DataColors } from "../../Util";

type Props = {
  id: number;
  enemy?: Enemy;
};

export const EnemyCard = (props: Props) => {
  const bgColor = DataColors.qualitative3(props.id);

  return (
    <div className="flex pl-1 min-w-fit" style={{ background: bgColor }}>
      <Card className="flex flex-auto flex-col gap-1">
        <EnemyTitle {...props} />
        <EnemyInfo {...props} />
        <EnemyResistances {...props} />
      </Card>
    </div>
  );
};

const EnemyTitle = ({ id }: Props) => {
  return (
    <div className="flex flex-row items-end gap-3">
      <div className="text-gray-400 text-lg" style={{ color: DataColors.qualitative5(id) }}>
        Target {id+1}
      </div>
    </div>
  );
};

const EnemyInfo = ({ enemy }: Props) => {
  return (
    <div className="flex flex-row font-mono gap-3 h-full items-center">
      <InfoItem name="lvl" value={enemy?.level} />
      <InfoItem name="hp" value={enemy?.hp} />
    </div>
  );
};

const InfoItem = ({ name, value }: { name: string, value?: number }) => {
  const { i18n } = useTranslation();

  if (value == null) {
    return null;
  }

  return (
    <div className="flex flex-row gap-1 text-xs items-center">
      <div className="text-gray-400">{name}</div>
      <div className="font-black text-current text-sm text-bp4-light-gray-500">
        {value?.toLocaleString(i18n.language)}
      </div>
    </div>
  );
};

const EnemyResistances = ({ enemy }: Props) => {
  return (
    <div className="grid grid-cols-4 gap-y-1 text-sm font-mono">
      <Resistance type="anemo" num={enemy?.resist?.["anemo"]} />
      <Resistance type="geo" num={enemy?.resist?.["geo"]} />
      <Resistance type="electro" num={enemy?.resist?.["electro"]} />
      <Resistance type="hydro" num={enemy?.resist?.["hydro"]} />
      <Resistance type="pyro" num={enemy?.resist?.["pyro"]} />
      <Resistance type="cryo" num={enemy?.resist?.["cryo"]} />
      <Resistance type="dendro" num={enemy?.resist?.["dendro"]} />
      <Resistance type="physical" num={enemy?.resist?.["physical"]} />
    </div>
  );
};

const Resistance = ({ type, num }: { type: string, num?: number }) => {
  const { i18n } = useTranslation();
  const format = (val?: number) => val?.toLocaleString(
      i18n.language, { maximumFractionDigits: 2, style: "percent" });

  return (
    <div className="flex flex-row gap-2 items-center">
      <Icon type={type} />
      <div>{format(num ?? 0)}</div>
    </div>
  );
};

const Icon = ({ type }: { type: string }) => {
  const size = "w-[16px] h-[16px] min-w-[16px] min-h-[16px]";
  switch (type) {
    case "electro":
      return <IconElectro className={`${size} text-electro`} />;
    case "pyro":
      return <IconPyro className={`${size} text-pyro`} />;
    case "cryo":
      return <IconCryo className={`${size} text-cryo`} />;
    case "hydro":
      return <IconHydro className={`${size} text-hydro`} />;
    case "geo":
      return <IconGeo className={`${size} text-geo`} />;
    case "anemo":
      return <IconAnemo className={`${size} text-anemo`} />;
    case "physical":
      return <IconPhysical className={`${size}`} />;
    case "dendro":
      return <IconDendro className={`${size} text-dendro`} />;
    default:
      return <IconHydro className={`${size} text-hydro`} />;
  }
};