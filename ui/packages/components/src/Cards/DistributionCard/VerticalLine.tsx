import { Line } from "@visx/shape";
import { Text } from "@visx/text";

type Props = {
  xScale: (x: number) => number;
  x?: number;
  yMax: number;
  color: string;
  label?: string;
  className?: string;
}

export const VerticalLine = ({ x, xScale, yMax, color, label, className }: Props) => {
  if (x == null) {
    return null;
  }
  
  const localX = xScale(x);
  const Label = ({}) => {
    if (label == null) return null;
    return (
      <Text x={xScale(localX)} y={yMax} dy="1em" className={className} textAnchor="middle">
        {label}
      </Text>
    );
  };
  
  return (
    <g>
      <Line
          from={{ x: localX, y: yMax }}
          to={{ x: localX, y: 0 }}
          stroke={color}
          strokeWidth={2}
          strokeDasharray="5,2"
          className={className} />
      <Label />
    </g>
  );
};