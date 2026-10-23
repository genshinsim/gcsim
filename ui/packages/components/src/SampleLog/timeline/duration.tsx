import { X } from "lucide-react";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import type { SimEvent } from "../events/types";
import { secondsLabel } from "./layout";
import { spanLabel } from "./model";

export type Span = { start: number; end: number; label: string };

export const spanOf = (e: SimEvent | null): Span | null =>
	e == null
		? null
		: { start: e.frame, end: e.end ?? e.frame, label: spanLabel(e) };

export function useClearOnEscape(active: boolean, clear: () => void) {
	useEffect(() => {
		if (!active) {
			return;
		}
		const onKey = (e: KeyboardEvent) => {
			if (e.key === "Escape" && !e.defaultPrevented) {
				clear();
			}
		};
		window.addEventListener("keydown", onKey);
		return () => window.removeEventListener("keydown", onKey);
	}, [active, clear]);
}

export function DurationClear({
	span,
	onClear,
}: {
	span: Span;
	onClear: () => void;
}) {
	const { t } = useTranslation();
	const title = t("sample.clear_duration", {
		label: span.label,
		start: span.start,
		end: Number.isFinite(span.end) ? span.end : t("sample.end"),
	});
	return (
		<button
			type="button"
			className="mr-1 flex min-w-0 flex-1 items-center gap-1 rounded-g-sm text-left text-g-warning hover:bg-g-surface-3"
			title={title}
			aria-label={title}
			onClick={onClear}
		>
			<span className="flex min-w-0 flex-1 flex-col">
				<span className="truncate">{span.label}</span>
				<span>
					{Number.isFinite(span.end)
						? secondsLabel(span.end - span.start)
						: t("sample.to_end")}
				</span>
			</span>
			<X className="size-3 shrink-0" />
		</button>
	);
}
