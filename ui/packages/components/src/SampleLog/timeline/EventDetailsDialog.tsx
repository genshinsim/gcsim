import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import { SampleEventDetails } from "../../SampleEventDetails";
import { type Chip, chipText } from "./model";

export function EventDetailsDialog({
	chip,
	onClose,
}: {
	chip: Chip | null;
	onClose: () => void;
}) {
	return (
		<Dialog
			open={chip != null}
			onOpenChange={(o) => {
				if (!o) {
					onClose();
				}
			}}
		>
			<DialogContent
				className="max-h-[90vh] overflow-y-auto sm:max-w-3xl"
				aria-describedby={undefined}
			>
				<DialogHeader>
					<DialogTitle>
						{chip != null && `${chip.frame} · ${chipText(chip)}`}
					</DialogTitle>
				</DialogHeader>
				{chip != null && <SampleEventDetails data={chip.event.raw} />}
			</DialogContent>
		</Dialog>
	);
}
