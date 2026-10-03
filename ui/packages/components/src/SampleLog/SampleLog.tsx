import type { Sample } from "@gcsim/types";
import React, { useMemo, useRef, useState } from "react";
import { useSampleFilter } from "./sampleFilter";
import { CommandBar } from "./timeline/CommandBar";
import { chipText, modelFromSample } from "./timeline/model";
import { type SearchHits, Strip, type StripHandle } from "./timeline/Strip";
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
	const hits = useMemo<SearchHits>(
		() => ({
			ids: new Set(matchList.map((c) => c.id)),
			frames: matchList.map((c) => c.frame),
		}),
		[matchList],
	);
	const [cursor, setCursor] = useState({ list: matchList, at: -1 });
	const at = cursor.list === matchList ? cursor.at : -1;

	const next = () => {
		if (matchList.length === 0) {
			return;
		}
		const i = (at + 1) % matchList.length;
		setCursor({ list: matchList, at: i });
		strip.current?.centerOn(matchList[i].frame);
	};

	return (
		<div ref={fit} className="flex h-[80vh] flex-col gap-2">
			<Strip
				ref={strip}
				model={model}
				enabled={enabled}
				hits={hits}
				onOpen={details.open}
				commandBar={
					<CommandBar
						model={model}
						filter={filter}
						onFilterChange={setFilter}
						search={search}
						onSearch={setSearch}
						onSearchNext={next}
						matchLabel={
							search.trim() === "" ? "" : `${at + 1}/${matchList.length}`
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
