import { Separator } from "@gcsim/primitives";
import type * as React from "react";
import { cn } from "../../lib/utils";

function Navbar({ className, ...props }: React.ComponentProps<"nav">) {
	return (
		<nav
			data-slot="navbar"
			className={cn(
				"flex h-12 w-full items-center gap-1 border-b border-border bg-background px-4 text-foreground",
				className,
			)}
			{...props}
		/>
	);
}

function NavbarGroup({
	className,
	align = "start",
	...props
}: React.ComponentProps<"div"> & { align?: "start" | "end" }) {
	return (
		<div
			data-slot="navbar-group"
			data-align={align}
			className={cn(
				"flex items-center gap-1",
				align === "end" && "ml-auto",
				className,
			)}
			{...props}
		/>
	);
}

function NavbarHeading({ className, ...props }: React.ComponentProps<"div">) {
	return (
		<div
			data-slot="navbar-heading"
			className={cn("mr-2 flex items-center font-medium", className)}
			{...props}
		/>
	);
}

function NavbarDivider({
	className,
	...props
}: React.ComponentProps<typeof Separator>) {
	return (
		<Separator
			data-slot="navbar-divider"
			orientation="vertical"
			className={cn("mx-2 h-6 self-center", className)}
			{...props}
		/>
	);
}

export { Navbar, NavbarDivider, NavbarGroup, NavbarHeading };
