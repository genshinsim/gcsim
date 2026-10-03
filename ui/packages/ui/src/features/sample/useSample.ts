import { useExecutor } from "@gcsim/components";
import type { Sample } from "@gcsim/types";
import { useCallback, useEffect, useRef, useState } from "react";

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
};

type SampleSource = {
	config?: string;
	autoSeed?: string | null;
	running?: boolean;
};

export function useSample({
	config,
	autoSeed = null,
	running = false,
}: SampleSource = {}): SampleState {
	const { exec } = useExecutor();
	const [sample, setSample] = useState<Sample | null>(null);
	const [seed, setSeed] = useState<string | null>(null);
	const [generating, setGenerating] = useState(false);
	const autoSampled = useRef(false);

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

	return { sample, seed, generating, generate };
}
