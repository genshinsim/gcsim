import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import React from "react";
import type { SampleItem } from "./parse";

export function SampleItemView({
	item,
	showBuffDuration,
}: {
	item: SampleItem;
	showBuffDuration: (e: SampleItem) => void;
}) {
	const [open, setOpen] = React.useState<boolean>(false);
	const handleClick = () => {
		setOpen(true);
	};
	return (
		<div
			className="flex flex-row gap-2 items-center pl-1 pr-1 pt-px pb-px rounded-g-md m-1 "
			style={{ backgroundColor: item.color }}
		>
			<button
				type="button"
				className="material-icons text-g-sm cursor-pointer"
				onClick={() => showBuffDuration(item)}
			>
				{item.icon}
			</button>
			<button
				type="button"
				className="flex-grow cursor-pointer text-left"
				onClick={handleClick}
			>
				{item.msg}
			</button>
			<div>{item.target}</div>
			<Dialog open={open} onOpenChange={setOpen}>
				<DialogContent className="max-h-[90vh] overflow-y-auto">
					<DialogHeader>
						<DialogTitle>{item.msg}</DialogTitle>
					</DialogHeader>
					<pre className="m-2 whitespace-pre-wrap">{item.raw}</pre>
				</DialogContent>
			</Dialog>
		</div>
	);
}
