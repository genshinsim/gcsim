import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import { useCallback, useMemo, useState } from "react";
import { SampleEventDetails } from "../../SampleEventDetails";
import { type Chip, chipText } from "./model";

export function useEventDetails() {
	const [chip, setChip] = useState<Chip | null>(null);
	const open = useCallback((c: Chip) => setChip(c), []);
	const dialog = useMemo(
		() => (
			<Dialog
				open={chip != null}
				onOpenChange={(o) => {
					if (!o) {
						setChip(null);
					}
				}}
			>
				<DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-3xl">
					<DialogHeader>
						<DialogTitle>
							{chip != null && `${chip.frame} · ${chipText(chip)}`}
						</DialogTitle>
					</DialogHeader>
					{chip != null && <SampleEventDetails data={chip.event.raw} />}
				</DialogContent>
			</Dialog>
		),
		[chip],
	);
	return { dialog, open };
}
