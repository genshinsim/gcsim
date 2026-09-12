import { Button } from "@blueprintjs/core";
import { useTranslation } from "react-i18next";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

// TODO: translation
export default ({}) => {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();

	return (
		<Button
			icon="cog"
			text={t<string>("simple.settings")}
			onClick={() => dispatch(appActions.setSettingsOpen(true))}
		/>
	);
};
