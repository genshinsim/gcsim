import {
	type AllGridColumnsProps,
	type AllGridRowsProps,
	Grid,
	GridColumns,
	GridRows,
	type GridScale,
} from "@visx/grid";
import type { ComponentProps } from "react";
import { DataColorsConst } from "./DataColors";

export const GraphGrid = <XScale extends GridScale, YScale extends GridScale>(
	props: ComponentProps<typeof Grid<XScale, YScale>>,
) => {
	return <Grid stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};

export const GraphGridRows = <Scale extends GridScale>(
	props: AllGridRowsProps<Scale>,
) => {
	return <GridRows stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};

export const GraphGridColumns = <Scale extends GridScale>(
	props: AllGridColumnsProps<Scale>,
) => {
	return <GridColumns stroke={DataColorsConst.gray} opacity={0.5} {...props} />;
};
