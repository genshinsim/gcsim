import { Badge, type BadgeProps } from "@gcsim/primitives";
import classNames from "classnames";
import { cn } from "../lib/utils";

type ItemProps = {
	title?: string;
	value?: string;
	intent?: BadgeProps["variant"];
	bright?: boolean;
	valueCase?: string;
	bold?: boolean;
	className?: string;
};

export const Item = ({
	title,
	value,
	intent = "default",
	bright,
	bold,
	valueCase = "uppercase",
	className = "",
}: ItemProps) => {
	if (value == null) {
		return null;
	}

	const titleCls = classNames("leading-4 align-bottom text-g-xs lowercase", {
		"text-g-ink-mute": !bright,
	});

	const cc = cn("font-g-mono", className);

	return (
		<Badge variant={intent} className={cc}>
			<span className="flex flex-row items-center gap-2 select-none">
				{title != null && <span className={titleCls}>{title}</span>}
				<span
					className={`${
						bold ? "font-bold" : ""
					} leading-4 align-bottom text-g-sm ${valueCase}`}
				>
					{value}
				</span>
			</span>
		</Badge>
	);
};
