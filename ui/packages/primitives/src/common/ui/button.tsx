import { cva, type VariantProps } from "class-variance-authority";
import { Slot as SlotPrimitive } from "radix-ui";
import type * as React from "react";

import { cn } from "../../lib/utils";

const buttonVariants = cva(
	"inline-flex shrink-0 items-center justify-center gap-g-base rounded-g-btn font-g-body text-g-body font-semibold whitespace-nowrap transition-all outline-none focus-visible:border-g-accent focus-visible:ring-[3px] focus-visible:ring-g-accent/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-g-danger aria-invalid:ring-g-danger/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
	{
		variants: {
			variant: {
				default: "bg-g-accent text-g-accent-fg hover:bg-g-accent-hover",
				destructive:
					"bg-g-danger text-white hover:bg-g-danger/90 focus-visible:ring-g-danger/20",
				outline:
					"border border-g-line bg-g-surface-2 text-g-ink hover:border-g-accent",
				secondary: "bg-g-surface-2 text-g-ink hover:bg-g-surface-3",
				ghost: "text-g-ink hover:bg-g-surface-2",
				link: "text-g-accent underline-offset-4 hover:underline",
			},
			size: {
				default: "h-g-ctrl px-g-btn-x py-2 has-[>svg]:px-3",
				xs: "h-6 gap-1 rounded-g-btn px-2 text-g-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3",
				sm: "h-8 gap-1.5 rounded-g-btn px-3 text-g-sm has-[>svg]:px-2.5",
				lg: "h-10 rounded-g-btn px-6 text-g-lg has-[>svg]:px-4",
				icon: "size-g-ctrl",
				"icon-xs": "size-6 rounded-g-btn [&_svg:not([class*='size-'])]:size-3",
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
