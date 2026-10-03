import {
	Button,
	cn,
	Input,
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@gcsim/primitives";
import { ArrowRight, Download, ListFilter } from "lucide-react";
import { useTranslation } from "react-i18next";
import {
	AdvancedPreset,
	AllSampleOptions,
	DebugPreset,
	SimplePreset,
	VerbosePreset,
} from "../SampleOptions";
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
}: Props) {
	const { t } = useTranslation();
	return (
		<div className="flex h-full flex-col justify-center gap-1.5 px-2">
			<div className="flex items-center gap-1">
				<Input
					className="h-7 min-w-0 flex-1 text-[12px]"
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
					className="size-7"
					aria-label={t("sample.next_match")}
					title={t("sample.next_match")}
					onClick={onSearchNext}
				>
					<ArrowRight />
				</Button>
			</div>
			<div className="flex items-center gap-1">
				<span className="min-w-0 flex-1 truncate font-g-mono text-[11px] text-g-ink-mute">
					{matchLabel}
				</span>
				<CategoryFilter
					model={model}
					filter={filter}
					onFilterChange={onFilterChange}
				/>
				{onDownload != null && (
					<Button
						variant="secondary"
						size="icon"
						className="size-7"
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
					size="sm"
					className="h-7 gap-1 px-2"
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
