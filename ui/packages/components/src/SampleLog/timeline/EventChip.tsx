import { cn } from "@gcsim/primitives";
import type React from "react";
import { useTranslation } from "react-i18next";
import { display } from "./display";
import { type Chip, chipShortText, chipText, hasDuration } from "./model";

export function EventChip({
	chip,
	onOpen,
	onDurationIcon,
	matched,
	inDuration,
	className,
}: {
	chip: Chip;
	onOpen: (c: Chip) => void;
	onDurationIcon: (c: Chip) => void;
	matched?: boolean;
	inDuration?: boolean;
	className?: string;
}) {
	const { t } = useTranslation();
	const cls = cn(
		"flex h-[18px] items-center gap-1 overflow-hidden whitespace-nowrap rounded-[3px] px-1 text-left text-[11px] leading-[18px] text-white",
		matched && "ring-2 ring-g-warning ring-offset-1 ring-offset-g-surface-2",
		inDuration && "ring-2 ring-g-ink ring-offset-1 ring-offset-g-surface-2",
		className,
	);
	const { color, icon } = display(chip.event.type);
	const text = chipShortText(chip);
	const title = `${chip.frame} · ${chip.event.type}: ${chipText(chip)}`;
	const style: React.CSSProperties = { backgroundColor: color };
	if (hasDuration(chip)) {
		const durationLabel = t("sample.show_duration", {
			start: chip.event.frame,
			end: Number.isFinite(chip.event.end) ? chip.event.end : t("sample.end"),
		});
		return (
			<div title={title} className={cls} style={style}>
				<button
					type="button"
					className="-ml-0.5 flex shrink-0 rounded-[2px] px-0.5 hover:bg-white/25"
					title={durationLabel}
					aria-label={durationLabel}
					onClick={() => onDurationIcon(chip)}
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
