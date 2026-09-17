import { MultiSelect } from "@gcsim/components";
import { dynamicKey } from "@gcsim/localization";
import { Badge, CommandItem } from "@gcsim/primitives";
import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { FaCheck, FaTimes } from "react-icons/fa";
import {
	charNames,
	FilterContext,
	FilterDispatchContext,
	ItemFilterState,
} from "./FilterComponents/Filter.utils";

export function CharacterQuickSelect() {
	const dispatch = useContext(FilterDispatchContext);
	const filter = useContext(FilterContext);
	const { t } = useTranslation();

	const includedChars = Object.entries(filter.charFilter)
		.filter(([, charState]) => charState.state === ItemFilterState.include)
		.map(([charName]) => charName);

	const translateCharName = (charName: string) =>
		t(dynamicKey("game:character_names." + charName));

	return (
		<div className="grow max-w-xl">
			<MultiSelect<string>
				items={charNames}
				itemKey={(charName) => charName}
				value={includedChars}
				placeholder={t("db.type_to_search")}
				itemPredicate={(charName, query) => {
					const normalizedQuery = query.toLocaleLowerCase();
					return (
						charName.toLowerCase().includes(normalizedQuery) ||
						translateCharName(charName).toLowerCase().includes(normalizedQuery)
					);
				}}
				itemRenderer={(charName, { selected, onSelect }) => (
					<CommandItem value={charName} onSelect={onSelect}>
						<img
							src={`/api/assets/avatar/${charName}.png`}
							alt=""
							className="w-6 h-6"
						/>
						<span className="flex-1">{translateCharName(charName)}</span>
						{selected && <FaCheck className="size-4" />}
					</CommandItem>
				)}
				tagRenderer={(charName, { onRemove }) => (
					<Badge variant="secondary" className="gap-1">
						<img
							src={`/api/assets/avatar/${charName}.png`}
							alt=""
							className="w-4 h-4"
						/>
						{translateCharName(charName)}
						<button
							type="button"
							aria-label={`remove ${translateCharName(charName)}`}
							onClick={onRemove}
						>
							<FaTimes className="size-3" />
						</button>
					</Badge>
				)}
				onChange={(next) => {
					const nextChars = new Set(next);
					const prevChars = new Set(includedChars);
					for (const charName of next) {
						if (!prevChars.has(charName)) {
							dispatch({ type: "includeChar", char: charName });
						}
					}
					for (const charName of includedChars) {
						if (!nextChars.has(charName)) {
							dispatch({ type: "removeChar", char: charName });
						}
					}
				}}
			/>
		</div>
	);
}
