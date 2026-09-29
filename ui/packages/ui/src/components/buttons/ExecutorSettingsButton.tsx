import { Button } from "@gcsim/primitives";
import { Settings } from "lucide-react";
import { useTranslation } from "react-i18next";
import { usePrefs } from "../../stores/AppState";

// TODO: translation
export default () => {
	const { t } = useTranslation();
	const { setSettingsOpen } = usePrefs();

	return (
		<Button variant="secondary" onClick={() => setSettingsOpen(true)}>
			<Settings />
			{t("simple.settings")}
		</Button>
	);
};
