import type { model } from "@gcsim/types";
import {
	AMBER_700,
	GRAY_400,
	GRAY_600,
	PRIMARY_BG,
	PRIMARY_FG,
	ROSE_700,
	SLATE_700,
} from "./colors";
import { FONT_FAMILY } from "./fonts";

type Intent = "default" | "warning" | "danger";

// Matches the live badge variants: default = slate-900 with a muted title, and
// warning/amber and danger/rose fill the whole pill with white text.
const intentStyles: Record<
	Intent,
	{ bg: string; titleFg: string; fg: string }
> = {
	default: { bg: PRIMARY_BG, titleFg: GRAY_400, fg: PRIMARY_FG },
	warning: { bg: AMBER_700, titleFg: "#f5f5f5", fg: "#ffffff" },
	danger: { bg: ROSE_700, titleFg: "#ffe4e6", fg: "#ffffff" },
};

type ItemProps = {
	title?: string;
	value: string;
	intent?: Intent;
	valueCase?: "uppercase" | "lowercase" | "none";
};

// A single metadata pill. Hook-free equivalent of the live Metadata Item/Badge
// (font-mono, text-sm bold value, text-xs title, px-2.5 py-1.5).
const Item = ({
	title,
	value,
	intent = "default",
	valueCase = "uppercase",
}: ItemProps) => {
	const { bg, titleFg, fg } = intentStyles[intent];
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				alignItems: "center",
				gap: 8,
				padding: "6px 10px",
				borderRadius: 4,
				border: "1px solid transparent",
				backgroundColor: bg,
				fontFamily: FONT_FAMILY,
			}}
		>
			{title != null ? (
				<span
					style={{
						fontSize: 12,
						lineHeight: "16px",
						color: titleFg,
						textTransform: "lowercase",
					}}
				>
					{title}
				</span>
			) : null}
			<span
				style={{
					fontSize: 14,
					lineHeight: "16px",
					fontWeight: 700,
					color: fg,
					textTransform: valueCase === "none" ? "none" : valueCase,
				}}
			>
				{value}
			</span>
		</div>
	);
};

const dpsFmt = new Intl.NumberFormat("en", {
	notation: "compact",
	minimumSignificantDigits: 3,
	maximumSignificantDigits: 3,
});
const compactFmt = new Intl.NumberFormat("en", { notation: "compact" });
const plainFmt = new Intl.NumberFormat("en");

type Props = {
	data: model.SimulationResult;
};

// The metadata row: build state, DPS/target, warnings, iterations and mode.
// Uses Intl.NumberFormat instead of react-i18next.
export const Metadata = ({ data }: Props) => {
	const row = (children: React.ReactNode) => (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				flexWrap: "wrap",
				gap: 8,
				justifyContent: "center",
				alignItems: "center",
				padding: 6,
				margin: "0 4px 4px",
				borderRadius: 4,
				border: `1px solid ${GRAY_600}`,
				backgroundColor: SLATE_700,
			}}
		>
			{children}
		</div>
	);

	if (data.schema_version == null) {
		return row(<Item value="legacy sim" intent="danger" />);
	}

	// @ts-ignore: generated proto uses lower-case `dps`, not `DPS`.
	let dps: number | undefined = data?.statistics?.dps?.mean;
	const targetCount = Object.keys(data?.statistics?.target_dps ?? {}).length;
	if (targetCount > 0 && dps !== undefined) {
		dps = dps / targetCount;
	} else {
		dps = undefined;
	}

	const isProd = data.key_type == null || data.key_type === "prod";
	const showDps = !data.modified && isProd;

	const warningCount = Object.entries(data?.statistics?.warnings ?? {}).filter(
		([, v]) => v as boolean,
	).length;

	const items: React.ReactNode[] = [];

	// Build state.
	if (isProd) {
		if (data.modified) {
			items.push(<Item key="dirty" value="dirty" intent="danger" />);
		}
	} else if (data.key_type === "dev") {
		items.push(<Item key="build" value="dev build" intent="danger" />);
	} else {
		items.push(<Item key="build" value="unofficial" intent="danger" />);
	}

	if (showDps) {
		items.push(
			<Item
				key="dps"
				title="dps/target"
				value={dps !== undefined ? dpsFmt.format(dps) : "n/a"}
			/>,
		);
	}

	if (warningCount > 0) {
		items.push(
			<Item
				key="warn"
				title="warnings"
				value={plainFmt.format(warningCount)}
				intent="warning"
			/>,
		);
	}

	if (data?.statistics?.iterations != null) {
		items.push(
			<Item
				key="itr"
				title="iterations"
				value={compactFmt.format(data.statistics.iterations)}
			/>,
		);
	}

	if (data.mode != null) {
		items.push(
			<Item
				key="mode"
				title="mode"
				value={data.mode === 2 ? "ttk" : "duration"}
			/>,
		);
	}

	return row(items);
};
