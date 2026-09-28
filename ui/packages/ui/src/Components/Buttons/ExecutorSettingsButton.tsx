import { Button } from "@gcsim/primitives";
import { Settings } from "lucide-react";
import { useTranslation } from "react-i18next";
import { prefs } from "../../Stores";

// TODO: translation
export default () => {
	const { t } = useTranslation();

	return (
		<Button
			variant="secondary"
			onClick={() => prefs.set((p) => ({ ...p, settingsOpen: true }))}
		>
			<Settings />
			{t("simple.settings")}
		</Button>
	);
};
