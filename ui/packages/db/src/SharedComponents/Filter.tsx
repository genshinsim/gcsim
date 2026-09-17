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
} from "@gcsim/primitives";
import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaArrowDown, FaArrowUp, FaFilter, FaSearch } from "react-icons/fa";
import useDebounce from "../SharedHooks/debounce";

import {
	charNames,
	FilterContext,
	FilterDispatchContext,
	ItemFilterState,
	SortByDirection,
	sortByParams,
} from "./FilterComponents/Filter.utils";

const activeFilterClasses = "bg-emerald-600 text-white hover:bg-emerald-600/90";

export function Filter() {
	// https://github.com/i18next/next-i18next/issues/1795
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;

	const dispatch = useContext(FilterDispatchContext);
	const [isOpen, setIsOpen] = useState(false);

	const [value, setValue] = useState<string>("");
	const debouncedValue = useDebounce<string>(value, 500);

	useEffect(() => {
		dispatch({ type: "setCustomFilter", customFilter: debouncedValue });
	}, [debouncedValue, dispatch]);

	return (
		<div>
			<Button
				className="h-12 w-12 p-3"
				onClick={() => setIsOpen(true)}
				aria-label={t("db.filter")}
			>
				<FaFilter size={24} className="opacity-80" />
			</Button>

			<Sheet open={isOpen} onOpenChange={setIsOpen}>
				<SheetContent side="left" className="overflow-y-auto">
					<SheetHeader>
						<div className="flex flex-row justify-between pr-6">
							<SheetTitle className="text-xl">{t("db.filter")}</SheetTitle>
							<ClearFilterButton />
						</div>
						<SheetDescription className="sr-only">
							{t("db.filter")}
						</SheetDescription>
					</SheetHeader>
					<div className="flex flex-col gap-2 overflow-y-auto overflow-x-hidden p-2">
						<Input
							placeholder={t("db.customFilter")}
							type="text"
							dir="auto"
							onChange={(e) => {
								setValue(e.target.value);
							}}
						/>
						<CharacterFilter />
						<TagFilter />
						<SortBy />
					</div>
				</SheetContent>
			</Sheet>
		</div>
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

function TagFilter() {
	const [tagIsOpen, setTagIsOpen] = useState(false);
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const sortedTagnames = Object.keys(tagData)
		.filter((key) => {
			return key !== "0" && key !== "1" && key !== "2";
		})
		.map((key) => {
			return {
				key: key,
				name: tagData[key]["display_name"],
			};
		});

	return (
		<div className="w-full  overflow-x-hidden no-scrollbar">
			<Button
				className="w-full justify-between"
				onClick={() => setTagIsOpen(!tagIsOpen)}
			>
				<div className=" grow">{t("db.tags")}</div>

				<div className="">{tagIsOpen ? "-" : "+"}</div>
			</Button>
			<Collapsible open={tagIsOpen}>
				<CollapsibleContent>
					<div className="grid grid-cols-3 gap-2 mt-2 bg-gray-800 p-1">
						{sortedTagnames.map((t) => (
							<TagFilterButton key={t.key} name={t.name} tag={t.key} />
						))}
					</div>
				</CollapsibleContent>
			</Collapsible>
		</div>
	);
}

function TagFilterButton({ tag, name }: { tag; name: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const handleClick = () => {
		dispatch({
			type: "handleTag",
			tag: tag,
		});
	};

	const state = filter.tagFilter[tag].state;
	return (
		<Button
			variant={state === ItemFilterState.exclude ? "destructive" : "secondary"}
			className={
				state === ItemFilterState.include ? activeFilterClasses : undefined
			}
			onClick={handleClick}
		>
			<div className="text-center">{name}</div>
		</Button>
	);
}

function CharacterFilter() {
	const [charIsOpen, setCharIsOpen] = useState(false);
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const sortedCharNames = charNames.sort((a, b) => {
		if (t(a) < t(b)) {
			return -1;
		}
		if (t(a) > t(b)) {
			return 1;
		}
		return 0;
	});
	const [charSearch, setCharSearch] = useState<string>("");

	const translateCharName = (charName: string) =>
		t("game:character_names." + charName);

	return (
		<div className="w-full  overflow-x-hidden no-scrollbar">
			<Button
				className="w-full justify-between"
				onClick={() => setCharIsOpen(!charIsOpen)}
			>
				<div className=" grow">{t("db.characters")}</div>

				<div className="">{charIsOpen ? "-" : "+"}</div>
			</Button>
			<Collapsible open={charIsOpen}>
				<CollapsibleContent>
					<div className="flex flex-col mt-2 bg-gray-800 p-1">
						<label
							htmlFor="email"
							className="relative text-gray-400 focus-within:text-gray-600 flex flex-row"
						>
							<FaSearch className="pointer-events-none w-4 h-4 absolute top-2 transform   right-2 " />

							<Input
								className="grow"
								type="text"
								dir="auto"
								onChange={(e) => {
									setCharSearch(e.target.value);
								}}
							/>
						</label>

						<div className="grid grid-cols-4 gap-1 mt-1 overflow-y-auto overflow-x-hidden">
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
				</CollapsibleContent>
			</Collapsible>
		</div>
	);
}

function CharFilterButton({ charName }: { charName: string }) {
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const handleClick = () => {
		dispatch({
			type: "handleChar",
			char: charName,
		});
	};

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
			onClick={handleClick}
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
				className="truncate h-16 object-contain"
			/>
			{travelerName !== "" ? (
				<div className="text-center">{travelerName}</div>
			) : (
				<></>
			)}
		</div>
	);
}

function SortBy() {
	const [sortIsOpen, setSortIsOpen] = useState(false);
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;

	return (
		<div className="w-full  overflow-x-hidden no-scrollbar">
			<Button
				className="w-full justify-between"
				onClick={() => setSortIsOpen(!sortIsOpen)}
			>
				<div className=" grow">{t("db.sort_by")}</div>

				<div className="">{sortIsOpen ? "-" : "+"}</div>
			</Button>
			<Collapsible open={sortIsOpen}>
				<CollapsibleContent>
					<div className="flex flex-col mt-2 bg-gray-800 p-1">
						<div className="flex flex-row gap-4">
							{sortByParams.map((param) => (
								<SortByParamButton
									key={param.sortKey}
									sortKey={param.sortKey}
									translation={t(param.translationKey)}
								/>
							))}
						</div>
					</div>
				</CollapsibleContent>
			</Collapsible>
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

	const handleClick = () => {
		dispatch({
			type: "handleSortBy",
			sortByKey: sortKey,
		});
	};

	const active = filter.sortBy.sortKey === sortKey;
	const direction = active ? filter.sortBy.sortByDirection : null;

	return (
		<Button
			onClick={handleClick}
			variant={direction === SortByDirection.dsc ? "destructive" : "secondary"}
			className={
				direction === SortByDirection.asc ? activeFilterClasses : undefined
			}
		>
			<div className="flex flex-row gap-1 justify-center items-center">
				{direction === SortByDirection.asc && <FaArrowUp />}
				{direction === SortByDirection.dsc && <FaArrowDown />}
				{translation}
			</div>
		</Button>
	);
}
