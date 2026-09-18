import { cva, type VariantProps } from "class-variance-authority";
import { Slot as SlotPrimitive } from "radix-ui";
import type * as React from "react";

import { cn } from "../../lib/utils";

const badgeVariants = cva(
	"inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-full border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:border-deprecated-ring focus-visible:ring-[3px] focus-visible:ring-deprecated-ring/50 aria-invalid:border-deprecated-destructive aria-invalid:ring-deprecated-destructive/20 dark:aria-invalid:ring-deprecated-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3",
	{
		variants: {
			variant: {
				default:
					"bg-deprecated-primary text-deprecated-primary-foreground [a&]:hover:bg-deprecated-primary/90",
				secondary:
					"bg-deprecated-secondary text-deprecated-secondary-foreground [a&]:hover:bg-deprecated-secondary/90",
				destructive:
					"bg-deprecated-destructive text-white focus-visible:ring-deprecated-destructive/20 dark:bg-deprecated-destructive/60 dark:focus-visible:ring-deprecated-destructive/40 [a&]:hover:bg-deprecated-destructive/90",
				success:
					"bg-deprecated-success text-deprecated-success-foreground [a&]:hover:bg-deprecated-success/90",
				warning:
					"bg-deprecated-warning text-deprecated-warning-foreground [a&]:hover:bg-deprecated-warning/90",
				outline:
					"border-deprecated-border text-deprecated-foreground [a&]:hover:bg-deprecated-accent [a&]:hover:text-deprecated-accent-foreground",
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
