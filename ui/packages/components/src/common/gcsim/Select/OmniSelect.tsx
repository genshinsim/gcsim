import {
	Command,
	CommandEmpty,
	CommandInput,
	CommandList,
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import type { ReactNode } from "react";
import { Fragment, useEffect, useState } from "react";
import { cn } from "../../../lib/utils";

export type OmniSelectItemState = {
	selected: boolean;
	onSelect: () => void;
};

export type OmniSelectProps<T> = {
	isOpen: boolean;
	onClose: () => void;
	items: T[];
	itemKey: (item: T) => string;
	itemPredicate: (item: T, query: string) => boolean;
	itemRenderer: (item: T, state: OmniSelectItemState) => ReactNode;
	onSelect: (item: T) => void;
	value?: T;
	title?: string;
	placeholder?: string;
	emptyMessage?: string;
	className?: string;
};

export function OmniSelect<T>({
	isOpen,
	onClose,
	items,
	itemKey,
	itemPredicate,
	itemRenderer,
	onSelect,
	value,
	title = "Select an item",
	placeholder,
	emptyMessage = "No results found.",
	className,
}: OmniSelectProps<T>) {
	const [query, setQuery] = useState("");

	useEffect(() => {
		if (isOpen) {
			setQuery("");
		}
	}, [isOpen]);

	const filtered = items.filter((item) => itemPredicate(item, query));
	const selectedKey = value !== undefined ? itemKey(value) : undefined;

	return (
		<Dialog
			open={isOpen}
			onOpenChange={(open) => {
				if (!open) {
					onClose();
				}
			}}
		>
			<DialogContent
				className={cn("gap-0 overflow-hidden p-0", className)}
				showCloseButton={false}
			>
				<DialogHeader className="sr-only">
					<DialogTitle>{title}</DialogTitle>
					<DialogDescription>{placeholder ?? title}</DialogDescription>
				</DialogHeader>
				<Command shouldFilter={false}>
					<CommandInput
						value={query}
						onValueChange={setQuery}
						placeholder={placeholder}
					/>
					<CommandList>
						<CommandEmpty>{emptyMessage}</CommandEmpty>
						{filtered.map((item) => (
							<Fragment key={itemKey(item)}>
								{itemRenderer(item, {
									selected: itemKey(item) === selectedKey,
									onSelect: () => onSelect(item),
								})}
							</Fragment>
						))}
					</CommandList>
				</Command>
			</DialogContent>
		</Dialog>
	);
}
