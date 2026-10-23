import { Button, NonIdealState } from "@gcsim/primitives";
import type { Sample } from "@gcsim/types";
import { FlaskConical, RefreshCw } from "lucide-react";
import React, { useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSampleFilter } from "./sampleFilter";
import { useSampleView } from "./sampleView";
import { CommandBar } from "./timeline/CommandBar";
import { EventDetailsDialog } from "./timeline/EventDetailsDialog";
import { type Chip, chipText, modelFromSample } from "./timeline/model";
import { type SearchHits, Strip, type StripHandle } from "./timeline/Strip";
import { useFitViewport } from "./timeline/useFitViewport";
import { VerticalStrip } from "./timeline/VerticalStrip";

export type SampleLogProps = {
	sample: Sample | null;
	onDownload?: (sample: Sample) => void;
	/** adds a Generate button; the caller decides what it opens */
	onGenerate?: () => void;
};

function SampleLogUI({ sample, ...props }: SampleLogProps) {
	if (sample == null) {
		return <EmptyLog onGenerate={props.onGenerate} />;
	}
	return <Timeline sample={sample} {...props} />;
}

function EmptyLog({ onGenerate }: { onGenerate?: () => void }) {
	const { t } = useTranslation();
	return (
		<NonIdealState
			className="h-[50vh]"
			icon={<FlaskConical />}
			action={
				onGenerate != null && (
					<Button onClick={onGenerate}>
						<RefreshCw />
						{t("viewer.generate")}
					</Button>
				)
			}
		/>
	);
}

function Timeline({
	sample,
	onDownload,
	onGenerate,
}: SampleLogProps & { sample: Sample }) {
	const [filter, setFilter] = useSampleFilter();
	const model = useMemo(() => modelFromSample(sample), [sample]);
	const enabled = useMemo(() => new Set(filter), [filter]);
	const [detail, setDetail] = useState<Chip | null>(null);
	const [search, setSearch] = useState("");
	const strip = useRef<StripHandle>(null);
	const fit = useFitViewport();
	const [view, setView] = useSampleView();

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

	const commandBar = (inline: boolean) => (
		<CommandBar
			model={model}
			filter={filter}
			onFilterChange={setFilter}
			search={search}
			onSearch={setSearch}
			onSearchNext={next}
			matchLabel={search.trim() === "" ? "" : `${at + 1}/${matchList.length}`}
			onDownload={onDownload != null ? () => onDownload(sample) : undefined}
			onGenerate={onGenerate}
			view={view}
			onViewChange={setView}
			inline={inline}
		/>
	);

	return (
		<div ref={fit} className="flex h-[80vh] flex-col gap-2">
			{view.orientation === "vertical" ? (
				<VerticalStrip
					ref={strip}
					model={model}
					enabled={enabled}
					hits={hits}
					onOpen={setDetail}
					commandBar={commandBar(true)}
					view={view}
					onOverview={(overview) => setView({ ...view, overview })}
				/>
			) : (
				<Strip
					ref={strip}
					model={model}
					enabled={enabled}
					hits={hits}
					onOpen={setDetail}
					commandBar={commandBar(false)}
				/>
			)}
			<EventDetailsDialog chip={detail} onClose={() => setDetail(null)} />
		</div>
	);
}

export const SampleLog = React.memo(SampleLogUI);
