import { SimResults } from "@gcsim/types";
import DPS from "./DPS";
import EPS from "./EPS";
import HPS from "./HPS";
import RPS from "./RPS";
import SHP from "./SHP";
import SimDur from "./SimDur";

// Qualitative colors from https://blueprintjs.com/docs/#core/colors
const colors = [
  "#D33D17",
  "#147EB3",
  "#9D3F9D",
  "#29A634",
  "#D1980B",
  "#00A396",
  "#DB2C6F",
  "#8EB125",
  "#946638",
  "#7961DB",
];

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  const cards = [DPS, EPS, RPS, HPS, SHP, SimDur];

  return (
    <div className="col-span-full flex flex-row flex-wrap gap-2 justify-center">
      {cards.map((e, i) => e({ data: data, color: colors[i] }))}
    </div>
  );
};
