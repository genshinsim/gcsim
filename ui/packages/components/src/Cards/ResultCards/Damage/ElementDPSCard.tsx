import { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
  CardTitle,
  FloatStatTooltipContent,
  NoData,
  OuterLabelPie,
  useDataColors,
  useRefreshWithTimer,
} from "../../../common/gcsim";

type Props = {
  data: model.SimulationResult | null;
  running: boolean;
};

export default ({ data, running }: Props) => {
  const { t } = useTranslation();
  const [dps, timer] = useRefreshWithTimer(
    (d) =>
      d?.statistics?.element_dps
        ? Object.fromEntries(
            Object.entries(d?.statistics?.element_dps).map(([k, v]) => [
              t("elements." + k),
              v,
            ])
          )
        : undefined,
    10000,
    data,
    running
  );

  return (
    <div className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle
        title={t("result.dist", { d: t("result.element_dps") })}
        timer={timer}
      />
      <DPSPie dps={dps} />
    </div>
  );
};

type PieProps = {
  dps?: { [k: string]: model.DescriptiveStats } | null;
};

const DPSPie = memo(({ dps }: PieProps) => {
  const { i18n } = useTranslation();
  const { data } = useData(dps);
  const { DataColors } = useDataColors();

  if (dps == null) {
    return <NoData />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
        <OuterLabelPie
          width={width}
          height={height}
          data={data}
          pieValue={(d) => d.pct}
          color={(d) => DataColors.element(d.label)}
          labelColor={(d) => DataColors.elementLabel(d.label)}
          labelText={(d) => d.label}
          labelValue={(d) => {
            return d.pct.toLocaleString(i18n.language, {
              maximumFractionDigits: 0,
              style: "percent",
            });
          }}
          tooltipContent={(d) => (
            <FloatStatTooltipContent
              title={d.label + " DPS"}
              data={d.value}
              color={DataColors.elementLabel(d.label)}
              percent={d.pct}
            />
          )}
        />
      )}
    </ParentSize>
  );
});

type ElementData = {
  label: string;
  value: model.DescriptiveStats;
  pct: number;
};

export function useData(dps?: { [k: string]: model.DescriptiveStats } | null): {
  data: ElementData[];
  total: number;
} {
  const total = useMemo(() => {
    if (dps == null) {
      return 0;
    }

    let out = 0;
    for (const key in dps) {
      out += dps[key].mean ?? 0;
    }
    return out;
  }, [dps]);

  const data: ElementData[] = useMemo(() => {
    if (dps == null) {
      return [];
    }
    const out: ElementData[] = [];
    for (const key in dps) {
      out.push({
        label: key,
        value: dps[key],
        pct: (dps[key].mean ?? 0) / total,
      });
    }
    return out;
  }, [dps, total]);

  return {
    data,
    total,
  };
}
