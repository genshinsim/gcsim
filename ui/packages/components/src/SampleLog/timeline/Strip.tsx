import { cn } from "@gcsim/primitives";
import { ChevronsDown, ChevronsUp, Search, X } from "lucide-react";
import {
	type ReactNode,
	type Ref,
	useEffect,
	useImperativeHandle,
	useMemo,
	useState,
} from "react";
import { useTranslation } from "react-i18next";
import type { SimEvent } from "../events/types";
import { EventChip } from "./EventChip";
import { FrameView } from "./FrameView";
import { LaneIcon, laneColor, laneTint, useLaneName } from "./lane";
import {
	CELL_PAD,
	COL_W,
	chipCap,
	chipsShown,
	columns,
	GAP_BIG,
	gutterWidth,
	HEAD_H,
	laneHeights,
	onFieldRuns,
	portraitSize,
	secondsLabel,
	visibleWindow,
} from "./layout";
import { Minimap } from "./Minimap";
import {
	type Chip,
	isSimLane,
	type Lane,
	lowerBound,
	onFieldLaneAt,
	spanLabel,
	type TimelineModel,
} from "./model";
import { useScrollView } from "./useScrollView";

export type StripHandle = { centerOn: (frame: number) => void };

export type SearchHits = { ids: Set<number>; frames: number[] };

type Highlight = { event: SimEvent; start: number; end: number; label: string };

export function Strip({
	ref,
	model,
	enabled,
	hits,
	onOpen,
	commandBar,
}: {
	ref?: Ref<StripHandle>;
	model: TimelineModel;
	enabled: Set<string>;
	hits: SearchHits;
	onOpen: (c: Chip) => void;
	commandBar: ReactNode;
}) {
	const { t } = useTranslation();
	const laneName = useLaneName();
	const { el, ref: scrollRef, view } = useScrollView();
	const [topOpen, setTopOpen] = useState(
		() => window.matchMedia?.("(min-width: 640px)").matches ?? true,
	);
	const [slot, setSlot] = useState<HTMLDivElement | null>(null);
	const [avail, setAvail] = useState(0);
	useEffect(() => {
		if (slot == null) {
			return;
		}
		const ro = new ResizeObserver(() => setAvail(slot.clientHeight));
		ro.observe(slot);
		return () => ro.disconnect();
	}, [slot]);
	const cap = chipCap(avail, model.lanes.length);
	const portrait = portraitSize(cap);
	const gutter = gutterWidth(portrait);
	const [expanded, setExpanded] = useState<Set<number>>(new Set());

	const cols = useMemo(() => columns(model, enabled), [model, enabled]);
	const onFieldByCol = useMemo(
		() => cols.map((c) => onFieldLaneAt(model.onField, c.frame)),
		[cols, model.onField],
	);
	const runs = useMemo(
		() => onFieldRuns(cols, onFieldByCol),
		[cols, onFieldByCol],
	);

	const wholeLogRows = useMemo(
		() =>
			model.lanes.map((l) =>
				cols.reduce((a, c) => Math.max(a, c.cells[l.index].length), 0),
			),
		[cols, model.lanes],
	);
	const onScreen = cols.slice(
		lowerBound(cols, view.left - COL_W + 1, (c) => c.x),
		lowerBound(cols, view.left + view.width - gutter, (c) => c.x),
	);
	const laneRows = model.lanes.map((l) =>
		expanded.has(l.index)
			? onScreen.reduce((a, c) => Math.max(a, c.cells[l.index].length), 0)
			: Math.min(wholeLogRows[l.index], cap),
	);
	const laneH = laneHeights(laneRows, portrait);
	const contentWidth = gutter + (cols[cols.length - 1]?.x ?? 0) + COL_W + 40;
	const totalH = HEAD_H + laneH.reduce((a, b) => a + b, 0);

	const px = visibleWindow(view.left, view.width);
	const from = lowerBound(cols, px[0] - COL_W, (c) => c.x);
	const shown = cols.slice(
		from,
		lowerBound(cols, px[1], (c) => c.x),
	);

	const firstIdx = lowerBound(cols, view.left, (c) => c.x);
	const lastIdx = Math.max(
		firstIdx,
		lowerBound(cols, view.left + view.width - gutter - COL_W, (c) => c.x),
	);
	const range: [number, number] =
		cols.length === 0
			? [0, model.maxFrame]
			: [
					cols[Math.min(firstIdx, cols.length - 1)].frame,
					cols[Math.min(lastIdx, cols.length - 1)].frame + 1,
				];

	const scrollToFrame = (frame: number, center: boolean) => {
		if (el == null || cols.length === 0) {
			return;
		}
		const i = Math.min(
			cols.length - 1,
			lowerBound(cols, frame, (c) => c.frame),
		);
		const left = center
			? cols[i].x - (el.clientWidth - gutter) / 2 + COL_W / 2
			: cols[i].x;
		el.scrollTo({ left, behavior: center ? "smooth" : "auto" });
	};
	useImperativeHandle(ref, () => ({
		centerOn: (frame) => scrollToFrame(frame, true),
	}));

	const [hl, setHl] = useState<Highlight | null>(null);
	const showDuration = (c: Chip) =>
		setHl((h) =>
			h?.event === c.event
				? null
				: {
						event: c.event,
						start: c.event.frame,
						end: c.event.end ?? c.event.frame,
						label: spanLabel(c.event),
					},
		);
	const [zoomed, setZoomed] = useState<number | null>(null);
	const [shownModel, setShownModel] = useState(model);
	if (model !== shownModel) {
		setShownModel(model);
		setHl(null);
		setZoomed(null);
		setExpanded(new Set());
	}
	useEffect(() => {
		if (hl == null || zoomed != null) {
			return;
		}
		const onKey = (e: KeyboardEvent) => {
			if (e.key === "Escape") {
				setHl(null);
			}
		};
		window.addEventListener("keydown", onKey);
		return () => window.removeEventListener("keydown", onKey);
	}, [hl, zoomed]);
	const hlI0 = hl == null ? 0 : lowerBound(cols, hl.start, (c) => c.frame);
	const hlI1 =
		hl == null ? -1 : lowerBound(cols, hl.end + 1, (c) => c.frame) - 1;
	const hlBand =
		hlI1 >= hlI0 ? { x0: cols[hlI0].x, x1: cols[hlI1].x + COL_W } : null;
	const hlTitle =
		hl == null
			? ""
			: t("sample.clear_duration", {
					label: hl.label,
					start: hl.start,
					end: hl.end > model.maxFrame ? t("sample.end") : hl.end,
				});

	const toggle = (i: number) =>
		setExpanded((s) => {
			const n = new Set(s);
			if (n.has(i)) {
				n.delete(i);
			} else {
				n.add(i);
			}
			return n;
		});

	const overviewLabel = topOpen
		? t("sample.hide_overview")
		: t("sample.show_overview");

	return (
		<>
			{topOpen && (
				<Minimap
					model={model}
					enabled={enabled}
					range={range}
					onPan={(f) => scrollToFrame(f, false)}
					matchFrames={hits.frames}
					aside={commandBar}
					selection={hl == null ? null : [hl.start, hl.end]}
				/>
			)}
			<div ref={setSlot} className="flex min-h-0 flex-1 flex-col">
				<div
					ref={scrollRef}
					className="relative min-h-0 overflow-auto [scrollbar-color:var(--g-text-mute)_transparent] [scrollbar-width:thin] rounded-g-md border border-g-line bg-g-surface-2"
				>
					<div
						className="relative"
						style={{ width: contentWidth, height: totalH }}
					>
						<div
							className="sticky top-0 z-20 border-b border-g-line bg-g-surface"
							style={{ height: HEAD_H, width: contentWidth }}
						>
							<div
								className="sticky left-0 z-10 flex h-full items-center border-r border-g-line bg-g-surface pr-1 pl-2 font-g-mono text-[10px] leading-tight text-g-ink-mute"
								style={{ width: gutter }}
							>
								{hl == null ? (
									<span className="flex min-w-0 flex-1 flex-col">
										<span>{t("sample.seconds")}</span>
										<span>{t("sample.frame")}</span>
									</span>
								) : (
									<button
										type="button"
										className="mr-1 flex min-w-0 flex-1 items-center gap-1 rounded-g-sm text-left text-g-warning hover:bg-g-surface-3"
										title={hlTitle}
										aria-label={hlTitle}
										onClick={() => setHl(null)}
									>
										<span className="flex min-w-0 flex-1 flex-col">
											<span className="truncate">{hl.label}</span>
											<span>
												{hl.end > model.maxFrame
													? t("sample.to_end")
													: secondsLabel(hl.end - hl.start)}
											</span>
										</span>
										<X className="size-3 shrink-0" />
									</button>
								)}
								<button
									type="button"
									className="flex size-6 shrink-0 items-center justify-center rounded-g-sm border border-g-line text-g-ink-dim hover:bg-g-surface-3 hover:text-g-ink"
									aria-label={overviewLabel}
									title={overviewLabel}
									aria-expanded={topOpen}
									onClick={() => setTopOpen((o) => !o)}
								>
									{topOpen ? (
										<ChevronsUp className="size-3.5" />
									) : (
										<ChevronsDown className="size-3.5" />
									)}
								</button>
							</div>
							{hlBand != null && (
								<div
									className="pointer-events-none absolute top-0 h-full border-b-[3px] border-g-warning bg-g-warning/15"
									style={{
										left: gutter + hlBand.x0,
										width: hlBand.x1 - hlBand.x0,
									}}
								/>
							)}
							{shown.map((c, k) => (
								<div key={c.frame}>
									{c.framesSincePrev > 1 && (
										<div
											className="absolute top-0 flex h-full items-center justify-center font-g-mono text-[9px] text-g-ink-mute"
											style={{ left: gutter + c.x - GAP_BIG, width: GAP_BIG }}
											title={t("sample.empty_frames", {
												count: c.framesSincePrev - 1,
											})}
										>
											<span className="-rotate-90 whitespace-nowrap">
												+{c.framesSincePrev}f
											</span>
										</div>
									)}
									<div
										className="absolute top-0 flex h-full items-center border-l border-g-line pr-1 pl-1.5 font-g-mono leading-tight"
										style={{ left: gutter + c.x, width: COL_W }}
									>
										<span className="flex min-w-0 flex-1 flex-col">
											<span className="text-[11px] text-g-ink">
												{secondsLabel(c.frame)}
											</span>
											<span className="text-[10px] text-g-ink-mute">
												{t("sample.frame_short", { frame: c.frame })}
											</span>
										</span>
										<button
											type="button"
											className="flex size-6 shrink-0 items-center justify-center rounded-g-sm text-g-ink-mute hover:bg-g-surface-3 hover:text-g-ink"
											aria-label={t("sample.open_frame", { frame: c.frame })}
											title={t("sample.open_frame", { frame: c.frame })}
											onClick={() => setZoomed(from + k)}
										>
											<Search className="size-3.5" />
										</button>
									</div>
								</div>
							))}
						</div>
						{hlBand != null && (
							<div
								className="pointer-events-none absolute border-x border-dashed border-g-warning/70 bg-g-warning/[0.06]"
								style={{
									top: HEAD_H,
									height: totalH - HEAD_H,
									left: gutter + hlBand.x0,
									width: hlBand.x1 - hlBand.x0,
								}}
							/>
						)}
						{model.lanes.map((lane) => {
							const rows = laneRows[lane.index];
							const isExpanded = expanded.has(lane.index);
							return (
								<div
									key={lane.key}
									className="relative border-b border-g-line"
									style={{ height: laneH[lane.index], width: contentWidth }}
								>
									<PortraitLabel
										lane={lane}
										name={laneName(lane)}
										size={portrait}
									>
										{isExpanded && (
											<button
												type="button"
												className="text-[10px] text-g-accent hover:underline"
												aria-expanded
												onClick={() => toggle(lane.index)}
											>
												▴ {t("sample.less")}
											</button>
										)}
									</PortraitLabel>
									{runs.map(
										(r) =>
											r.lane === lane.index &&
											r.x1 > px[0] - COL_W &&
											r.x0 < px[1] && (
												<div
													key={r.x0}
													className="pointer-events-none absolute inset-y-[3px] rounded-g-md border-2"
													style={{
														left: gutter + r.x0 + 2,
														width: r.x1 - r.x0 - 4,
														borderColor: laneColor(lane),
														background: laneTint(lane, 9),
													}}
												/>
											),
									)}
									{shown.map((c, k) => {
										const items = c.cells[lane.index];
										const active = onFieldByCol[from + k] === lane.index;
										const fit = chipsShown(items.length, rows, isExpanded);
										return (
											<div
												key={c.frame}
												className="absolute top-0 flex h-full flex-col gap-[2px] overflow-hidden"
												style={{
													left: gutter + c.x,
													width: COL_W,
													padding: CELL_PAD,
													borderLeft: active
														? undefined
														: "1px solid var(--g-border-soft)",
												}}
											>
												{items.slice(0, fit).map((e) => (
													<EventChip
														key={e.id}
														chip={e}
														onOpen={onOpen}
														onDurationIcon={showDuration}
														matched={hits.ids.has(e.id)}
														inDuration={hl?.event === e.event}
														className="w-full shrink-0"
													/>
												))}
												{fit < items.length && (
													<button
														type="button"
														className="h-[18px] shrink-0 rounded-[3px] border border-dashed border-g-line text-[10px] text-g-ink-dim hover:text-g-ink"
														aria-expanded={false}
														onClick={() => toggle(lane.index)}
													>
														{t("sample.more", { count: items.length - fit })}
													</button>
												)}
											</div>
										);
									})}
								</div>
							);
						})}
					</div>
				</div>
			</div>
			<FrameView
				model={model}
				cols={cols}
				index={zoomed}
				onIndex={(i) => {
					setZoomed(i);
					scrollToFrame(cols[i].frame, true);
				}}
				onClose={() => setZoomed(null)}
				matches={hits.ids}
				onOpen={onOpen}
				onDurationIcon={(c) => {
					setZoomed(null);
					showDuration(c);
				}}
			/>
		</>
	);
}

function PortraitLabel({
	lane,
	name,
	size,
	children,
}: {
	lane: Lane;
	name: string;
	size: number;
	children?: ReactNode;
}) {
	return (
		<div
			className={cn(
				"sticky left-0 z-10 flex h-full flex-col items-center gap-1 border-r border-g-line bg-g-surface-2 px-1.5",
				isSimLane(lane) ? "justify-center" : "pt-1",
			)}
			style={{ width: gutterWidth(size) }}
			title={name}
		>
			<LaneIcon lane={lane} name={name} style={{ width: size, height: size }} />
			{children}
		</div>
	);
}
