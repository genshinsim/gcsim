import { cva, type VariantProps } from "class-variance-authority";
import { Slot as SlotPrimitive } from "radix-ui";
import type * as React from "react";

import { cn } from "../../lib/utils";

const buttonVariants = cva(
	"inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-deprecated-ring focus-visible:ring-[3px] focus-visible:ring-deprecated-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-deprecated-destructive aria-invalid:ring-deprecated-destructive/20 dark:aria-invalid:ring-deprecated-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	{
		variants: {
			variant: {
				default:
					"bg-deprecated-primary text-deprecated-primary-foreground hover:bg-deprecated-primary/90",
				destructive:
					"bg-deprecated-destructive text-white hover:bg-deprecated-destructive/90 focus-visible:ring-deprecated-destructive/20 dark:bg-deprecated-destructive/60 dark:focus-visible:ring-deprecated-destructive/40",
				outline:
					"border bg-deprecated-background shadow-xs hover:bg-deprecated-accent hover:text-deprecated-accent-foreground dark:border-deprecated-input dark:bg-deprecated-input/30 dark:hover:bg-deprecated-input/50",
				secondary:
					"bg-deprecated-secondary text-deprecated-secondary-foreground hover:bg-deprecated-secondary/80",
				ghost:
					"hover:bg-deprecated-accent hover:text-deprecated-accent-foreground dark:hover:bg-deprecated-accent/50",
				link: "text-deprecated-primary underline-offset-4 hover:underline",
			},
			size: {
				default: "h-9 px-4 py-2 has-[>svg]:px-3",
				xs: "h-6 gap-1 rounded-md px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3",
				sm: "h-8 gap-1.5 rounded-md px-3 has-[>svg]:px-2.5",
				lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
				icon: "size-9",
				"icon-xs": "size-6 rounded-md [&_svg:not([class*='size-'])]:size-3",
				"icon-sm": "size-8",
				"icon-lg": "size-10",
			},
		},
		defaultVariants: {
			variant: "default",
			size: "default",
		},
	},
);

function Button({
	className,
	variant = "default",
	size = "default",
	asChild = false,
	...props
}: React.ComponentProps<"button"> &
	VariantProps<typeof buttonVariants> & {
		asChild?: boolean;
	}) {
	const Comp = asChild ? SlotPrimitive.Slot : "button";

	return (
		<Comp
			data-slot="button"
			data-variant={variant}
			data-size={size}
			className={cn(buttonVariants({ variant, size, className }))}
			{...props}
		/>
	);
}

export { Button, buttonVariants };
