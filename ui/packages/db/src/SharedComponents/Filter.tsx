import tagData from "@gcsim/data/src/tags.json";
import { dynamicKey } from "@gcsim/localization";
import {
	Badge,
	Button,
	Input,
	Separator,
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@gcsim/primitives";
import { useContext, useState } from "react";
import { useTranslation } from "react-i18next";
import {
	FaArrowDown,
	FaArrowUp,
	FaBan,
	FaCheck,
	FaFilter,
	FaSearch,
} from "react-icons/fa";

import {
	charNames,
	FilterContext,
	FilterDispatchContext,
	type FilterState,
	ItemFilterState,
	SortByDirection,
	sortByParams,
	tagBaseState,
} from "./FilterComponents/Filter.utils";

function activeCount(filter: FilterState): number {
	let n = 0;
	for (const c of Object.values(filter.charFilter)) {
		if (c.state !== ItemFilterState.none) n++;
	}
	for (const key of Object.keys(filter.tagFilter)) {
		if (filter.tagFilter[key].state !== tagBaseState(key)) n++;
	}
	return n;
}

export function Filter() {
	// https://github.com/i18next/next-i18next/issues/1795
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const filter = useContext(FilterContext);
	const [isOpen, setIsOpen] = useState(false);
	const n = activeCount(filter);

	return (
		<Sheet open={isOpen} onOpenChange={setIsOpen}>
			<SheetTrigger asChild>
				<Button
					variant="outline"
					className="gap-g-base-sm"
					aria-label={t("db.filter")}
				>
					<FaFilter size={13} className="text-g-accent" />
					{t("db.filter")}
					{n > 0 && (
						<Badge className="bg-g-accent-weak text-g-accent">{n}</Badge>
					)}
				</Button>
			</SheetTrigger>
			<SheetContent side="left" className="w-[min(92vw,360px)] overflow-y-auto">
				<SheetHeader>
					<div className="flex flex-row items-center justify-between pr-6">
						<SheetTitle className="text-g-h3">{t("db.filter")}</SheetTitle>
						<ClearFilterButton />
					</div>
					<SheetDescription className="sr-only">
						{t("db.filter")}
					</SheetDescription>
				</SheetHeader>
				<Separator />
				<div className="flex flex-col gap-g-section p-2">
					<Section title={t("db.sort_by")}>
						<SortControl />
					</Section>
					<Section title={t("db.tags")}>
						<TagPicker />
					</Section>
					<Section title={t("db.characters")}>
						<CharacterPicker />
					</Section>
				</div>
			</SheetContent>
		</Sheet>
	);
}

function ClearFilterButton() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const dispatch = useContext(FilterDispatchContext);
	return (
		<Button
			variant="ghost"
			size="sm"
			onClick={() => dispatch({ type: "clearFilter" })}
		>
			{t("db.clear")}
		</Button>
	);
}

function Section({
	title,
	children,
}: {
	title: string;
	children: React.ReactNode;
}) {
	return (
		<div className="flex flex-col gap-g-base-sm">
			<div className="text-g-xs font-semibold uppercase tracking-wide text-g-ink-mute">
				{title}
			</div>
			{children}
		</div>
	);
}

function SortControl() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;

	return (
		<div className="flex flex-wrap gap-g-base-sm">
			{sortByParams.map((param) => (
				<SortByParamButton
					key={param.sortKey}
					sortKey={param.sortKey}
					translation={t(param.translationKey)}
				/>
			))}
		</div>
	);
}

function SortByParamButton({
	sortKey,
	translation,
}: {
	sortKey: string;
	translation: string;
}) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const active = filter.sortBy.sortKey === sortKey;
	const direction = active ? filter.sortBy.sortByDirection : null;

	return (
		<Button
			size="sm"
			variant={active ? "default" : "secondary"}
			onClick={() => dispatch({ type: "handleSortBy", sortByKey: sortKey })}
		>
			{translation}
			{direction === SortByDirection.asc && <FaArrowUp size={10} />}
			{direction === SortByDirection.dsc && <FaArrowDown size={10} />}
		</Button>
	);
}

function TagPicker() {
	const tags = Object.keys(tagData)
		.filter((key) => key !== "0" && key !== "1" && key !== "2")
		.map((key) => ({ key, name: tagData[key].display_name }));

	return (
		<div className="flex flex-wrap gap-g-base-sm">
			{tags.map((tag) => (
				<TagPill key={tag.key} name={tag.name} tag={tag.key} />
			))}
		</div>
	);
}

function TagPill({ tag, name }: { tag: string; name: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const state = filter.tagFilter[tag].state;
	const cls =
		state === ItemFilterState.include
			? "border-transparent bg-g-success/20 text-g-success"
			: state === ItemFilterState.exclude
				? "border-transparent bg-g-danger/20 text-g-danger"
				: "border-g-line-soft bg-g-surface-2 text-g-ink-dim hover:border-g-line";
	return (
		<button
			type="button"
			onClick={() => dispatch({ type: "handleTag", tag })}
			className={`inline-flex items-center gap-1 rounded-g-pill border px-2.5 py-1 text-g-xs font-medium transition-colors ${cls}`}
		>
			{state === ItemFilterState.include && <FaCheck size={9} />}
			{state === ItemFilterState.exclude && <FaBan size={9} />}
			{name}
		</button>
	);
}

function CharacterPicker() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const translateCharName = (charName: string) =>
		t(`game:character_names.${charName}`);

	const [charSearch, setCharSearch] = useState<string>("");
	const sorted = [...charNames].sort((a, b) =>
		translateCharName(a).localeCompare(translateCharName(b)),
	);
	const visible = sorted.filter((charName) =>
		translateCharName(charName)
			.toLocaleLowerCase()
			.includes(charSearch.toLocaleLowerCase()),
	);

	return (
		<div className="flex flex-col gap-g-base">
			<div className="relative">
				<FaSearch
					className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-g-ink-mute"
					size={13}
				/>
				<Input
					className="pl-8"
					type="text"
					dir="auto"
					placeholder={t("db.type_to_search")}
					value={charSearch}
					onChange={(e) => setCharSearch(e.target.value)}
				/>
			</div>
			<div className="grid max-h-[46vh] grid-cols-4 gap-g-base-sm overflow-y-auto overflow-x-hidden no-scrollbar">
				{visible.map((charName) => (
					<CharCard
						key={charName}
						charName={charName}
						label={translateCharName(charName)}
					/>
				))}
			</div>
		</div>
	);
}

function CharCard({ charName, label }: { charName: string; label: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const state = filter.charFilter[charName].state;
	const set = state !== ItemFilterState.none;
	return (
		<button
			type="button"
			title={label}
			onClick={() => dispatch({ type: "handleChar", char: charName })}
			className={`relative flex flex-col items-center rounded-g-md border p-1 transition-colors ${
				set
					? "border-g-accent bg-g-accent-weak"
					: "border-g-line-soft bg-g-surface-2 hover:border-g-line"
			}`}
		>
			<img
				src={`/api/assets/avatar/${charName}.png`}
				alt=""
				className="h-12 w-12 object-contain"
			/>
			<span className="w-full truncate text-center text-g-xs text-g-ink-dim">
				{label}
			</span>
			<TriRing state={state} />
		</button>
	);
}

function TriRing({ state }: { state: ItemFilterState }) {
	if (state === ItemFilterState.none) return null;
	const include = state === ItemFilterState.include;
	const cls = include
		? "bg-g-success text-g-accent-fg"
		: "bg-g-danger text-g-accent-fg";
	return (
		<span
			className={`absolute right-1 top-1 flex size-4 items-center justify-center rounded-full ${cls}`}
		>
			{include ? <FaCheck size={8} /> : <FaBan size={8} />}
		</span>
	);
}
