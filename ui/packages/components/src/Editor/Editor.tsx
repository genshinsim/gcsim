import { Editor as ConfigEditor, DEFAULT_EDITOR_THEME } from "@gcsim/editor";
import { useTranslation } from "react-i18next";
import { EditorSettings } from "./EditorSettings";
import type { EditorProps } from "./types";

export const Editor = ({
	value,
	onChange,
	fontSize = 14,
	theme = DEFAULT_EDITOR_THEME,
	onAppearanceChange,
	maxLines,
	error,
}: EditorProps) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-col gap-1">
			{onAppearanceChange ? (
				<div className="flex justify-end gap-1">
					<EditorSettings
						appearance={{ fontSize, theme }}
						onChange={onAppearanceChange}
					/>
				</div>
			) : null}
			<ConfigEditor
				id="config_editor"
				value={value}
				onChange={onChange}
				fontSize={fontSize}
				theme={theme}
				maxLines={maxLines}
				error={error}
				errorTitle={t("viewer.config_invalid")}
			/>
		</div>
	);
};
