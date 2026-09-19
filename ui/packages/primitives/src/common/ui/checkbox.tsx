import { CheckIcon } from "lucide-react";
import { Checkbox as CheckboxPrimitive } from "radix-ui";
import type * as React from "react";

import { cn } from "../../lib/utils";

function Checkbox({
	className,
	...props
}: React.ComponentProps<typeof CheckboxPrimitive.Root>) {
	return (
		<CheckboxPrimitive.Root
			data-slot="checkbox"
			className={cn(
				"peer size-4 shrink-0 rounded-g-sm border border-g-line bg-g-surface-2 transition-shadow outline-none focus-visible:border-g-accent focus-visible:ring-[3px] focus-visible:ring-g-accent/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-g-danger aria-invalid:ring-g-danger/20 data-[state=checked]:border-g-accent data-[state=checked]:bg-g-accent data-[state=checked]:text-g-accent-fg",
				className,
			)}
			{...props}
		>
			<CheckboxPrimitive.Indicator
				data-slot="checkbox-indicator"
				className="grid place-content-center text-current transition-none"
			>
				<CheckIcon className="size-3.5" />
			</CheckboxPrimitive.Indicator>
		</CheckboxPrimitive.Root>
	);
}

export { Checkbox };
