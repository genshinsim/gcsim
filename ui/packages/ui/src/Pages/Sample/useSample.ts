import { DefaultSampleOptions, useExecutor } from "@gcsim/components";
import type { Sample } from "@gcsim/types";
import { useCallback, useEffect, useRef, useState } from "react";

const SETTINGS_KEY = "gcsim-sample-settings";

type SettingsStorage = Pick<Storage, "getItem" | "setItem">;

export function loadSampleSettings(storage: SettingsStorage): string[] {
	try {
		const raw = storage.getItem(SETTINGS_KEY);
		return (raw && JSON.parse(raw)) || DefaultSampleOptions;
	} catch {
		return DefaultSampleOptions;
	}
}

export function saveSampleSettings(storage: SettingsStorage, val: string[]) {
	storage.setItem(SETTINGS_KEY, JSON.stringify(val));
}

export function autoSampleSeed(
	linkSeed: string | null,
	sampleOnLoad: boolean,
	resultSeed: string | undefined,
): string | null {
	if (linkSeed != null) {
		return linkSeed;
	}
	return sampleOnLoad ? (resultSeed ?? null) : null;
}

export type SampleState = {
	sample: Sample | null;
	seed: string | null;
	generating: boolean;
	generate: (seed: string) => void;
	settings: string[];
	setSettings: (val: string[]) => void;
};

type SampleSource = {
	config?: string;
	autoSeed?: string | null;
	running?: boolean;
};

// Without a source (the standalone sample pages) only the persisted settings are used.
export function useSample({
	config,
	autoSeed = null,
	running = false,
}: SampleSource = {}): SampleState {
	const { exec } = useExecutor();
	const [settings, setSettingsState] = useState(() =>
		loadSampleSettings(localStorage),
	);
	const [sample, setSample] = useState<Sample | null>(null);
	const [seed, setSeed] = useState<string | null>(null);
	const [generating, setGenerating] = useState(false);
	const autoSampled = useRef(false);

	const setSettings = useCallback((val: string[]) => {
		setSettingsState(val);
		saveSampleSettings(localStorage, val);
	}, []);

	const generate = useCallback(
		(next: string) => {
			if (config == null) {
				return;
			}
			setGenerating(true);
			setSeed(next);
			exec()
				.sample(config, next)
				.then(setSample)
				.finally(() => setGenerating(false));
		},
		[exec, config],
	);

	useEffect(() => {
		if (running) {
			setSample(null);
			autoSampled.current = false;
		}
	}, [running]);

	useEffect(() => {
		if (autoSampled.current || running || autoSeed == null || config == null) {
			return;
		}
		autoSampled.current = true;
		generate(autoSeed);
	}, [autoSeed, config, running, generate]);

	return { sample, seed, generating, generate, settings, setSettings };
}
