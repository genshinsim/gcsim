import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Item } from "./Item";

type Props = {
	standard?: string;
};

export const Standard = memo(({ standard }: Props) => {
	const { t } = useTranslation();
	if (standard == null) {
		return null;
	}

	return (
		<Item
			title={t<string>("result.metadata_standard")}
			value={standard}
			intent="success"
			bold
		/>
	);
});
