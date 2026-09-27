import { memo } from "react";

type Props = {
	title: string;
};

const CardTitle = ({ title }: Props) => (
	<div className="flex flex-row text-g-lg text-g-ink-mute items-center gap-2 outline-0">
		{title}
	</div>
);

export default memo(CardTitle);
