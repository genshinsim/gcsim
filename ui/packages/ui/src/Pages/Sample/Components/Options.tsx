import {
	Button,
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@gcsim/primitives";
import { Trans, useTranslation } from "react-i18next";
import { eventColor } from "./parse";

export interface OptionsProp {
	isOpen: boolean;
	handleClose: () => void;
	handleToggle: (opt: string) => void;
	handleClear: () => void;
	handleResetDefault: () => void;
	handleSetPresets: (opt: "simple" | "advanced" | "verbose" | "debug") => void;
	selected: string[];
	options: string[];
}

export function Options(props: OptionsProp) {
	const { t } = useTranslation();

	const cols = props.options.map((o) => {
		return (
			<div className="flex flex-row gap-1 p-1 items-center" key={o}>
				<label className="cursor-pointer">
					<input
						type="checkbox"
						checked={props.selected.indexOf(o) > -1}
						className="checkbox cursor-pointer"
						onChange={() => props.handleToggle(o)}
					/>
					<span
						className="font-medium text-sm pl-1"
						style={{ color: eventColor(o) }}
					>
						{o}
					</span>
				</label>
			</div>
		);
	});

	return (
		<Dialog
			open={props.isOpen}
			onOpenChange={(open) => !open && props.handleClose()}
		>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>
						<Trans>viewer.log_options</Trans>
					</DialogTitle>
				</DialogHeader>
				<div className="grid grid-cols-2 sm:grid-cols-3">{cols}</div>
				<DialogFooter className="!flex !flex-col !gap-1.5 sm:!flex-row sm:!gap-0">
					<Button
						variant="secondary"
						onClick={() => props.handleSetPresets("simple")}
					>
						{t("viewer.simple")}
					</Button>
					<Button
						variant="secondary"
						onClick={() => props.handleSetPresets("advanced")}
					>
						{t("viewer.advanced")}
					</Button>
					<Button
						variant="secondary"
						onClick={() => props.handleSetPresets("verbose")}
					>
						{t("viewer.verbose")}
					</Button>
					<Button
						variant="secondary"
						onClick={() => props.handleSetPresets("debug")}
					>
						{t("viewer.debug")}
					</Button>
					<Button variant="destructive" onClick={props.handleClear}>
						{t("viewer.clear")}
					</Button>
					<Button variant="outline" onClick={props.handleClose}>
						{t("viewer.close")}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
