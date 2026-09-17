import { Button } from "@gcsim/primitives";
import { Settings } from "lucide-react";
import { useTranslation } from "react-i18next";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

// TODO: translation
export default () => {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();

	return (
		<Button
			variant="secondary"
			onClick={() => dispatch(appActions.setSettingsOpen(true))}
		>
			<Settings />
			{t("simple.settings")}
		</Button>
	);
};
