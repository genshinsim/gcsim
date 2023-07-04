import { ElementDPS, FloatStat } from "@gcsim/types";
import { useData } from "@gcsim/ui/src/Pages/Viewer/Components/Damage/CharacterDPSCard";
import { useData as useEleData } from "@gcsim/ui/src/Pages/Viewer/Components/Damage/ElementDPSCard";

import { DataColors, OuterLabelPie } from "@gcsim/ui/src/Pages/Viewer/Components/Util";
import { NoDataIcon } from "@gcsim/ui/src/Pages/Viewer/Components/Util/NoData";
import { ParentSize } from "@visx/responsive";

type Props = {
  dps?: FloatStat[];
}

export const CharacterDPSPie = ({ dps }: Props) => {
  const { data } = useData(dps, ["1", "2", "3", "4"]);

  if (dps == null) {
    return <NoDataIcon className="h-16" />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
        <OuterLabelPie
            width={width}
            height={height}
            data={data}
            pieValue={d => d.pct}
            color={d => DataColors.character(d.index)}
            margin={0}
            pieRadius={0.8}
            outlineWidth={0.5}
        />
      )}
    </ParentSize>
  );
};

type ElementProps = {
  dps?: ElementDPS;
}

export const ElementDPSPie = ({ dps }: ElementProps) => {
  const { data } = useEleData(dps);

  if (dps == null) {
    return <NoDataIcon className="h-16" />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
        <OuterLabelPie
            width={width}
            height={height}
            data={data}
            pieValue={d => d.pct}
            color={d =>  DataColors.element(d.label)}
            margin={0}
            pieRadius={0.8}
            outlineWidth={0.5}
        />
      )}
    </ParentSize>
  );
};