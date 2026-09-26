import type React from "react";

export interface SectionDividerProps {
	children: React.ReactNode;
	fontClass?: string;
}

export function SectionDivider({
	children,
	fontClass = "font-bold text-g-lg",
}: SectionDividerProps) {
	return (
		<div className="mt-2 mb-2 flex flex-row place-items-center">
			<div className="mt-1 mr-1 ml-1 h-0 flex-grow border-t"></div>
			<span className={fontClass}>{children}</span>
			<div className="mt-1 mr-1 ml-1 h-0 flex-grow border-t"></div>
		</div>
	);
}
