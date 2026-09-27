import {
	Button,
	Input,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { RefreshCw } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export type SeedName = "sample" | "min" | "max" | "p25" | "p50" | "p75";
export type NamedSeeds = Partial<Record<SeedName, string>>;

type Choice = SeedName | "custom";

const statisticSeeds: Choice[] = ["min", "max", "p25", "p50", "p75"];

export function namedSeeds(result: model.SimulationResult): NamedSeeds {
	const stats = result.statistics;
	return {
		sample: result.sample_seed,
		min: stats?.min_seed,
		max: stats?.max_seed,
		p25: stats?.p25_seed,
		p50: stats?.p50_seed,
		p75: stats?.p75_seed,
	};
}

function initialChoice(seeds: NamedSeeds, value: string | null): Choice {
	if (value == null || value === seeds.sample) {
		return "sample";
	}
	const named = statisticSeeds.find((c) => seeds[c as SeedName] === value);
	return named ?? "custom";
}

export type SeedPickerProps = {
	seeds: NamedSeeds;
	value: string | null;
	onPick: (seed: string) => void;
	running?: boolean;
};

export const SeedPicker = ({
	seeds,
	value,
	onPick,
	running = false,
}: SeedPickerProps) => {
	const { t } = useTranslation();
	const [choice, setChoice] = useState<Choice>(() =>
		initialChoice(seeds, value),
	);
	// Seeds are uint64 (serialized as strings) and routinely exceed 2^53, so the
	// custom seed is kept as a string; Number() would silently truncate it.
	const [customSeed, setCustomSeed] = useState<string>(
		() => value ?? String(Math.floor(Number.MAX_SAFE_INTEGER * Math.random())),
	);

	const options: { label: string; value: Choice }[] = [
		{ label: t("viewer.seed_sample"), value: "sample" },
		{ label: t("viewer.seed_min"), value: "min" },
		{ label: t("viewer.seed_max"), value: "max" },
		{ label: t("viewer.seed_p", { p: 25 }), value: "p25" },
		{ label: t("viewer.seed_p", { p: 50 }), value: "p50" },
		{ label: t("viewer.seed_p", { p: 75 }), value: "p75" },
		{ label: t("viewer.seed_custom"), value: "custom" },
	];

	const pick = () => {
		onPick(choice === "custom" ? customSeed : (seeds[choice] ?? "0"));
	};

	return (
		<div className="flex flex-col gap-2 w-full mx-auto">
			<Select value={choice} onValueChange={(v) => setChoice(v as Choice)}>
				<SelectTrigger className="w-full">
					<SelectValue />
				</SelectTrigger>
				<SelectContent>
					{options.map((o) => (
						<SelectItem key={o.value} value={o.value}>
							{o.label}
						</SelectItem>
					))}
				</SelectContent>
			</Select>
			{choice === "custom" ? (
				<Input
					value={customSeed}
					onChange={(e) => setCustomSeed(e.target.value)}
					inputMode="numeric"
					className="w-full"
				/>
			) : null}
			<Button
				size="lg"
				className="w-full"
				disabled={running && statisticSeeds.includes(choice)}
				onClick={pick}
			>
				<RefreshCw />
				{t("viewer.generate")}
			</Button>
		</div>
	);
};
