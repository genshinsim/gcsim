import { cn } from "@gcsim/primitives";
import { PanelRightClose, PanelRightOpen, Search } from "lucide-react";
import {
	type ReactNode,
	type Ref,
	useCallback,
	useEffect,
	useImperativeHandle,
	useMemo,
	useState,
} from "react";
import { useTranslation } from "react-i18next";
import type { SimEvent } from "../events/types";
import type { SampleView } from "../sampleView";
import { DurationClear, spanOf, useClearOnEscape } from "./duration";
import { EventChip } from "./EventChip";
import { FrameView } from "./FrameView";
import { LaneIcon, laneColor, laneTint, useLaneName } from "./lane";
import { columns, secondsLabel } from "./layout";
import { Minimap } from "./Minimap";
import {
	type Chip,
	chipShortText,
	lowerBound,
	onFieldLaneAt,
	type TimelineModel,
} from "./model";
import type { SearchHits, StripHandle } from "./Strip";
import {
	chipTextWidth,
	frameOfY,
	rowsHeight,
	slimChip,
	V_CELL_PAD_X,
	V_CELL_PAD_Y,
	V_GUTTER_W,
	V_HEAD_H,
	V_MIN_COL_W,
	type VRow,
	vColumnWidth,
	vRows,
	vShown,
	wrappedChip,
	yOfFrame,
} from "./vlayout";

const OVERSCAN_PX = 600;
const LINE_PX = 40;

type ScrollView = { top: number; height: number; width: number };

function useScrollBox() {
	const [el, setEl] = useState<HTMLDivElement | null>(null);
	const [view, setView] = useState<ScrollView>({
		top: 0,
		height: 800,
		width: 0,
	});
	useEffect(() => {
		if (el == null) {
			return;
		}
		const on = () =>
			setView((v) =>
				v.top === el.scrollTop &&
				v.height === el.clientHeight &&
				v.width === el.clientWidth
					? v
					: {
							top: el.scrollTop,
							height: el.clientHeight,
							width: el.clientWidth,
						},
			);
		on();
		el.addEventListener("scroll", on, { passive: true });
		const ro = new ResizeObserver(on);
		ro.observe(el);
		return () => {
			el.removeEventListener("scroll", on);
			ro.disconnect();
		};
	}, [el]);
	return { el, ref: setEl, view };
}

/** the sample log with time running down and a column per lane */
export function VerticalStrip({
	ref,
	model,
	enabled,
	hits,
	onOpen,
	commandBar,
	view,
	onOverview,
}: {
	ref?: Ref<StripHandle>;
	model: TimelineModel;
	enabled: Set<string>;
	hits: SearchHits;
	onOpen: (c: Chip) => void;
	commandBar: ReactNode;
	view: SampleView;
	onOverview: (open: boolean) => void;
}) {
	const { t } = useTranslation();
	const laneName = useLaneName();
	const { el, ref: scrollRef, view: sv } = useScrollBox();
	const [expanded, setExpanded] = useState<Set<number>>(new Set());
	const [hlEvent, setHl] = useState<SimEvent | null>(null);
	const [zoomed, setZoomed] = useState<number | null>(null);

	const cols = useMemo(() => columns(model, enabled), [model, enabled]);
	const onField = useMemo(
		() => cols.map((c) => onFieldLaneAt(model.onField, c.frame)),
		[cols, model.onField],
	);
	const lanes = useMemo(() => model.lanes.map((l) => l.index), [model]);
	const cap = view.slim ? view.cap : null;

	const [shown, setShown] = useState({ model, cols, cap });
	if (shown.model !== model || shown.cols !== cols || shown.cap !== cap) {
		setShown({ model, cols, cap });
		setExpanded(new Set());
		if (shown.model !== model) {
			setHl(null);
			setZoomed(null);
		}
	}

	const rows = useMemo(() => {
		const font = el == null ? "" : `11px ${getComputedStyle(el).fontFamily}`;
		const chipH =
			cap != null || el == null
				? slimChip
				: wrappedChip(
						vColumnWidth(sv.width, lanes.length) - 2 * V_CELL_PAD_X - 1,
						chipShortText,
						chipTextWidth(font),
					);
		return vRows(cols, lanes, { cap, expanded, chipH });
	}, [cols, lanes, cap, expanded, el, sv.width]);

	const bodyH = sv.height - V_HEAD_H;
	const range: [number, number] = [
		frameOfY(rows, sv.top),
		frameOfY(rows, sv.top + bodyH),
	];
	const onScreen = rows.slice(
		lowerBound(rows, sv.top - OVERSCAN_PX, (r) => r.y + r.h),
		lowerBound(rows, sv.top + bodyH + OVERSCAN_PX, (r) => r.y),
	);

	const centerOn = (frame: number) =>
		el?.scrollTo({
			top: yOfFrame(rows, frame) - bodyH / 2,
			behavior: "smooth",
		});
	useImperativeHandle(ref, () => ({ centerOn }));

	const hl = spanOf(hlEvent);
	useClearOnEscape(
		hlEvent != null && zoomed == null,
		useCallback(() => setHl(null), []),
	);
	const showDuration = (c: Chip) =>
		setHl((h) => (h === c.event ? null : c.event));
	const hlI0 = hl == null ? 0 : lowerBound(rows, hl.start, (r) => r.col.frame);
	const hlI1 =
		hl == null ? -1 : lowerBound(rows, hl.end + 1, (r) => r.col.frame) - 1;

	const toggleRow = (i: number) =>
		setExpanded((s) => {
			const n = new Set(s);
			if (n.has(i)) {
				n.delete(i);
			} else {
				n.add(i);
			}
			return n;
		});

	const overviewLabel = view.overview
		? t("sample.hide_overview_short")
		: t("sample.show_overview_short");
	const overviewToggle = (
		<button
			type="button"
			className="flex size-6 shrink-0 items-center justify-center rounded-g-sm border border-g-line text-g-ink-dim hover:bg-g-surface-3 hover:text-g-ink"
			aria-label={overviewLabel}
			title={overviewLabel}
			aria-expanded={view.overview}
			onClick={() => onOverview(!view.overview)}
		>
			{view.overview ? (
				<PanelRightClose className="size-3.5" />
			) : (
				<PanelRightOpen className="size-3.5" />
			)}
		</button>
	);

	return (
		<>
			<div className="shrink-0 rounded-g-md border border-g-line bg-g-surface">
				{commandBar}
			</div>
			<div className="flex min-h-0 flex-1 gap-2">
				<div
					ref={scrollRef}
					className="relative min-h-0 min-w-0 flex-1 overflow-auto [scrollbar-color:var(--g-text-mute)_transparent] [scrollbar-width:thin] rounded-g-md border border-g-line bg-g-surface-2"
				>
					<div
						className="relative"
						style={{
							minWidth: V_GUTTER_W + lanes.length * V_MIN_COL_W,
							height: V_HEAD_H + rowsHeight(rows),
						}}
					>
						<div
							className="sticky top-0 z-20 flex border-b border-g-line bg-g-surface"
							style={{ height: V_HEAD_H }}
						>
							<div
								className="sticky left-0 z-10 flex shrink-0 items-center border-r border-g-line bg-g-surface pr-1 pl-2 font-g-mono text-[10px] leading-tight text-g-ink-mute"
								style={{ width: V_GUTTER_W }}
							>
								{hl == null ? (
									<span className="flex flex-col">
										<span>{t("sample.seconds")}</span>
										<span>{t("sample.frame")}</span>
									</span>
								) : (
									<DurationClear span={hl} onClear={() => setHl(null)} />
								)}
							</div>
							{model.lanes.map((lane) => (
								<div
									key={lane.key}
									className="flex min-w-0 flex-1 items-center gap-1.5 border-l border-g-line px-2"
									title={laneName(lane)}
								>
									<LaneIcon
										lane={lane}
										name={laneName(lane)}
										style={{ width: 28, height: 28 }}
									/>
									<span className="truncate text-[12px] text-g-ink-dim">
										{laneName(lane)}
									</span>
								</div>
							))}
						</div>
						{onScreen.map((row) => (
							<Row
								key={row.col.frame}
								row={row}
								model={model}
								onField={onField}
								cap={cap}
								expanded={expanded.has(row.i)}
								onToggle={() => toggleRow(row.i)}
								onZoom={() => setZoomed(row.i)}
								inDuration={
									hl != null &&
									row.col.frame >= hl.start &&
									row.col.frame <= hl.end
								}
								chip={(e) => (
									<EventChip
										key={e.id}
										chip={e}
										onOpen={onOpen}
										onDurationIcon={showDuration}
										matched={hits.ids.has(e.id)}
										inDuration={hlEvent === e.event}
										className="w-full shrink-0"
										wrap={cap == null}
									/>
								)}
								hidesMatch={(cs) => cs.some((e) => hits.ids.has(e.id))}
							/>
						))}
						{hlI1 >= hlI0 && (
							<div
								className="pointer-events-none absolute right-0 z-10 border-y border-dashed border-g-warning/70 bg-g-warning/[0.06]"
								style={{
									left: V_GUTTER_W,
									top: V_HEAD_H + rows[hlI0].y + rows[hlI0].gap,
									height:
										rows[hlI1].y + rows[hlI1].h - rows[hlI0].y - rows[hlI0].gap,
								}}
							/>
						)}
					</div>
				</div>
				{view.overview ? (
					<Minimap
						vertical
						model={model}
						enabled={enabled}
						range={range}
						onPan={(f) => el?.scrollTo({ top: yOfFrame(rows, f) })}
						onWheel={(delta, mode) => {
							if (el != null) {
								const px =
									mode === WheelEvent.DOM_DELTA_LINE
										? LINE_PX
										: mode === WheelEvent.DOM_DELTA_PAGE
											? el.clientHeight
											: 1;
								el.scrollTo({ top: el.scrollTop + delta * px });
							}
						}}
						matchFrames={hits.frames}
						selection={hl == null ? null : [hl.start, hl.end]}
						aside={overviewToggle}
					/>
				) : (
					<div className="flex w-8 shrink-0 flex-col items-center rounded-g-md border border-g-line bg-g-surface py-1">
						{overviewToggle}
					</div>
				)}
			</div>
			<FrameView
				model={model}
				cols={cols}
				index={zoomed}
				onIndex={(i) => {
					setZoomed(i);
					centerOn(cols[i].frame);
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

function Row({
	row,
	model,
	onField,
	cap,
	expanded,
	onToggle,
	onZoom,
	inDuration,
	chip,
	hidesMatch,
}: {
	row: VRow;
	model: TimelineModel;
	/** lane on field per row, -1 for none */
	onField: number[];
	cap: number | null;
	expanded: boolean;
	onToggle: () => void;
	onZoom: () => void;
	inDuration: boolean;
	chip: (c: Chip) => ReactNode;
	hidesMatch: (cs: Chip[]) => boolean;
}) {
	const { t } = useTranslation();
	const { col, gap, i } = row;
	return (
		<div
			className="absolute inset-x-0 flex border-b border-g-line"
			style={{ top: V_HEAD_H + row.y, height: row.h }}
		>
			<div
				className={cn(
					"sticky left-0 z-10 flex shrink-0 flex-col border-r border-g-line bg-g-surface",
					inDuration && "bg-g-warning/15",
				)}
				style={{ width: V_GUTTER_W }}
			>
				{gap > 0 && (
					<div
						className="flex shrink-0 items-center justify-center border-b border-dashed border-g-line font-g-mono text-[9px] text-g-ink-mute"
						style={{ height: gap }}
						title={t("sample.empty_frames", {
							count: col.framesSincePrev - 1,
						})}
					>
						+{col.framesSincePrev}f
					</div>
				)}
				<div className="flex items-start gap-1 pt-1 pr-1 pl-2 font-g-mono leading-tight">
					<span className="flex min-w-0 flex-1 flex-col">
						<span className="text-[11px] text-g-ink">
							{secondsLabel(col.frame)}
						</span>
						<span className="text-[10px] text-g-ink-mute">
							{t("sample.frame_short", { frame: col.frame })}
						</span>
					</span>
					<button
						type="button"
						className="flex size-5 shrink-0 items-center justify-center rounded-g-sm text-g-ink-mute hover:bg-g-surface-3 hover:text-g-ink"
						aria-label={t("sample.open_frame", { frame: col.frame })}
						title={t("sample.open_frame", { frame: col.frame })}
						onClick={onZoom}
					>
						<Search className="size-3" />
					</button>
				</div>
				{expanded && (
					<button
						type="button"
						className="mt-auto mb-1 text-[10px] text-g-accent hover:underline"
						aria-expanded
						onClick={onToggle}
					>
						▴ {t("sample.less")}
					</button>
				)}
			</div>
			{model.lanes.map((lane) => {
				const items = col.cells[lane.index];
				const fit = vShown(items.length, cap, expanded);
				const active = onField[i] === lane.index;
				const starts = active && onField[i - 1] !== lane.index;
				const ends = active && onField[i + 1] !== lane.index;
				return (
					<div
						key={lane.key}
						className="relative min-w-0 flex-1 border-l border-g-line-soft"
					>
						{active && (
							<div
								className={cn(
									"pointer-events-none absolute inset-x-[2px] border-x-2",
									starts && "rounded-t-g-md border-t-2",
									ends && "rounded-b-g-md border-b-2",
								)}
								style={{
									top: starts ? gap + 2 : 0,
									bottom: ends ? 2 : -1,
									borderColor: laneColor(lane),
									background: laneTint(lane, 9),
								}}
							/>
						)}
						<div
							className="relative flex flex-col gap-[2px]"
							style={{
								paddingTop: gap + V_CELL_PAD_Y,
								paddingInline: V_CELL_PAD_X,
							}}
						>
							{items.slice(0, fit).map(chip)}
							{fit < items.length && (
								<button
									type="button"
									className={cn(
										"h-[18px] shrink-0 rounded-[3px] border border-dashed border-g-line text-[10px] text-g-ink-dim hover:text-g-ink",
										hidesMatch(items.slice(fit)) &&
											"border-solid border-g-warning text-g-warning",
									)}
									aria-expanded={false}
									onClick={onToggle}
								>
									{t("sample.more", { count: items.length - fit })}
								</button>
							)}
						</div>
					</div>
				);
			})}
		</div>
	);
}
