import tagData from "@gcsim/data/src/tags.json";
import { dynamicKey } from "@gcsim/localization";
import { Button, Input } from "@gcsim/primitives";
import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaBan, FaSearch, FaTimes } from "react-icons/fa";
import useDebounce from "../SharedHooks/debounce";
import { Filter } from "./Filter";
import {
	FilterContext,
	FilterDispatchContext,
	ItemFilterState,
	tagBaseState,
} from "./FilterComponents/Filter.utils";

export function ActionBar({ simCount }: { simCount: number | null }) {
	const { t } = useTranslation();

	return (
		<div className="sticky top-0 z-10 -mx-8 flex flex-col gap-g-base border-b border-g-line-soft bg-g-canvas/90 px-8 py-g-base backdrop-blur">
			<div className="flex flex-wrap items-center gap-g-base">
				<Filter />
				<CustomFilterSearch />
				<span className="ml-auto whitespace-nowrap font-g-mono text-g-sm text-g-ink-dim">
					{t("db.showing_simulations", { i: simCount ?? 0 })}
				</span>
			</div>
			<ActiveFilterChips />
		</div>
	);
}

function CustomFilterSearch() {
	const { t } = useTranslation();
	const dispatch = useContext(FilterDispatchContext);
	const [value, setValue] = useState<string>("");
	const debouncedValue = useDebounce<string>(value, 500);

	useEffect(() => {
		dispatch({ type: "setCustomFilter", customFilter: debouncedValue });
	}, [debouncedValue, dispatch]);

	return (
		<div className="relative min-w-[180px] flex-1">
			<FaSearch
				className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-g-ink-mute"
				size={13}
			/>
			<Input
				className="pl-8"
				type="text"
				dir="auto"
				placeholder={t("db.search")}
				value={value}
				onChange={(e) => setValue(e.target.value)}
			/>
		</div>
	);
}

function ActiveFilterChips() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const chars = Object.values(filter.charFilter).filter(
		(c) => c.state !== ItemFilterState.none,
	);
	const tags = Object.values(filter.tagFilter).filter(
		(tg) => tg.state !== tagBaseState(tg.tag),
	);
	if (chars.length === 0 && tags.length === 0) return null;

	const clearAll = () => {
		for (const c of chars) dispatch({ type: "removeChar", char: c.charName });
		for (const tg of tags) dispatch({ type: "resetTag", tag: tg.tag });
	};

	return (
		<div className="flex flex-wrap items-center gap-g-base-sm">
			{chars.map((c) => {
				const excluded = c.state === ItemFilterState.exclude;
				return (
					<button
						type="button"
						key={`char-${c.charName}`}
						onClick={() => dispatch({ type: "removeChar", char: c.charName })}
						className={`inline-flex items-center gap-1.5 rounded-g-pill py-1 pl-1.5 pr-2 text-g-xs ${
							excluded
								? "bg-g-danger/20 text-g-danger line-through"
								: "bg-g-accent-weak text-g-accent"
						}`}
					>
						{excluded && <FaBan size={9} className="no-underline" />}
						<img
							src={`/api/assets/avatar/${c.charName}.png`}
							alt=""
							className="size-4"
						/>
						{t("game:character_names." + c.charName)}
						<FaTimes size={9} />
					</button>
				);
			})}
			{tags.map((tg) => {
				const excluded = tg.state === ItemFilterState.exclude;
				return (
					<button
						type="button"
						key={`tag-${tg.tag}`}
						onClick={() => dispatch({ type: "resetTag", tag: tg.tag })}
						className={`inline-flex items-center gap-1.5 rounded-g-pill px-2 py-1 text-g-xs ${
							excluded
								? "bg-g-danger/20 text-g-danger line-through"
								: "bg-g-success/20 text-g-success"
						}`}
					>
						{excluded && <FaBan size={9} className="no-underline" />}
						{tagData[tg.tag]?.display_name}
						<FaTimes size={9} />
					</button>
				);
			})}
			<Button variant="ghost" size="xs" onClick={clearAll}>
				{t("db.clear")}
			</Button>
		</div>
	);
}
