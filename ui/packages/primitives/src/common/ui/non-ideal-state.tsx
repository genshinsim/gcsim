import { cva, type VariantProps } from "class-variance-authority";
import { LoaderCircle } from "lucide-react";
import type * as React from "react";

import { cn } from "../../lib/utils";

const nonIdealStateVariants = cva(
	"flex items-center justify-center gap-4 text-center text-deprecated-muted-foreground",
	{
		variants: {
			layout: {
				vertical: "flex-col",
				horizontal: "flex-row text-left",
			},
		},
		defaultVariants: {
			layout: "vertical",
		},
	},
);

type NonIdealStateProps = Omit<React.ComponentProps<"div">, "title"> &
	VariantProps<typeof nonIdealStateVariants> & {
		icon?: React.ReactNode;
		title?: React.ReactNode;
		description?: React.ReactNode;
		action?: React.ReactNode;
		loading?: boolean;
	};

function NonIdealState({
	className,
	layout = "vertical",
	icon,
	title,
	description,
	action,
	loading = false,
	children,
	...props
}: NonIdealStateProps) {
	const media = loading ? (
		<LoaderCircle className="size-8 animate-spin" />
	) : (
		icon
	);

	return (
		<div
			data-slot="non-ideal-state"
			className={cn(nonIdealStateVariants({ layout }), className)}
			{...props}
		>
			{media != null && (
				<div
					data-slot="non-ideal-state-icon"
					className="text-deprecated-muted-foreground/70 [&_svg]:size-8"
				>
					{media}
				</div>
			)}
			{(title != null || description != null) && (
				<div data-slot="non-ideal-state-text" className="flex flex-col gap-1">
					{title != null && (
						<div
							data-slot="non-ideal-state-title"
							className="font-medium text-deprecated-foreground"
						>
							{title}
						</div>
					)}
					{description != null && (
						<div
							data-slot="non-ideal-state-description"
							className="text-sm text-deprecated-muted-foreground"
						>
							{description}
						</div>
					)}
				</div>
			)}
			{children}
			{action != null && <div data-slot="non-ideal-state-action">{action}</div>}
		</div>
	);
}

export type { NonIdealStateProps };
export { NonIdealState, nonIdealStateVariants };
