import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/db/$id")({
	component: Removed,
});

function Removed() {
	const { t } = useTranslation();
	return (
		<>
			<title>gcsim - db entry unavailable</title>
			<div className="m-2 text-center">{t("src.db_entry_unavailable")}</div>
		</>
	);
}
