import {
	Editor as ConfigEditor,
	DEFAULT_EDITOR_THEME,
	type EditorHandle,
} from "@gcsim/editor";
import { Button } from "@gcsim/primitives";
import { TextAlignStart } from "lucide-react";
import { useRef } from "react";
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
	const editorRef = useRef<EditorHandle>(null);

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
				<Button
					variant="ghost"
					size="sm"
					aria-label={t("simple.format_config")}
					onClick={() => editorRef.current?.format()}
				>
					<TextAlignStart />
					<span className="hidden sm:inline">{t("simple.format_config")}</span>
				</Button>
				{onAppearanceChange ? (
					<EditorSettings
						appearance={{ fontSize, theme }}
						onChange={onAppearanceChange}
					/>
				) : null}
			</div>
			<ConfigEditor
				ref={editorRef}
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
