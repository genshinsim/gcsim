import type { JSX, ReactNode } from "react";
import { useTranslation } from "react-i18next";

export type TooltipFormat = (n?: number) => string | undefined;

type TitleProps = {
	title: string | JSX.Element;
	color?: string;
	percent?: number;
};

export const TooltipTitle = ({ title, color, percent }: TitleProps) => {
	const { i18n } = useTranslation();
	const value = percent?.toLocaleString(i18n.language, {
		minimumFractionDigits: 2,
		maximumFractionDigits: 2,
		style: "percent",
	});

	if (typeof title !== "string") {
		return title;
	}

	return (
		<div className="flex flex-row flex-nowrap justify-start text-g-ink-mute gap-2">
			<span className="whitespace-nowrap" style={{ color: color }}>
				{title}
			</span>
			{value != null && <span>{"(" + value + ")"}</span>}
		</div>
	);
};

export const TooltipList = ({ children }: { children: ReactNode }) => (
	<ul className="list-disc pl-4 grid grid-cols-[repeat(2,_max-content)] gap-x-2 justify-start">
		{children}
	</ul>
);

type RowProps = {
	name: string;
	value?: number | null;
	color?: string;
	format?: TooltipFormat;
	suffix?: string;
};

export const TooltipRow = ({ name, value, color, format, suffix }: RowProps) => {
	const { i18n } = useTranslation();
	const num =
		format == null
			? value?.toLocaleString(i18n.language, {
					minimumFractionDigits: 2,
					maximumFractionDigits: 2,
				})
			: format(value ?? 0);

	return (
		<>
			<span className="text-g-ink-mute list-item" style={{ color: color }}>
				{name}
			</span>
			<span>
				{num}
				{suffix}
			</span>
		</>
	);
};
