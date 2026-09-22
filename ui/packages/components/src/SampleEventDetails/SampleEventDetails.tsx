import {
	collapseAllNested,
	defaultStyles,
	JsonView,
} from "react-json-view-lite";
import "react-json-view-lite/dist/index.css";
import { cn } from "../lib/utils";

type StyleProps = typeof defaultStyles;

// Value colours: the library class is nothing but a hardcoded colour, so a
// bare token replaces it outright.
const replacedValueColours = {
	label: "mr-[5px] font-semibold text-g-accent",
	nullValue: "text-g-ink-mute",
	undefinedValue: "text-g-ink-mute",
	numberValue: "text-g-hydro",
	stringValue: "text-g-success",
	booleanValue: "text-g-warning",
	otherValue: "text-g-ink",
	punctuation: "text-g-ink-dim",
};

// Structural classes carry behaviour we keep (expand/collapse glyphs, clickable
// cursor); we merge a token tint on top rather than replace them.
const tintedStructuralClasses = {
	// The library's container hardcodes a light `#eee` background; drop it so the
	// dialog's themed `bg-g-surface` shows through in both light and dark mode.
	container: cn(defaultStyles.container, "bg-transparent!"),
	clickableLabel: cn(defaultStyles.clickableLabel, "text-g-accent"),
	collapseIcon: cn(defaultStyles.collapseIcon, "text-g-ink-mute"),
	expandIcon: cn(defaultStyles.expandIcon, "text-g-ink-mute"),
	collapsedContent: cn(defaultStyles.collapsedContent, "text-g-ink-mute"),
};

const themedStyles: StyleProps = {
	...defaultStyles,
	...replacedValueColours,
	...tintedStructuralClasses,
};

export interface SampleEventDetailsProps {
	/** Structured event object, rendered as a collapsible tree. Preferred. */
	data?: unknown;
	/** Pretty-printed JSON string, shown when `data` is absent. */
	raw?: string;
	/** Optional heading shown above the viewer (e.g. the event message). */
	title?: string;
	className?: string;
}

export function SampleEventDetails({
	data,
	raw,
	title,
	className,
}: SampleEventDetailsProps) {
	const hasStructured =
		data !== null && data !== undefined && typeof data === "object";
	return (
		<div className={cn("w-full font-g-mono text-g-ink text-g-xs", className)}>
			{title ? (
				<div className="mb-2 font-g-body font-semibold text-g-ink text-g-sm">
					{title}
				</div>
			) : null}
			{hasStructured ? (
				<JsonView
					data={data as object}
					style={themedStyles}
					shouldExpandNode={collapseAllNested}
				/>
			) : raw ? (
				<pre className="whitespace-pre-wrap break-words">{raw}</pre>
			) : (
				<div className="text-g-ink-mute">No event data.</div>
			)}
		</div>
	);
}
