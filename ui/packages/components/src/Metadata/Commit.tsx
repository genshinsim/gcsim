import { Badge } from "@gcsim/primitives";
import { memo } from "react";
import { useTranslation } from "react-i18next";
import { cn } from "../lib/utils";

type Props = {
	commit?: string;
	className?: string;
};

export const Commit = memo(({ commit, className }: Props) => {
	const { t } = useTranslation();
	if (commit == null || commit === "") {
		return null;
	}

	const shortCommit = commit?.substring(0, 7);
	const url = "https://github.com/genshinsim/gcsim/commits/" + commit;

	const cc = cn("text-sm font-mono", className);

	return (
		<Badge className={cc}>
			<span className="flex flex-row items-center gap-2">
				<span className="text-gray-400">{t("result.metadata_commit")}</span>
				<a href={url} target="_blank" rel="noreferrer" className="font-bold">
					{shortCommit}
				</a>
			</span>
		</Badge>
	);
});
