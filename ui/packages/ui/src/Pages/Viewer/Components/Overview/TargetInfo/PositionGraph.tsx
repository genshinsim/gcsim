import { Coord, Enemy } from "@gcsim/types";
import { Group } from "@visx/group";
import { scaleLinear } from "@visx/scale";
import { Circle } from "@visx/shape";
import { useTooltip } from "@visx/tooltip";
import { useMemo } from "react";
import { DataColorsConst, GraphAxisBottom, GraphAxisLeft, GraphAxisRight, GraphGrid, NoData } from "../../Util";
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
  const { data, max, centerX, centerY } = useData(enemies, player);
  const size = Math.min(width - margin.left - margin.right, height - margin.top - margin.bottom);
  const gridSize = Math.max(Math.round(2 * max + 1), 1) / size;
  const marginLeft = (width - size) / 2;
  const domain = (size * gridSize) / 2;

  const xScale = scaleLinear<number>({
    range: [0, size],
    domain: [centerX - domain, centerX + domain],
  });

  const yScale = scaleLinear<number>({
    range: [size, 0], 
    domain: [centerY - domain, centerY + domain],
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
            const opacity = (tooltip.tooltipData?.index == i && !tooltip.tooltipData.player)
                ? 0.75 : 0.25;
            return (
              <Circle
                key={`enemy-${i}`}
                cx={xScale(e.x)}
                cy={yScale(e.y)}
                r={sizeScale(e.r)}
                fillOpacity={opacity}
                fill={DataColorsConst.qualitative3(i)}
                stroke={DataColorsConst.qualitative3(i)}
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
              fill={DataColorsConst.gray}
              stroke={DataColorsConst.gray}
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
  centerX: number;
  centerY: number;
}

function useData(enemies?: Enemy[], player?: Coord): PositionData {
  return useMemo(() => {
    if (enemies == null || player == null) {
      return { data: [], max: 0, centerX: 0, centerY: 0 };
    }

    const playerX = player.x ?? 0;
    const playerY = player.y ?? 0;

    let max = 0;
    const data: Position[] = enemies.map((e) => {
      const x = e.position?.x ?? 0;
      const y = e.position?.y ?? 0;
      const r = e.position?.r ?? 1;
      const dist = Math.sqrt(Math.pow(playerX - x, 2) + Math.pow(playerY - y, 2));
      max = Math.max(max, dist + r);
      return { x: x, y: y, r: r};
    });

    return {
      data: data,
      max: max,
      centerX: playerX,
      centerY: playerY,
    };
  }, [enemies, player]);
}