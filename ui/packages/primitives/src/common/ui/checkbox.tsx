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
				"peer size-4 shrink-0 rounded-[4px] border border-deprecated-input shadow-xs transition-shadow outline-none focus-visible:border-deprecated-ring focus-visible:ring-[3px] focus-visible:ring-deprecated-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-deprecated-destructive aria-invalid:ring-deprecated-destructive/20 data-[state=checked]:border-deprecated-primary data-[state=checked]:bg-deprecated-primary data-[state=checked]:text-deprecated-primary-foreground dark:bg-deprecated-input/30 dark:aria-invalid:ring-deprecated-destructive/40 dark:data-[state=checked]:bg-deprecated-primary",
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
