import { curveBasis } from "@visx/curve";
import { Group } from "@visx/group";
import { LinePath } from "@visx/shape";
import { PieArcDatum } from "@visx/shape/lib/shapes/Pie";
import { useMemo } from "react";

type Props<Datum> = {
  arcs: PieArcDatum<Datum>[];
  labelRadius: number;
  pieRadius: number;

  labelColor: (d: Datum) => string;
  labelText: (d: Datum) => string;
  labelValue: (d: Datum) => string;

  xPadding?: number;
  yPadding?: number;
  tail?: number;

  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, index: number, data: Datum) => void;
};

type Point = {
  x: number;
  y: number;
};

type LabelPosition = {
  index: number;
  angle: number;
  x: number;
  y: number;
  value: number;
};

export const OuterLabels = <Datum,>({
  arcs,
  labelRadius,
  pieRadius,
  labelColor,
  labelText,
  labelValue,

  mouseHover,
  mouseLeave,

  xPadding = 8,
  yPadding = 18,
  tail = 15,
}: Props<Datum>) => {
  const labelPositions = useLabelPositions(
    arcs,
    labelRadius,
    xPadding,
    yPadding
  );
  const linePoints = useLinePoints(labelPositions, pieRadius, tail);

  if (labelRadius == 0) {
    return null;
  }

  return (
    <>
      {arcs.map((arc, index) => {
        const left = midAngle(arc) > Math.PI;
        return (
          <Group key={"g-" + index}>
            <Group
              key={"label-" + index}
              left={labelPositions.get(index)?.x}
              top={labelPositions.get(index)?.y}
              onMouseMove={(e) => mouseHover(e, index, arc.data)}
              onMouseLeave={() => mouseLeave()}
            >
              <text
                dx={left ? "-.5em" : ".5em"}
                textAnchor={left ? "end" : "start"}
                className="text-xs font-mono font-thin fill-gray-400 cursor-default"
              >
                <tspan fill={labelColor(arc.data)}>
                  {labelText(arc.data) + ": "}
                </tspan>
                <tspan>{labelValue(arc.data)}</tspan>
              </text>
            </Group>
            <LinePath<Point>
              key={"line-" + index}
              curve={curveBasis}
              data={linePoints.get(index)}
              stroke={labelColor(arc.data)}
              x={(p) => p.x}
              y={(p) => p.y}
            />
          </Group>
        );
      })}
    </>
  );
};

function useLabelPositions<Datum>(
  arcs: PieArcDatum<Datum>[],
  labelRadius: number,
  xPadding: number,
  yPadding: number
): Map<number, LabelPosition> {
  return useMemo(() => {
    const initial: LabelPosition[] = arcs.map((arc, index) => {
      const mid = midAngle(arc);
      return {
        index: index,
        angle: mid,
        x: Math.sin(mid) * labelRadius,
        y: -Math.cos(mid) * labelRadius,
        value: (arc.startAngle - arc.endAngle) / (2 * Math.PI),
      };
    });
    initial.sort((a, b) => b.angle - a.angle);

    let prev: LabelPosition | null = null;
    return new Map(
      initial.map((pos) => {
        if (
          prev != null &&
          Math.abs(pos.y - prev.y) < yPadding &&
          isLeft(pos.angle) == isLeft(prev.angle)
        ) {
          const r = labelRadius;
          const y = prev.y + yPadding * (isLeft(pos.angle) ? 1 : -1);
          const x = -Math.sqrt(r * r - y * y);
          pos.y = y;
          pos.x = x;
        }
        pos.x += xPadding * (isLeft(pos.angle) ? -1 : 1);
        return [pos.index, (prev = pos)];
      })
    );
  }, [arcs, labelRadius, xPadding, yPadding]);
}

function useLinePoints(
  labelPositions: Map<number, LabelPosition>,
  pieRadius: number,
  tail: number
) {
  const linePoints = new Map<number, Point[]>();
  labelPositions.forEach((pos, index) => {
    const dx = Math.sin(pos.angle);
    const dy = -Math.cos(pos.angle);

    linePoints.set(index, [
      { x: dx * pieRadius, y: dy * pieRadius },
      { x: dx * (pieRadius + tail), y: dy * (pieRadius + tail) },
      { x: pos.x, y: pos.y - 3 },
    ]);
  });
  return linePoints;
}

function midAngle<Datum>(arc: PieArcDatum<Datum>): number {
  return arc.startAngle + (arc.endAngle - arc.startAngle) / 2;
}

function isLeft(angle: number): boolean {
  return angle > Math.PI;
}
