import { cva, type VariantProps } from "class-variance-authority";
import { Slot as SlotPrimitive } from "radix-ui";
import type * as React from "react";

import { cn } from "../../lib/utils";

const badgeVariants = cva(
	"inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-g-pill border border-transparent px-2 py-0.5 text-g-xs font-semibold whitespace-nowrap transition-[color,box-shadow] focus-visible:border-g-accent focus-visible:ring-[3px] focus-visible:ring-g-accent/50 aria-invalid:border-g-danger aria-invalid:ring-g-danger/20 [&>svg]:pointer-events-none [&>svg]:size-3",
	{
		variants: {
			variant: {
				default: "bg-g-accent-weak text-g-accent",
				secondary: "bg-g-surface-2 text-g-ink-dim",
				destructive: "bg-g-danger/15 text-g-danger",
				success: "bg-g-success/15 text-g-success",
				warning: "bg-g-warning/15 text-g-warning",
				outline: "border-g-line text-g-ink-dim [a&]:hover:bg-g-surface-2",
			},
		},
		defaultVariants: {
			variant: "default",
		},
	},
);

export type BadgeProps = React.ComponentProps<"span"> &
	VariantProps<typeof badgeVariants> & { asChild?: boolean };

function Badge({
	className,
	variant = "default",
	asChild = false,
	...props
}: BadgeProps) {
	const Comp = asChild ? SlotPrimitive.Slot : "span";

	return (
		<Comp
			data-slot="badge"
			data-variant={variant}
			className={cn(badgeVariants({ variant }), className)}
			{...props}
		/>
	);
}

export { Badge, badgeVariants };
