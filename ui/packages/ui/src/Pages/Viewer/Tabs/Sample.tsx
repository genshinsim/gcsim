import {
	DefaultSampleOptions,
	namedSeeds,
	SampleLog,
	SeedPicker,
} from "@gcsim/components";
import { NonIdealState } from "@gcsim/primitives";
import type { model, Sample } from "@gcsim/types";
import { FlaskConical } from "lucide-react";
import queryString from "query-string";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { downloadSample } from "../../Sample/downloadSample";

const SAVED_SAMPLE_KEY = "gcsim-sample-settings";

type UseSampleData = {
	sample?: Sample;
	seed: string | null;
	settings: string[];
	generating: boolean;
	setGenerating: (val: boolean) => void;
	setSample: (sample?: Sample) => void;
	setSettings: (val: string[]) => void;
	setSeed: (val: string | null) => void;
};

type Props = {
	sampler: (cfg: string, seed: string) => Promise<Sample>;
	data: model.SimulationResult | null;
	sample: UseSampleData;
	running: boolean;
};

export default ({ sampler, data, sample, running }: Props) => {
	if (
		data?.character_details == null ||
		data.config_file == null ||
		sample.generating
	) {
		return <NonIdealState loading />;
	}

	const generate = (seed: string) => {
		const parsed = queryString.parse(location.hash);
		parsed.sample = seed;
		location.hash = queryString.stringify(parsed);

		sample.setGenerating(true);
		sample.setSeed(seed);
		sampler(data.config_file ?? "", seed).then((out) => {
			sample.setSample(out);
			sample.setGenerating(false);
		});
	};

	const picker = (
		<SeedPicker
			seeds={namedSeeds(data)}
			value={sample.seed}
			onPick={generate}
			running={running}
		/>
	);

	if (sample.sample == null) {
		return (
			<NonIdealState
				icon={<FlaskConical />}
				action={picker}
				className="!px-2"
			/>
		);
	}

	return (
		<div className="w-full 2xl:mx-auto 2xl:container flex flex-grow flex-col gap-[15px] px-2">
			{picker}
			<SampleLog
				sample={sample.sample}
				settings={sample.settings}
				onSettingsChange={sample.setSettings}
				onDownload={downloadSample}
			/>
		</div>
	);
};

export function useSample(
	running: boolean,
	data: model.SimulationResult | null,
	sampleOnLoad: boolean,
	sampler: (cfg: string, seed: string) => Promise<Sample>,
): UseSampleData {
	const [selected, setSelected] = useState<string[]>(() => {
		const saved = localStorage.getItem(SAVED_SAMPLE_KEY);
		if (saved) {
			const initialValue = JSON.parse(saved);
			return initialValue || DefaultSampleOptions;
		}
		return DefaultSampleOptions;
	});

	// TODO(react19): drop this useCallback and inline the function once the
	// React Compiler is enabled — it auto-memoizes. Don't remove it before
	// then: useExhaustiveDependencies is now an error and would fail the build.
	const setAndStore = useCallback((val: string[]) => {
		setSelected(val);
		localStorage.setItem(SAVED_SAMPLE_KEY, JSON.stringify(val));
	}, []);

	const [sample, SetSample] = useState<Sample | undefined>(undefined);
	const [generating, setGenerating] = useState(false);
	const [seed, setSeed] = useState<string | null>(null);

	// Special case where sim is rerunning. Want to reset any generated sample state
	useEffect(() => {
		if (running) {
			SetSample(undefined);
		}
	}, [running]);

	const initQuery = useRef(queryString.parse(location.hash));

	// if seed in url or sampleOnLoad is checked, load sample on viewer load
	useEffect(() => {
		const linkSeed = initQuery.current.sample as string;
		if (
			(sampleOnLoad || linkSeed) &&
			sample == null &&
			!generating &&
			data?.config_file != null
		) {
			const seed = linkSeed ?? data.sample_seed;

			setGenerating(true);
			setSeed(seed);
			sampler(data.config_file ?? "", seed).then((out) => {
				SetSample(out);
				setGenerating(false);
			});
		}
	}, [
		data?.config_file,
		data?.sample_seed,
		generating,
		sample,
		sampleOnLoad,
		sampler,
	]);

	return useMemo(() => {
		return {
			sample: sample,
			seed: seed,
			settings: selected,
			generating: generating,
			setGenerating: setGenerating,
			setSample: SetSample,
			setSettings: setAndStore,
			setSeed: setSeed,
		};
	}, [generating, sample, seed, selected, setAndStore]);
}
