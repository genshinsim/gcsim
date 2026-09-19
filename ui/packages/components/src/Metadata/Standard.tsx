// TODO(gauge-migration): not ported — still on --deprecated-* (story-only — not reachable from the web/db runtime import graph; only referenced by Standard.stories and a commented-out PreviewCard usage)
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
			title={t("result.metadata_standard")}
			value={standard}
			intent="success"
			bold
		/>
	);
});
