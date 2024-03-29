import { model } from "@gcsim/types";
import { useData } from "../../ResultCards/Damage/CharacterDPSCard";
import { useData as useEleData } from "../../ResultCards/Damage/ElementDPSCard";

import { ParentSize } from "@visx/responsive";
import {
  NoDataIcon,
  OuterLabelPie,
  useDataColors,
} from "../../../common/gcsim";

type Props = {
  dps?: model.DescriptiveStats[] | null;
};

export const CharacterDPSPie = ({ dps }: Props) => {
  const { DataColors } = useDataColors();
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
          pieValue={(d) => d.pct}
          color={(d) => DataColors.character(d.index)}
          margin={0}
          pieRadius={0.8}
          outlineWidth={0.5}
        />
      )}
    </ParentSize>
  );
};

type ElementProps = {
  dps?: { [k: string]: model.DescriptiveStats } | null;
};

export const ElementDPSPie = ({ dps }: ElementProps) => {
  const { DataColors } = useDataColors();
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
          pieValue={(d) => d.pct}
          color={(d) => DataColors.element(d.label)}
          margin={0}
          pieRadius={0.8}
          outlineWidth={0.5}
        />
      )}
    </ParentSize>
  );
};
