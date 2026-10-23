import {
	Button,
	cn,
	Input,
	Label,
	Popover,
	PopoverContent,
	PopoverTrigger,
	Switch,
} from "@gcsim/primitives";
import {
	ArrowRight,
	Columns3,
	Download,
	ListFilter,
	RefreshCw,
	Rows3,
} from "lucide-react";
import { useId } from "react";
import { useTranslation } from "react-i18next";
import {
	AdvancedPreset,
	AllSampleOptions,
	DebugPreset,
	SimplePreset,
	VerbosePreset,
} from "../SampleOptions";
import {
	CAP_MAX,
	CAP_MIN,
	type Orientation,
	type SampleView,
} from "../sampleView";
import { display } from "./display";
import type { TimelineModel } from "./model";

type Props = {
	model: TimelineModel;
	filter: string[];
	onFilterChange: (v: string[]) => void;
	search: string;
	onSearch: (v: string) => void;
	onSearchNext: () => void;
	matchLabel: string;
	onDownload?: () => void;
	onGenerate?: () => void;
	view: SampleView;
	onViewChange: (v: SampleView) => void;
	/** one wrapping row, for a bar across the top rather than a panel beside the overview */
	inline?: boolean;
};

const PRESETS = [
	["viewer.simple", SimplePreset],
	["viewer.advanced", AdvancedPreset],
	["viewer.verbose", VerbosePreset],
	["viewer.debug", DebugPreset],
	["sample.all", AllSampleOptions],
	["viewer.clear", []],
] as const;

export function CommandBar({
	model,
	filter,
	onFilterChange,
	search,
	onSearch,
	onSearchNext,
	matchLabel,
	onDownload,
	onGenerate,
	view,
	onViewChange,
	inline = false,
}: Props) {
	const { t } = useTranslation();
	const group = inline ? "contents" : "flex items-center gap-1";
	return (
		<div
			className={
				inline
					? "flex flex-wrap items-center gap-1 p-1.5"
					: "flex h-full flex-col justify-center gap-1.5 px-2"
			}
		>
			<div className={group}>
				<Input
					className={cn(
						"min-w-0 flex-1 text-[12px]",
						inline && "min-w-[140px] max-w-sm",
					)}
					placeholder={t("sample.search_placeholder")}
					aria-label={t("viewer.search")}
					value={search}
					onChange={(e) => onSearch(e.target.value)}
					onKeyDown={(e) => {
						if (e.key === "Enter") {
							onSearchNext();
						}
					}}
				/>
				<Button
					variant="secondary"
					size="icon"
					aria-label={t("sample.next_match")}
					title={t("sample.next_match")}
					onClick={onSearchNext}
				>
					<ArrowRight />
				</Button>
			</div>
			<div className={group}>
				<span
					className={cn(
						"min-w-0 truncate font-g-mono text-[11px] text-g-ink-mute",
						inline ? "mr-auto px-1" : "flex-1",
					)}
				>
					{matchLabel}
				</span>
				<CategoryFilter
					model={model}
					filter={filter}
					onFilterChange={onFilterChange}
				/>
				<ViewMenu view={view} onViewChange={onViewChange} />
				{onGenerate != null && (
					<Button
						variant="secondary"
						size="icon"
						aria-label={t("viewer.generate")}
						title={t("viewer.generate")}
						onClick={onGenerate}
					>
						<RefreshCw />
					</Button>
				)}
				{onDownload != null && (
					<Button
						variant="secondary"
						size="icon"
						aria-label={t("viewer.download")}
						title={t("viewer.download")}
						onClick={onDownload}
					>
						<Download />
					</Button>
				)}
			</div>
		</div>
	);
}

function CategoryFilter({
	model,
	filter,
	onFilterChange,
}: {
	model: TimelineModel;
	filter: string[];
	onFilterChange: (v: string[]) => void;
}) {
	const { t } = useTranslation();
	const cats = AllSampleOptions.filter((c) => (model.counts.get(c) ?? 0) > 0);
	const toggle = (c: string) =>
		onFilterChange(
			filter.includes(c) ? filter.filter((s) => s !== c) : [...filter, c],
		);
	const shownCount = cats.filter((c) => filter.includes(c)).length;
	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					variant="secondary"
					className="gap-1"
					title={t("viewer.log_options")}
					aria-label={t("viewer.log_options")}
				>
					<ListFilter />
					<span className="font-g-mono text-[11px] text-g-ink-mute">
						{shownCount}/{cats.length}
					</span>
				</Button>
			</PopoverTrigger>
			<PopoverContent
				align="end"
				className="w-[min(440px,calc(100vw-2rem))] p-3"
			>
				<div className="mb-2 flex flex-wrap items-center gap-1">
					<span className="mr-1 text-[11px] text-g-ink-mute">
						{t("sample.preset")}
					</span>
					{PRESETS.map(([label, p]) => (
						<Button
							key={label}
							size="sm"
							variant="ghost"
							className="h-7 px-2 text-[11px]"
							onClick={() => onFilterChange([...p])}
						>
							{t(label)}
						</Button>
					))}
				</div>
				<div className="flex flex-wrap gap-1">
					{cats.map((c) => {
						const on = filter.includes(c);
						const { color } = display(c);
						return (
							<button
								type="button"
								key={c}
								aria-pressed={on}
								onClick={() => toggle(c)}
								className={cn(
									"flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-[11px] transition-opacity",
									on
										? "border-g-line bg-g-surface-3 text-g-ink"
										: "border-dashed border-g-line text-g-ink-mute opacity-60",
								)}
							>
								<span
									className="size-2 rounded-full"
									style={{
										background: on ? color : "transparent",
										outline: `1px solid ${color}`,
									}}
								/>
								{c}
								<span className="font-g-mono text-g-ink-mute">
									{model.counts.get(c)}
								</span>
							</button>
						);
					})}
				</div>
			</PopoverContent>
		</Popover>
	);
}

const ORIENTATIONS: [Orientation, typeof Rows3][] = [
	["horizontal", Columns3],
	["vertical", Rows3],
];

function ViewMenu({
	view,
	onViewChange,
}: {
	view: SampleView;
	onViewChange: (v: SampleView) => void;
}) {
	const { t } = useTranslation();
	const id = useId();
	const set = (p: Partial<SampleView>) => onViewChange({ ...view, ...p });
	const vertical = view.orientation === "vertical";
	const Icon = vertical ? Rows3 : Columns3;
	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					variant="secondary"
					size="icon"
					title={t("sample.view")}
					aria-label={t("sample.view")}
				>
					<Icon />
				</Button>
			</PopoverTrigger>
			<PopoverContent align="end" className="w-64 p-3 text-[12px]">
				<div className="mb-1 text-[11px] text-g-ink-mute">
					{t("sample.layout")}
				</div>
				<div className="grid grid-cols-2 gap-1">
					{ORIENTATIONS.map(([o, OIcon]) => (
						<Button
							key={o}
							size="sm"
							variant={view.orientation === o ? "default" : "secondary"}
							aria-pressed={view.orientation === o}
							onClick={() => set({ orientation: o })}
						>
							<OIcon />
							{t(`sample.${o}`)}
						</Button>
					))}
				</div>
				<fieldset
					disabled={!vertical}
					className="mt-3 flex flex-col gap-2 border-t border-g-line pt-3 disabled:opacity-50"
				>
					<div className="flex items-center gap-2">
						<Switch
							id={`${id}-slim`}
							checked={view.slim}
							onCheckedChange={(slim) => set({ slim })}
						/>
						<Label htmlFor={`${id}-slim`} className="font-normal">
							{t("sample.slim_chips")}
						</Label>
					</div>
					<label
						className={cn(
							"flex flex-col gap-1 text-g-ink-dim",
							!view.slim && "opacity-50",
						)}
					>
						{t("sample.chip_cap", { count: view.cap })}
						<input
							type="range"
							min={CAP_MIN}
							max={CAP_MAX}
							step={1}
							value={view.cap}
							disabled={!view.slim}
							style={{ accentColor: "var(--g-accent)" }}
							onChange={(e) => set({ cap: Number(e.target.value) })}
						/>
					</label>
					<p className="text-[11px] text-g-ink-mute">
						{vertical ? t("sample.slim_chips_hint") : t("sample.vertical_only")}
					</p>
				</fieldset>
			</PopoverContent>
		</Popover>
	);
}
