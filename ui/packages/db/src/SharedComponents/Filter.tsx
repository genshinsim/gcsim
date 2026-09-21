import tagData from "@gcsim/data/src/tags.json";
import { dynamicKey } from "@gcsim/localization";
import {
	Button,
	Collapsible,
	CollapsibleContent,
	Input,
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@gcsim/primitives";
import { useContext, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaArrowDown, FaArrowUp, FaFilter, FaSearch } from "react-icons/fa";

import {
	charNames,
	FilterContext,
	FilterDispatchContext,
	ItemFilterState,
	SortByDirection,
	sortByParams,
} from "./FilterComponents/Filter.utils";

const activeFilterClasses = "bg-g-accent text-g-accent-fg hover:bg-g-accent/90";

export function Filter() {
	// https://github.com/i18next/next-i18next/issues/1795
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;

	const [isOpen, setIsOpen] = useState(false);

	return (
		<Sheet open={isOpen} onOpenChange={setIsOpen}>
			<SheetTrigger asChild>
				<Button variant="outline" aria-label={t("db.filter")}>
					<FaFilter size={13} className="opacity-80" />
					{t("db.filter")}
				</Button>
			</SheetTrigger>
			<SheetContent side="left" className="overflow-y-auto">
				<SheetHeader>
					<div className="flex flex-row justify-between pr-6">
						<SheetTitle className="text-g-h3">{t("db.filter")}</SheetTitle>
						<ClearFilterButton />
					</div>
					<SheetDescription className="sr-only">
						{t("db.filter")}
					</SheetDescription>
				</SheetHeader>
				<div className="flex flex-col gap-g-base-sm overflow-y-auto overflow-x-hidden p-2">
					<CharacterFilter />
					<TagFilter />
					<SortBy />
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
			variant="destructive"
			size="sm"
			onClick={() => dispatch({ type: "clearFilter" })}
		>
			{t("db.clear")}
		</Button>
	);
}

function FilterSection({
	label,
	children,
}: {
	label: string;
	children: React.ReactNode;
}) {
	const [isOpen, setIsOpen] = useState(false);
	return (
		<div className="w-full overflow-x-hidden no-scrollbar">
			<Button
				variant="secondary"
				className="w-full justify-between"
				onClick={() => setIsOpen(!isOpen)}
			>
				<span className="grow text-left">{label}</span>
				<span>{isOpen ? "-" : "+"}</span>
			</Button>
			<Collapsible open={isOpen}>
				<CollapsibleContent>
					<div className="mt-2 rounded-g-md bg-g-surface-2 p-1">{children}</div>
				</CollapsibleContent>
			</Collapsible>
		</div>
	);
}

function TagFilter() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const sortedTagnames = Object.keys(tagData)
		.filter((key) => key !== "0" && key !== "1" && key !== "2")
		.map((key) => ({ key, name: tagData[key]["display_name"] }));

	return (
		<FilterSection label={t("db.tags")}>
			<div className="grid grid-cols-3 gap-2">
				{sortedTagnames.map((tag) => (
					<TagFilterButton key={tag.key} name={tag.name} tag={tag.key} />
				))}
			</div>
		</FilterSection>
	);
}

function TagFilterButton({ tag, name }: { tag: string; name: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const state = filter.tagFilter[tag].state;
	return (
		<Button
			variant={state === ItemFilterState.exclude ? "destructive" : "secondary"}
			className={
				state === ItemFilterState.include ? activeFilterClasses : undefined
			}
			onClick={() => dispatch({ type: "handleTag", tag })}
		>
			<div className="text-center">{name}</div>
		</Button>
	);
}

function CharacterFilter() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const sortedCharNames = charNames.sort((a, b) => {
		if (t(a) < t(b)) return -1;
		if (t(a) > t(b)) return 1;
		return 0;
	});
	const [charSearch, setCharSearch] = useState<string>("");

	const translateCharName = (charName: string) =>
		t("game:character_names." + charName);

	return (
		<FilterSection label={t("db.characters")}>
			<div className="flex flex-col">
				<div className="relative flex flex-row text-g-ink-mute">
					<FaSearch className="pointer-events-none absolute right-2 top-3 size-4" />
					<Input
						className="grow"
						type="text"
						dir="auto"
						placeholder={t("db.type_to_search")}
						onChange={(e) => setCharSearch(e.target.value)}
					/>
				</div>

				<div className="mt-1 grid grid-cols-4 gap-1 overflow-y-auto overflow-x-hidden">
					{sortedCharNames
						.filter((charName) =>
							translateCharName(charName)
								.toLocaleLowerCase()
								.includes(charSearch.toLocaleLowerCase()),
						)
						.map((charName) => (
							<CharFilterButton key={charName} charName={charName} />
						))}
				</div>
			</div>
		</FilterSection>
	);
}

function CharFilterButton({ charName }: { charName: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const state = filter.charFilter[charName].state;
	return (
		<Button
			variant={
				state === ItemFilterState.exclude
					? "destructive"
					: state === ItemFilterState.include
						? "default"
						: "secondary"
			}
			className={`block h-auto ${
				state === ItemFilterState.include ? activeFilterClasses : ""
			}`}
			onClick={() => dispatch({ type: "handleChar", char: charName })}
		>
			<CharFilterButtonChild charName={charName} />
		</Button>
	);
}

function CharFilterButtonChild({ charName }: { charName: string }) {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const displayCharName = t("game:character_names." + charName);

	const travelerName = (
		charName.includes("lumine") || charName.includes("aether")
			? displayCharName
			: ""
	).replace(/.*?\((\S+)\).*?/, "$1");

	return (
		<div className="flex flex-col truncate gap-1">
			<img
				alt={displayCharName}
				src={`/api/assets/avatar/${charName}.png`}
				className="h-16 truncate object-contain"
			/>
			{travelerName !== "" ? (
				<div className="text-center">{travelerName}</div>
			) : null}
		</div>
	);
}

function SortBy() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;

	return (
		<FilterSection label={t("db.sort_by")}>
			<div className="flex flex-row flex-wrap gap-g-base">
				{sortByParams.map((param) => (
					<SortByParamButton
						key={param.sortKey}
						sortKey={param.sortKey}
						translation={t(param.translationKey)}
					/>
				))}
			</div>
		</FilterSection>
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
			onClick={() => dispatch({ type: "handleSortBy", sortByKey: sortKey })}
			variant={direction === SortByDirection.dsc ? "destructive" : "secondary"}
			className={
				direction === SortByDirection.asc ? activeFilterClasses : undefined
			}
		>
			<div className="flex flex-row items-center justify-center gap-1">
				{direction === SortByDirection.asc && <FaArrowUp />}
				{direction === SortByDirection.dsc && <FaArrowDown />}
				{translation}
			</div>
		</Button>
	);
}
