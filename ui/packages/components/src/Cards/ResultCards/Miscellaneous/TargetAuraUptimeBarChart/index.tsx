import { dynamicKey } from "@gcsim/localization";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useDeferredValue, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, useDataColors } from "../../../../common/gcsim";
import { BarChart } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
};

export default ({ data }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();
	const deferred = useDeferredValue(data);
	const targetAuraUptime = deferred?.statistics?.target_aura_uptime;
	const auraUptime = useMemo(
		() =>
			targetAuraUptime?.map((s) =>
				s.sources
					? {
							sources: Object.fromEntries(
								Object.entries(s.sources).map(([k, v]) => [
									t(dynamicKey("elements." + k)),
									v,
								]),
							),
						}
					: {},
			),
		[targetAuraUptime, t],
	);

	const targets = useMemo(() => {
		if (auraUptime == null) {
			return [];
		}

		const targets = new Set<string>();
		for (let i = 0; i < auraUptime.length; i++) {
			targets.add(i.toString());
		}
		return Array.from(targets);
	}, [auraUptime]);
	const [target, setTarget] = useState("0");

	const auras = useMemo(() => {
		if (auraUptime == null) {
			return [];
		}

		const auras = new Set<string>();
		for (const key in auraUptime[target]?.sources) {
			auras.add(key);
		}
		return Array.from(auras).sort(
			(a, b) =>
				DataColors.reactableModifierKeys.indexOf(a) -
				DataColors.reactableModifierKeys.indexOf(b),
		);
	}, [auraUptime, DataColors.reactableModifierKeys, target]);

	return (
		<Card className="flex flex-col col-span-3 min-h-[384px] p-5">
			<div className="flex flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle title={t("result.target_aura_uptime")} />
					<Options target={target} setTarget={setTarget} targets={targets} />
				</div>
			</div>
			<ParentSize className="flex-grow">
				{({ width, height }) => (
					<BarChart
						width={width}
						height={height}
						auraUptime={auraUptime}
						auras={auras}
						target={target}
					/>
				)}
			</ParentSize>
		</Card>
	);
};

const Options = ({
	target,
	setTarget,
	targets,
}: {
	target: string;
	setTarget: (v: string) => void;
	targets: string[];
}) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-row items-center gap-2 mb-2">
			<span className="text-g-xs font-g-mono text-g-ink-mute">
				{t("viewer.target")}
			</span>
			<select
				className="rounded-g-md border border-g-line bg-transparent px-2 py-1 text-g-sm"
				value={target}
				onChange={(e) => setTarget(e.target.value)}
			>
				{targets.map((target) => (
					<option key={target} value={target}>
						{Number(target) + 1}
					</option>
				))}
			</select>
		</div>
	);
};
