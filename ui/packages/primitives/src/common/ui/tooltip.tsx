import { Tooltip as TooltipPrimitive } from "radix-ui";
import * as React from "react";

import { cn } from "../../lib/utils";

const InsideTooltipProvider = React.createContext(false);

function TooltipProvider({
	delayDuration = 0,
	...props
}: React.ComponentProps<typeof TooltipPrimitive.Provider>) {
	return (
		<InsideTooltipProvider.Provider value={true}>
			<TooltipPrimitive.Provider
				data-slot="tooltip-provider"
				delayDuration={delayDuration}
				{...props}
			/>
		</InsideTooltipProvider.Provider>
	);
}

function Tooltip({
	...props
}: React.ComponentProps<typeof TooltipPrimitive.Root>) {
	const root = <TooltipPrimitive.Root data-slot="tooltip" {...props} />;
	if (React.useContext(InsideTooltipProvider)) {
		return root;
	}
	return <TooltipProvider>{root}</TooltipProvider>;
}

function TooltipTrigger({
	...props
}: React.ComponentProps<typeof TooltipPrimitive.Trigger>) {
	return <TooltipPrimitive.Trigger data-slot="tooltip-trigger" {...props} />;
}

function TooltipContent({
	className,
	sideOffset = 4,
	...props
}: React.ComponentProps<typeof TooltipPrimitive.Content>) {
	return (
		<TooltipPrimitive.Portal>
			<TooltipPrimitive.Content
				data-slot="tooltip-content"
				sideOffset={sideOffset}
				className={cn(
					"z-50 w-fit origin-(--radix-tooltip-content-transform-origin) animate-in rounded-g-lg border border-g-line bg-g-surface px-2 py-1 text-g-xs text-g-ink shadow-g-pop fade-in-0 zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95",
					className,
				)}
				{...props}
			/>
		</TooltipPrimitive.Portal>
	);
}

export { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger };
