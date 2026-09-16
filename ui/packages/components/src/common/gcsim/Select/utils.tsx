import { CommandItem } from "@gcsim/primitives";
import { Check } from "lucide-react";

export function createKeyOrLabelPredicate<T extends string>(
	label: (item: T) => string,
) {
	return (item: T, query: string) => {
		const normalized = query.trim().toLowerCase();
		if (normalized.length === 0) {
			return true;
		}
		return `${item} ${label(item)}`.toLowerCase().includes(normalized);
	};
}

export function renderLabeledItem<T extends string>(
	item: T,
	label: string,
	state: { selected: boolean; onSelect: () => void },
) {
	return (
		<CommandItem value={item} onSelect={state.onSelect}>
			<span className="flex-1">{label}</span>
			{state.selected && <Check className="size-4" />}
		</CommandItem>
	);
}
