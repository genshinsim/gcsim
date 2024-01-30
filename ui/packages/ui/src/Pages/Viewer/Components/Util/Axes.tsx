import { AxisBottom, AxisLeft, AxisRight, AxisScale, SharedAxisProps, TickFormatter, TickRendererProps } from "@visx/axis";
import { ScaleInput } from "@visx/scale";
import { Text } from "@visx/text";
import { DataColorsConst } from "./DataColors";

type CustomProps<Scale extends AxisScale> = {
  tickFormat?: TickFormatter<ScaleInput<Scale>>
  tickLabelX?: string | number; 
  tickLabelY?: string | number;
}

export const GraphAxisLeft = <Scale extends AxisScale,>({
      tickLabelX,
      tickLabelY = "0.25em",
      ...restProps
    }: CustomProps<Scale> & SharedAxisProps<Scale>) => {
  return (
    <AxisLeft
        stroke={DataColorsConst.gray}
        tickStroke={DataColorsConst.gray}
        tickLineProps={{ opacity: 0.5 }}
        labelClassName="fill-gray-400 text-lg"
        tickClassName="fill-gray-400 font-mono text-xs"
        tickComponent={(props) => (
            <TickLabel {...props} dx={tickLabelX} dy={tickLabelY} textAnchor="end" />
        )}
        {...restProps} />
  );
};

export const GraphAxisBottom = <Scale extends AxisScale,>({
      tickLabelX,
      tickLabelY,
      ...restProps
    }: CustomProps<Scale> & SharedAxisProps<Scale>) => {
  return (
    <AxisBottom
        stroke={DataColorsConst.gray}
        tickStroke={DataColorsConst.gray}
        tickLineProps={{ opacity: 0.5 }}
        labelClassName="fill-gray-400 font-mono text-base"
        tickClassName="fill-gray-400 font-mono text-xs"
        tickComponent={(props) => (
            <TickLabel {...props} dx={tickLabelX} dy={tickLabelY} textAnchor="middle" />
        )}
        {...restProps} />
  );
};

export const GraphAxisRight = <Scale extends AxisScale,>({
      tickLabelX,
      tickLabelY,
      ...restProps
    }: CustomProps<Scale> & SharedAxisProps<Scale>) => {
  return (
  <AxisRight
      stroke={DataColorsConst.gray}
      tickStroke={DataColorsConst.gray}
      tickLineProps={{ opacity: 0.5 }}
      labelClassName="fill-gray-400 text-lg"
      tickClassName="fill-gray-400 font-mono text-xs"
      tickComponent={(props) => (
          <TickLabel {...props} dx={tickLabelX} dy={tickLabelY} textAnchor="start" />
      )}
      {...restProps} />
  );
};

const TickLabel = (props: TickRendererProps) => {
  return (
    <Text
        x={props.x}
        y={props.y}
        dx={props.dx}
        dy={props.dy}
        textAnchor={props.textAnchor}
        className="cursor-default">
      {props.formattedValue}
    </Text>
  );
};