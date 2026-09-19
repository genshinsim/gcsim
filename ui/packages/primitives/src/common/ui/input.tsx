import type * as React from "react";

import { cn } from "../../lib/utils";

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
	return (
		<input
			type={type}
			data-slot="input"
			className={cn(
				"h-g-ctrl w-full min-w-0 rounded-g-input border border-g-line bg-g-surface-2 px-3 py-1 font-g-body text-g-body transition-[color,box-shadow] outline-none selection:bg-g-accent selection:text-g-accent-fg file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-g-sm file:font-medium file:text-g-ink placeholder:text-g-ink-mute disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50",
				"focus-visible:border-g-accent focus-visible:ring-[3px] focus-visible:ring-g-accent/50",
				"aria-invalid:border-g-danger aria-invalid:ring-g-danger/20",
				className,
			)}
			{...props}
		/>
	);
}

export { Input };
