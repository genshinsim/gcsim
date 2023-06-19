import { Coord, Enemy } from "@gcsim/types";
import { Group } from "@visx/group";
import { scaleLinear } from "@visx/scale";
import { Circle } from "@visx/shape";
import { useTooltip } from "@visx/tooltip";
import { useMemo } from "react";
import { DataColors, GraphAxisBottom, GraphAxisLeft, GraphAxisRight, GraphGrid, NoData } from "../../Util";
import { RenderTooltip, TooltipData, useTooltipHandles } from "./PositionGraphTooltip";

type Props = {
  width: number;
  height: number;
  margin?: { left: number, right: number, top: number, bottom: number };
  enemies?: Enemy[];
  player?: Coord;
}

const defaultMargin = { top: 10, left: 10, right: 0, bottom: 0 };

export const PositionGraph = ({
      width,
      height,
      enemies,
      player,
      margin = defaultMargin, 
    }: Props) => {
  const numXTicks = 10;
  const numYTicks = 10;
  const { data } = useData(enemies);
  const size = Math.min(width - margin.left - margin.right, height - margin.top - margin.bottom);
  const gridSize = 10 / size;
  const marginLeft = (width - size) / 2;

  const xScale = scaleLinear<number>({
    range: [0, size],
    domain: [(-size* gridSize) / 2, (size * gridSize) / 2],
    nice: true,
  });

  const yScale = scaleLinear<number>({
    range: [size, 0],
    domain: [(-size * gridSize) / 2, (size * gridSize) / 2],
    nice: true,
  });

  const sizeScale = (size: number) => size / gridSize;

  const tooltip = useTooltip<TooltipData>();
  const tooltipHandles = useTooltipHandles(tooltip.showTooltip, tooltip.hideTooltip);

  if (enemies == null || data.length == 0) {
    return <NoData />;
  }

  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={marginLeft} top={margin.top}>
          <GraphGrid
              opacity={0.35}
              strokeDasharray="0 5 0"
              numTicksColumns={numXTicks}
              xScale={xScale}
              numTicksRows={numYTicks}
              yScale={yScale}
              width={size}
              height={size}
          />
          <GraphAxisLeft
              hideTicks
              numTicks={0}
              scale={yScale}
              axisLineClassName="stroke-2"
          />
          <GraphAxisBottom
              hideTicks
              numTicks={0}
              top={size}
              scale={xScale}
              axisLineClassName="stroke-2"
          />
          <GraphAxisBottom
              hideTicks
              numTicks={0}
              top={0}
              scale={xScale}
              axisLineClassName="stroke-2"
          />
          <GraphAxisRight
              hideTicks
              numTicks={0}
              left={size}
              scale={yScale}
              axisLineClassName="stroke-2"
          />
          {data.map((e, i) => {
            const opacity = tooltip.tooltipData?.index == i ? 0.75 : 0.25;
            return (
              <Circle
                key={`enemy-${i}`}
                cx={xScale(e.x)}
                cy={yScale(e.y)}
                r={sizeScale(e.r)}
                fillOpacity={opacity}
                fill={DataColors.qualitative3(i)}
                stroke={DataColors.qualitative3(i)}
                strokeWidth={1}
                onMouseMove={(ev) => tooltipHandles.mouseHover(ev, {
                  player: false,
                  index: i,
                  x: e.x,
                  y: e.y,
                  r: e.r,
                })}
                onMouseLeave={() => tooltipHandles.mouseLeave()}
              />
            );
          })}
          {player != null && (
            <Circle
              cx={xScale(player.x ?? 0)}
              cy={yScale(player.y ?? 0)}
              r={sizeScale(player.r ?? 0.3)}
              fillOpacity={tooltip.tooltipData?.player ? 0.75 : 0}
              fill={DataColors.gray}
              stroke={DataColors.gray}
              strokeWidth={2}
              onMouseMove={(ev) => tooltipHandles.mouseHover(ev, {
                player: true,
                index: 0,
                x: player.x,
                y: player.y,
                r: player.r,
              })}
              onMouseLeave={() => tooltipHandles.mouseLeave()}
            />
          )}
        </Group>
      </svg>
      <RenderTooltip
          tooltipOpen={tooltip.tooltipOpen}
          tooltipData={tooltip.tooltipData}
          tooltipLeft={tooltip.tooltipLeft}
          tooltipTop={tooltip.tooltipTop}
          handles={tooltipHandles}
          showTooltip={tooltip.showTooltip}
      />
    </div>
  );
};

type Position = {
  x: number;
  y: number;
  r: number;
}

type PositionData = {
  data: Position[];
  max: number;
}

function useData(enemies?: Enemy[]): PositionData {
  return useMemo(() => {
    if (enemies == null) {
      return { data: [], max: 0 };
    }

    let max = 0;
    const data: Position[] = enemies.map((e) => {
      const x = e.position?.x ?? 0;
      const y = e.position?.y ?? 0;
      const r = e.position?.r ?? 1;
      const edgeX = Math.abs(x) + r;
      const edgeY = Math.abs(y) + r; 
      max = Math.max(max, edgeX, edgeY);
      return { x: x, y: y, r: r};
    });

    return {
      data: data,
      max: max,
    };
  }, [enemies]);
}