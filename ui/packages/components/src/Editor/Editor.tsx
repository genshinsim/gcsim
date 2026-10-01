import { Editor as ConfigEditor, DEFAULT_EDITOR_THEME } from "@gcsim/editor";
import { useTranslation } from "react-i18next";
import { EditorSettings } from "./EditorSettings";
import {
	ImportedCharacterInsert,
	insertCharacterBlock,
} from "./ImportedCharacterInsert";
import { ShortcutSearch } from "./ShortcutSearch";
import type { EditorProps } from "./types";

export const Editor = ({
	value,
	onChange,
	fontSize = 14,
	theme = DEFAULT_EDITOR_THEME,
	onAppearanceChange,
	maxLines,
	error,
	importedCharacters,
}: EditorProps) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-col gap-1">
			<div className="flex justify-end gap-1">
				{importedCharacters ? (
					<ImportedCharacterInsert
						characters={importedCharacters}
						onInsert={(option) => onChange(insertCharacterBlock(value, option))}
					/>
				) : null}
				<ShortcutSearch />
				{onAppearanceChange ? (
					<EditorSettings
						appearance={{ fontSize, theme }}
						onChange={onAppearanceChange}
					/>
				) : null}
			</div>
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
