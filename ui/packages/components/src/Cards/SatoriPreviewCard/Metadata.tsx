import type { model } from "@gcsim/types";
import { GRAY_400, GRAY_600, GRAY_700, SLATE_700 } from "./colors";

type Intent = "default" | "warning" | "danger";

const intentStyles: Record<Intent, { bg: string; fg: string }> = {
	default: { bg: GRAY_700, fg: "#e5e7eb" },
	warning: { bg: "#78350f", fg: "#fcd34d" }, // amber-900 / amber-300
	danger: { bg: "#7f1d1d", fg: "#fca5a5" }, // red-900 / red-300
};

type ItemProps = {
	title?: string;
	value: string;
	intent?: Intent;
	valueCase?: "uppercase" | "lowercase" | "none";
};

// A single metadata pill. Hook-free equivalent of the live Metadata Item/Badge.
const Item = ({
	title,
	value,
	intent = "default",
	valueCase = "uppercase",
}: ItemProps) => {
	const { bg, fg } = intentStyles[intent];
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				alignItems: "center",
				gap: 4,
				padding: "1px 6px",
				borderRadius: 4,
				backgroundColor: bg,
				fontFamily: "Inter",
			}}
		>
			{title != null ? (
				<span
					style={{ fontSize: 10, color: GRAY_400, textTransform: "lowercase" }}
				>
					{title}
				</span>
			) : null}
			<span
				style={{
					fontSize: 12,
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
				padding: 4,
				margin: "0 4px",
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
				valueCase="lowercase"
			/>,
		);
	}

	return row(items);
};
