import { Editor as ConfigEditor, type EditorHandle } from "@gcsim/editor";
import { Button } from "@gcsim/primitives";
import { TextAlignStart } from "lucide-react";
import { useRef } from "react";
import { useTranslation } from "react-i18next";
import { EditorSettings } from "./EditorSettings";
import { useEditorPrefs } from "./editorPrefs";
import {
	ImportedCharacterInsert,
	insertCharacterBlock,
} from "./ImportedCharacterInsert";
import { ShortcutSearch } from "./ShortcutSearch";
import type { EditorProps } from "./types";

export const Editor = ({
	value,
	onChange,
	maxLines,
	error,
	importedCharacters,
}: EditorProps) => {
	const { t } = useTranslation();
	const editorRef = useRef<EditorHandle>(null);
	const [prefs, setPrefs] = useEditorPrefs();

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
				<EditorSettings appearance={prefs} onChange={setPrefs} />
			</div>
			<ConfigEditor
				ref={editorRef}
				id="config_editor"
				value={value}
				onChange={onChange}
				fontSize={prefs.fontSize}
				theme={prefs.theme}
				maxLines={maxLines}
				error={error}
				errorTitle={t("viewer.config_invalid")}
			/>
		</div>
	);
};
