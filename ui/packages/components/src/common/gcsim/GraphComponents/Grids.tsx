import { Grid } from "@visx/grid";
import { GridProps } from "@visx/grid/lib/grids/Grid";
import GridColumns, { AllGridColumnsProps } from "@visx/grid/lib/grids/GridColumns";
import GridRows, { AllGridRowsProps } from "@visx/grid/lib/grids/GridRows";
import { GridScale } from "@visx/grid/lib/types";
import { DataColorsConst } from "./DataColors";

export const GraphGrid = <XScale extends GridScale, YScale extends GridScale>(
      props: GridProps<XScale, YScale>
    ) => {
  return <Grid stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};

export const GraphGridRows = <Scale extends GridScale,>(
      props: AllGridRowsProps<Scale>
    ) => {
  return <GridRows stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};

export const GraphGridColumns = <Scale extends GridScale,>(
      props: AllGridColumnsProps<Scale>
    ) => {
  return <GridColumns stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};