import { cn } from "@gcsim/primitives";
import type React from "react";
import { useTranslation } from "react-i18next";
import { display } from "./display";
import { type Chip, chipText } from "./model";

export const hasDuration = (c: Chip) => c.event.end !== undefined;

/** where the event's duration ends; one that never ends runs past the sample */
export const endOf = (c: Chip) => c.event.end ?? Number.POSITIVE_INFINITY;

export function EventChip({
	chip,
	onOpen,
	onDuration,
	matched,
	lit,
	className,
}: {
	chip: Chip;
	onOpen: (c: Chip) => void;
	/** clicking the icon of a chip with a duration calls this instead */
	onDuration?: (c: Chip) => void;
	matched?: boolean;
	/** part of the highlighted duration */
	lit?: boolean;
	className?: string;
}) {
	const { t } = useTranslation();
	const cls = cn(
		"flex h-[18px] items-center gap-1 overflow-hidden whitespace-nowrap rounded-[3px] px-1 text-left text-[11px] leading-[18px] text-white",
		matched && "ring-2 ring-g-warning ring-offset-1 ring-offset-g-surface-2",
		lit && "ring-2 ring-g-ink ring-offset-1 ring-offset-g-surface-2",
		className,
	);
	const { color, icon } = display(chip.event.type);
	const text = chipText(chip);
	const title = `${chip.frame} · ${chip.event.type}: ${text}`;
	const style: React.CSSProperties = { backgroundColor: color };
	if (onDuration != null && hasDuration(chip)) {
		const durationLabel = t("sample.show_duration", {
			start: chip.event.frame,
			end: chip.event.end ?? t("sample.end"),
		});
		return (
			<div title={title} className={cls} style={style}>
				<button
					type="button"
					className="-ml-0.5 flex shrink-0 rounded-[2px] px-0.5 hover:bg-white/25"
					title={durationLabel}
					aria-label={durationLabel}
					onClick={() => onDuration(chip)}
				>
					<span className="material-icons !text-[12px]" aria-hidden>
						{icon}
					</span>
				</button>
				<button
					type="button"
					className="min-w-0 flex-1 truncate text-left hover:brightness-125"
					onClick={() => onOpen(chip)}
				>
					{text}
				</button>
			</div>
		);
	}
	return (
		<button
			type="button"
			title={title}
			onClick={() => onOpen(chip)}
			className={cn(cls, "hover:brightness-125")}
			style={style}
		>
			<span className="material-icons shrink-0 !text-[12px]" aria-hidden>
				{icon}
			</span>
			<span className="truncate">{text}</span>
		</button>
	);
}
