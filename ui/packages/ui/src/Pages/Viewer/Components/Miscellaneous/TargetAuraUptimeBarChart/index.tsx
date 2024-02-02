import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useMemo, useState } from "react";
import { CardTitle, useDataColors, useRefreshWithTimer } from "../../Util";
import { BarChart } from "./BarChart";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults | null;
  running: boolean;
};

export const TargetAuraUptimeCard = ({ data, running }: Props) => {
  const { DataColors } = useDataColors();
  const { t } = useTranslation();
  const [stats, timer] = useRefreshWithTimer(
    (d) => {
      return {
        data: d?.statistics?.target_aura_uptime ? d?.statistics?.target_aura_uptime.map(
          (s) => s.sources ? { sources: Object.fromEntries(Object.entries(s.sources).map(([k, v]) => [t<string>("elements."+k), v])) } : {}
        ) : undefined,
      };
    },
    5000,
    data,
    running
  );

  const targets = useMemo(() => {
    if (stats.data == null) {
      return [];
    }

    const targets = new Set<string>();
    for (let i = 0; i < stats.data.length; i++) {
      targets.add(i.toString());
    }
    return Array.from(targets);
  }, [stats.data]);
  const [target, setTarget] = useState("0");

  const auras = useMemo(() => {
    if (stats.data == null) {
      return [];
    }

    const auras = new Set<string>();
    for (const key in stats.data[target]?.sources) {
      auras.add(key);
    }
    return Array.from(auras).sort(
      (a, b) =>
        DataColors.reactableModifierKeys.indexOf(a) -
        DataColors.reactableModifierKeys.indexOf(b)
    );
  }, [stats.data, DataColors.reactableModifierKeys, target]);

  return (
    <Card className="flex flex-col col-span-3 min-h-[384px]">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle
            title={t<string>("result.target_aura_uptime")}
            tooltip="x"
            timer={timer}
          />
          <Options target={target} setTarget={setTarget} targets={targets} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <BarChart
            width={width}
            height={height}
            auraUptime={stats.data}
            auras={auras}
            target={target}
          />
        )}
      </ParentSize>
    </Card>
  );
};

const Options = ({
  target,
  setTarget,
  targets,
}: {
  target: string;
  setTarget: (v: string) => void;
  targets: string[];
}) => {
  const { t } = useTranslation();
  const label = <span className="text-xs font-mono text-gray-400">{t<string>("viewer.target")}</span>;

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={target} onChange={(e) => setTarget(e.target.value)}>
        {targets.map((target) => (
          <option key={target} value={target}>{Number(target) + 1}</option>
        ))}
      </HTMLSelect>
    </FormGroup>
  );
};
