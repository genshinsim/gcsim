import type { Sample } from "@gcsim/types";
import React, { useMemo, useRef, useState } from "react";
import { useSampleFilter } from "./sampleFilter";
import { CommandBar } from "./timeline/CommandBar";
import { chipText, modelFromSample } from "./timeline/model";
import { Strip, type StripHandle } from "./timeline/Strip";
import { useEventDetails } from "./timeline/useEventDetails";
import { useFitViewport } from "./timeline/useFitViewport";

export type SampleLogProps = {
	sample: Sample;
	onDownload?: (sample: Sample) => void;
};

function SampleLogUI({ sample, onDownload }: SampleLogProps) {
	const [filter, setFilter] = useSampleFilter();
	const model = useMemo(() => modelFromSample(sample), [sample]);
	const enabled = useMemo(() => new Set(filter), [filter]);
	const details = useEventDetails();
	const [search, setSearch] = useState("");
	const [cursor, setCursor] = useState(-1);
	const strip = useRef<StripHandle>(null);
	const fit = useFitViewport();

	const matchList = useMemo(() => {
		const needle = search.trim().toLowerCase();
		if (needle === "") {
			return [];
		}
		return model.chips.filter(
			(c) =>
				enabled.has(c.event.type) && chipText(c).toLowerCase().includes(needle),
		);
	}, [model, enabled, search]);
	const matches = useMemo(
		() => new Set(matchList.map((c) => c.id)),
		[matchList],
	);
	const matchFrames = useMemo(() => matchList.map((c) => c.frame), [matchList]);

	const next = () => {
		if (matchList.length === 0) {
			return;
		}
		const i = (cursor + 1) % matchList.length;
		setCursor(i);
		strip.current?.centerOn(matchList[i].frame);
	};

	return (
		<div
			ref={fit.ref}
			className="flex flex-col gap-2"
			style={{ height: fit.height ?? "80vh" }}
		>
			<Strip
				ref={strip}
				model={model}
				enabled={enabled}
				matches={matches}
				matchFrames={matchFrames}
				onOpen={details.open}
				commandBar={
					<CommandBar
						model={model}
						filter={filter}
						onFilterChange={setFilter}
						search={search}
						onSearch={(v) => {
							setSearch(v);
							setCursor(-1);
						}}
						onSearchNext={next}
						matchLabel={
							search.trim() === ""
								? ""
								: `${cursor >= 0 ? cursor + 1 : 0}/${matchList.length}`
						}
						onDownload={
							onDownload != null ? () => onDownload(sample) : undefined
						}
					/>
				}
			/>
			{details.dialog}
		</div>
	);
}

export const SampleLog = React.memo(SampleLogUI);
