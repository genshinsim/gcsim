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

export function autoSampleKey(
	config: string | undefined,
	autoSeed: string | null,
): string | null {
	if (config == null || autoSeed == null) {
		return null;
	}
	return `${autoSeed}:${config}`;
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
};

export function useSample({
	config,
	autoSeed = null,
}: SampleSource = {}): SampleState {
	const { exec } = useExecutor();
	const [sample, setSample] = useState<Sample | null>(null);
	const [seed, setSeed] = useState<string | null>(null);
	const [generating, setGenerating] = useState(false);
	const autoSampled = useRef<string | null>(null);
	const request = useRef(0);

	const generate = useCallback(
		(next: string) => {
			if (config == null) {
				return;
			}
			const id = ++request.current;
			setGenerating(true);
			setSeed(next);
			exec()
				.sample(config, next)
				.then((result) => {
					if (id === request.current) {
						setSample(result);
					}
				})
				.finally(() => {
					if (id === request.current) {
						setGenerating(false);
					}
				});
		},
		[exec, config],
	);

	useEffect(() => {
		const key = autoSampleKey(config, autoSeed);
		if (key == null) {
			request.current++;
			autoSampled.current = null;
			setSample(null);
			setGenerating(false);
			return;
		}
		if (autoSampled.current === key || autoSeed == null) {
			return;
		}
		autoSampled.current = key;
		generate(autoSeed);
	}, [autoSeed, config, generate]);

	return { sample, seed, generating, generate };
}
