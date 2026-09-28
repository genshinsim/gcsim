import { Alert, AlertDescription, AlertTitle } from "@gcsim/primitives";
import { useTranslation } from "react-i18next";

export function ConfigError({ error }: { error: string | null }) {
	const { t } = useTranslation();
	if (!error) {
		return null;
	}
	return (
		<Alert variant="destructive">
			<AlertTitle>
				{t("viewer.error_encountered") + t("viewer.config_invalid")}
			</AlertTitle>
			<AlertDescription>
				<pre className="whitespace-pre-wrap">{error}</pre>
			</AlertDescription>
		</Alert>
	);
}
