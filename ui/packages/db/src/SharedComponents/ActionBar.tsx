import { dynamicKey } from "@gcsim/localization";
import { Button, Input } from "@gcsim/primitives";
import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaArrowDown, FaArrowUp, FaSearch, FaTimes } from "react-icons/fa";
import useDebounce from "../SharedHooks/debounce";
import { Filter } from "./Filter";
import {
	FilterContext,
	FilterDispatchContext,
	ItemFilterState,
	SortByDirection,
	sortByParams,
} from "./FilterComponents/Filter.utils";

export function ActionBar({ simCount }: { simCount: number | null }) {
	const { t } = useTranslation();

	return (
		<div className="sticky top-0 z-10 -mx-8 flex flex-col gap-g-base border-b border-g-line-soft bg-g-canvas/90 px-8 py-g-base backdrop-blur">
			<div className="flex flex-wrap items-center gap-g-base">
				<Filter />
				<CustomFilterSearch />
				<div className="hidden sm:block">
					<SortControl />
				</div>
				<span className="ml-auto whitespace-nowrap font-g-mono text-g-sm text-g-ink-dim">
					{t("db.showing_simulations", { i: simCount ?? 0 })}
				</span>
			</div>
			<SelectedCharChips />
		</div>
	);
}

/** Debounced free-text query, dispatched as the server-side custom filter. */
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
				placeholder={t("db.customFilter")}
				value={value}
				onChange={(e) => setValue(e.target.value)}
			/>
		</div>
	);
}

function SelectedCharChips() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	const selected = Object.values(filter.charFilter).filter(
		(c) => c.state === ItemFilterState.include,
	);
	if (selected.length === 0) return null;

	return (
		<div className="flex flex-wrap items-center gap-g-base-sm">
			{selected.map((c) => (
				<button
					type="button"
					key={c.charName}
					onClick={() => dispatch({ type: "removeChar", char: c.charName })}
					className="inline-flex items-center gap-1.5 rounded-g-pill bg-g-accent-weak py-1 pl-1.5 pr-2 text-g-xs text-g-accent"
				>
					<img
						src={`/api/assets/avatar/${c.charName}.png`}
						alt=""
						className="size-4"
					/>
					{t("game:character_names." + c.charName)}
					<FaTimes size={9} />
				</button>
			))}
			<Button
				variant="ghost"
				size="xs"
				onClick={() => {
					for (const c of selected) {
						dispatch({ type: "removeChar", char: c.charName });
					}
				}}
			>
				{t("db.clear")}
			</Button>
		</div>
	);
}

/** Compact sort selector with a direction toggle, backed by the filter reducer. */
function SortControl() {
	const { t: translation } = useTranslation();
	const t = (s: string) => translation(dynamicKey(s)) as string;
	const filter = useContext(FilterContext);
	const dispatch = useContext(FilterDispatchContext);

	return (
		<div className="flex gap-g-base-sm">
			{sortByParams.map((param) => {
				const active = filter.sortBy.sortKey === param.sortKey;
				const dir = active ? filter.sortBy.sortByDirection : null;
				return (
					<Button
						key={param.sortKey}
						size="sm"
						variant={active ? "default" : "secondary"}
						onClick={() =>
							dispatch({ type: "handleSortBy", sortByKey: param.sortKey })
						}
					>
						{t(param.translationKey)}
						{dir === SortByDirection.asc && <FaArrowUp size={10} />}
						{dir === SortByDirection.dsc && <FaArrowDown size={10} />}
					</Button>
				);
			})}
		</div>
	);
}
