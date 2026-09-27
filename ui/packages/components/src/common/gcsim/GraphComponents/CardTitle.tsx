import { ReloadIcon } from "@radix-ui/react-icons";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type Props = {
	title: string;
	stale?: boolean;
};

const CardTitle = ({ title, stale = false }: Props) => (
	<div className="flex flex-row justify-between items-center gap-4">
		<div className="flex flex-row text-g-lg text-g-ink-mute items-center gap-2 outline-0">
			{title}
		</div>
		{stale ? <UpdatingStatus /> : null}
	</div>
);

const UpdatingStatus = () => {
	const { t } = useTranslation();
	const label = t("result.updating");

	return (
		<div
			role="status"
			aria-label={label}
			title={label}
			className="text-g-ink-mute text-g-xs cursor-default"
		>
			<ReloadIcon className="animate-spin" />
		</div>
	);
};

export default memo(CardTitle);
