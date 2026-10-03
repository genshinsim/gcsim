import { Dialog, DialogContent, DialogTitle } from "@gcsim/primitives";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { AllSampleOptions } from "../SampleOptions";
import { EventChip } from "./EventChip";
import { LaneIcon, laneColor, laneTint, useLaneName } from "./lane";
import { type Column, secondsLabel } from "./layout";
import { type Chip, onFieldLaneAt, type TimelineModel } from "./model";

const catRank = (c: string) => {
	const i = AllSampleOptions.indexOf(c);
	return i === -1 ? AllSampleOptions.length : i;
};

export function FrameView({
	model,
	cols,
	index,
	onIndex,
	onClose,
	matches,
	onOpen,
	onDurationIcon,
}: {
	model: TimelineModel;
	cols: Column[];
	index: number | null;
	onIndex: (i: number) => void;
	onClose: () => void;
	matches: Set<number>;
	onOpen: (c: Chip) => void;
	onDurationIcon: (c: Chip) => void;
}) {
	const { t } = useTranslation();
	const laneName = useLaneName();
	const [byCat, setByCat] = useState(false);
	const body = useRef<HTMLDivElement>(null);
	const col = index != null ? cols[index] : undefined;

	useEffect(() => {
		if (index == null) {
			return;
		}
		const onKey = (e: KeyboardEvent) => {
			if (e.key !== "ArrowLeft" && e.key !== "ArrowRight") {
				return;
			}
			if (e.key === "ArrowLeft" && index > 0) {
				onIndex(index - 1);
			} else if (e.key === "ArrowRight" && index < cols.length - 1) {
				onIndex(index + 1);
			}
		};
		window.addEventListener("keydown", onKey);
		return () => window.removeEventListener("keydown", onKey);
	}, [index, cols.length, onIndex]);

	useEffect(() => {
		if (body.current != null && col != null) {
			body.current.scrollTop = 0;
		}
	}, [col]);

	const field = col != null ? onFieldLaneAt(model.onField, col.frame) : -1;
	const i = index ?? 0;
	const prev = cols[i - 1];
	const next = cols[i + 1];
	const total = col?.cells.reduce((a, c) => a + c.length, 0) ?? 0;

	return (
		<Dialog
			open={col != null}
			onOpenChange={(o) => {
				if (!o) {
					onClose();
				}
			}}
		>
			<DialogContent className="flex h-[92vh] w-[96vw] max-w-none flex-col gap-0 overflow-hidden p-0 sm:max-w-none">
				{col != null && (
					<>
						<div className="flex shrink-0 items-center gap-3 border-b border-g-line px-4 py-2.5 pr-12">
							<button
								type="button"
								disabled={prev == null}
								onClick={() => onIndex(i - 1)}
								className="flex items-center rounded-g-sm p-1 text-g-ink-mute hover:bg-g-surface-2 hover:text-g-ink disabled:opacity-30"
								aria-label={t("sample.previous_frame")}
								title={t("sample.previous_frame")}
							>
								<ChevronLeft className="size-4" />
							</button>
							<DialogTitle className="font-g-mono text-g-base">
								{t("sample.frame_title", { frame: col.frame })}
								<span className="ml-2 text-g-ink-mute">
									{secondsLabel(col.frame, 3)}
								</span>
							</DialogTitle>
							<button
								type="button"
								disabled={next == null}
								onClick={() => onIndex(i + 1)}
								className="flex items-center rounded-g-sm p-1 text-g-ink-mute hover:bg-g-surface-2 hover:text-g-ink disabled:opacity-30"
								aria-label={t("sample.next_frame")}
								title={t("sample.next_frame")}
							>
								<ChevronRight className="size-4" />
							</button>
							<span className="font-g-mono text-[11px] text-g-ink-mute">
								{t("sample.events_in_frame", { count: total })}
								{prev != null &&
									` · ${t("sample.since_last", { count: col.frame - prev.frame })}`}
							</span>
							<label className="ml-auto flex cursor-pointer items-center gap-1.5 text-[11px] text-g-ink-dim">
								<input
									type="checkbox"
									checked={byCat}
									onChange={(e) => setByCat(e.target.checked)}
								/>
								{t("sample.group_by_category")}
							</label>
						</div>
						<div ref={body} className="relative min-h-0 flex-1 overflow-y-auto">
							<div
								className="grid min-h-full"
								style={{
									gridTemplateColumns: `repeat(${model.lanes.length}, minmax(0, 1fr))`,
								}}
							>
								{model.lanes.map((lane) => {
									const cell = col.cells[lane.index];
									const events = byCat
										? [...cell].sort(
												(a, b) =>
													catRank(a.event.type) - catRank(b.event.type) ||
													a.id - b.id,
											)
										: cell;
									const onField = field === lane.index;
									return (
										<div
											key={lane.key}
											className="min-w-0 border-l border-g-line first:border-l-0"
											style={
												onField
													? {
															background: laneTint(lane, 6),
														}
													: undefined
											}
										>
											<div
												className="sticky top-0 z-10 flex items-center gap-2 border-b bg-g-surface px-3 py-1.5"
												style={{
													borderBottomColor: onField
														? laneColor(lane)
														: "var(--g-border)",
												}}
												title={laneName(lane)}
											>
												<LaneIcon
													lane={lane}
													name={laneName(lane)}
													className="size-7"
												/>
												{onField && (
													<span className="rounded-g-sm bg-g-surface-3 px-1 text-[10px] text-g-ink-dim">
														{t("sample.on_field")}
													</span>
												)}
												<span className="ml-auto font-g-mono text-[11px] text-g-ink-mute">
													{cell.length}
												</span>
											</div>
											<div className="flex flex-col gap-1 p-1.5">
												{cell.length === 0 && (
													<span className="px-1.5 py-1 text-[11px] text-g-ink-mute">
														{t("sample.nothing_this_frame")}
													</span>
												)}
												{events.map((e, k) => (
													<div key={e.id}>
														{byCat &&
															events[k - 1]?.event.type !== e.event.type && (
																<div className="mt-1 mb-0.5 px-1 text-[10px] font-medium tracking-wide text-g-ink-mute uppercase">
																	{e.event.type.replaceAll("_", " ")}
																</div>
															)}
														<EventChip
															chip={e}
															onOpen={onOpen}
															onDurationIcon={onDurationIcon}
															matched={matches.has(e.id)}
															className="w-full"
														/>
													</div>
												))}
											</div>
										</div>
									);
								})}
							</div>
						</div>
					</>
				)}
			</DialogContent>
		</Dialog>
	);
}
