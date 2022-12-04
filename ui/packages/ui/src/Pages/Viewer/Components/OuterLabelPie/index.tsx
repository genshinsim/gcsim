import { Group } from "@visx/group";
import { Pie } from "@visx/shape";
import { OuterLabels } from "./OuterLabels";


type Props<Datum> = {
  width: number;
  height: number;

  data: Datum[];
  pieValue: (d: Datum) => number;
  color: (d: Datum) => string;
  labelColor?: (d: Datum) => string;
  labelText: (d: Datum) => string;
  labelValue: (d: Datum) => string;

  pieRadius?: number;
  labelRadius?: number;
  outline?: string;
  tail?: number;
}

export default <Datum,>(
    {
      width,
      height,
      data,
      color,
      labelColor = color,
      labelText,
      labelValue,
      pieValue,
      pieRadius = 0.6,
      labelRadius = 0.8,
      outline = "#FFF",
      tail = 15
    }: Props<Datum>) => {

  const radius = Math.min(width, height) / 2;

  return (
    <svg width={width} height={height}>
      <Group left={width / 2} top={height / 2}>
        {/* tooltip hover arcs */}

        {/* pie arcs */}
        <Pie
            data={data}
            pieValue={pieValue}
            outerRadius={radius * pieRadius}>
          {(pie) => {
            return pie.arcs.map((arc, index) => {
              return (
                <path
                    key={"arc-" + index}
                    d={pie.path(arc) ?? ""}
                    fill={color(arc.data)}
                    stroke={outline}
                    strokeWidth={1}
                />
              );
            });
          }}
        </Pie>

        {/* label arcs */}
        <Pie
            data={data}
            pieValue={pieValue}
            innerRadius={radius * labelRadius}
            outerRadius={radius * labelRadius}>
          {(pie) => (
            <OuterLabels
                arcs={pie.arcs}
                labelRadius={radius * labelRadius}
                pieRadius={radius * pieRadius}
                labelColor={labelColor}
                labelText={labelText}
                labelValue={labelValue}
                tail={tail} />
          )}
        </Pie>
      </Group>
    </svg>
  );
};