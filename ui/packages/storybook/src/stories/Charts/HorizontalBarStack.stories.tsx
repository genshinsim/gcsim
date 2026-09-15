import {
	FloatStatTooltipContent,
	HorizontalBarStack,
	NoData,
	useDataColors,
} from "@gcsim/components";
import type { model } from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { range } from "lodash-es";
import { sampleResult } from "../samples";

type Datum = { name: string; data: model.DescriptiveStats; index: number };

type Props = {
	result?: model.SimulationResult;
	width: number;
	height: number;
};

const CharacterDPSBarStack = ({ result, width, height }: Props) => {
	const { DataColors } = useDataColors();

	const dps = result?.statistics?.character_dps;
	const names = result?.character_details?.map((c) => c.name ?? "");
	if (dps == null || names == null) {
		return <NoData />;
	}

	const data: Datum[] = dps.map((d, i) => ({
		name: names[i],
		data: d,
		index: i,
	}));
	const xMax = data.reduce(
		(m, d) =>
			Math.max(m, d.data.max ?? 0, (d.data.mean ?? 0) + (d.data.sd ?? 0)),
		0,
	);

	return (
		<HorizontalBarStack<Datum, number>
			width={width}
			height={height}
			xDomain={[0, xMax]}
			yDomain={names}
			y={(d) => d.name}
			data={data}
			keys={range(names.length)}
			value={(d, k) => (d.index === k ? (d.data.mean ?? 0) : 0)}
			stat={(d) => d.data}
			barColor={(k) => DataColors.character(k)}
			hoverColor={(k) => DataColors.characterLabel(k)}
			bottomLabel="DPS"
			tooltipContent={(d, k) => (
				<FloatStatTooltipContent
					title={d.name + " DPS"}
					data={d.data}
					color={DataColors.characterLabel(k)}
					percent={1}
				/>
			)}
		/>
	);
};

const meta: Meta<typeof CharacterDPSBarStack> = {
	title: "Charts/HorizontalBarStack",
	component: CharacterDPSBarStack,
	tags: ["autodocs"],
	args: {
		width: 520,
		height: 220,
	},
	decorators: [
		(Story, ctx) => (
			<div style={{ width: ctx.args.width }}>
				<Story />
			</div>
		),
	],
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		result: sampleResult as model.SimulationResult,
	},
};

export const Running: Story = {
	args: {
		result: undefined,
	},
};
