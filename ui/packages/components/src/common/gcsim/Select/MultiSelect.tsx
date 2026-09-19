import {
	Button,
	Command,
	CommandEmpty,
	CommandInput,
	CommandList,
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@gcsim/primitives";
import type { ReactNode } from "react";
import { Fragment, useState } from "react";
import { cn } from "../../../lib/utils";

export type MultiSelectItemState = {
	selected: boolean;
	onSelect: () => void;
};

export type MultiSelectTagState = {
	onRemove: () => void;
};

export type MultiSelectProps<T> = {
	items: T[];
	itemKey: (item: T) => string;
	itemPredicate: (item: T, query: string) => boolean;
	itemRenderer: (item: T, state: MultiSelectItemState) => ReactNode;
	tagRenderer: (item: T, state: MultiSelectTagState) => ReactNode;
	value: T[];
	onChange: (value: T[]) => void;
	placeholder?: string;
	emptyMessage?: string;
	className?: string;
};

export function MultiSelect<T>({
	items,
	itemKey,
	itemPredicate,
	itemRenderer,
	tagRenderer,
	value,
	onChange,
	placeholder = "Select items...",
	emptyMessage = "No results found.",
	className,
}: MultiSelectProps<T>) {
	const [open, setOpen] = useState(false);
	const [query, setQuery] = useState("");

	const selectedKeys = new Set(value.map(itemKey));
	const filtered = items.filter((item) => itemPredicate(item, query));

	function toggle(item: T) {
		const key = itemKey(item);
		if (selectedKeys.has(key)) {
			onChange(value.filter((selected) => itemKey(selected) !== key));
		} else {
			onChange([...value, item]);
		}
	}

	function remove(item: T) {
		const key = itemKey(item);
		onChange(value.filter((selected) => itemKey(selected) !== key));
	}

	return (
		<Popover
			open={open}
			onOpenChange={(next) => {
				setOpen(next);
				if (next) {
					setQuery("");
				}
			}}
		>
			<div
				className={cn(
					"flex flex-wrap items-center gap-1 rounded-g-md border border-g-line px-2 py-1",
					className,
				)}
			>
				{value.map((item) => (
					<Fragment key={itemKey(item)}>
						{tagRenderer(item, { onRemove: () => remove(item) })}
					</Fragment>
				))}
				<PopoverTrigger asChild>
					<Button type="button" variant="ghost" size="sm">
						{placeholder}
					</Button>
				</PopoverTrigger>
			</div>
			<PopoverContent className="w-72 p-0" align="start">
				<Command shouldFilter={false}>
					<CommandInput
						value={query}
						onValueChange={setQuery}
						placeholder={placeholder}
					/>
					<CommandList>
						<CommandEmpty>{emptyMessage}</CommandEmpty>
						{filtered.map((item) => {
							const key = itemKey(item);
							return (
								<Fragment key={key}>
									{itemRenderer(item, {
										selected: selectedKeys.has(key),
										onSelect: () => toggle(item),
									})}
								</Fragment>
							);
						})}
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	);
}
