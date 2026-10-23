import { Editor as ConfigEditor, type EditorHandle } from "@gcsim/editor";
import { Button, cn } from "@gcsim/primitives";
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

const FILL_MIN_LINES = 10;

export const Editor = ({
	value,
	onChange,
	maxLines,
	fill = false,
	scrollMarginBottom,
	error,
	importedCharacters,
}: EditorProps) => {
	const { t } = useTranslation();
	const editorRef = useRef<EditorHandle>(null);
	const toolbarRef = useRef<HTMLDivElement>(null);
	const [prefs, setPrefs] = useEditorPrefs();
	const scrollMargins = () => ({
		top: toolbarRef.current?.offsetHeight ?? 0,
		bottom: scrollMarginBottom?.() ?? 0,
	});

	return (
		<div className={cn("flex flex-col gap-1", fill && "flex-1")}>
			<div
				ref={toolbarRef}
				id="config_toolbar"
				className={cn(
					"flex justify-end gap-1",
					fill && "sticky top-0 z-10 bg-g-canvas py-1",
				)}
			>
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
				<EditorSettings prefs={prefs} onChange={setPrefs} />
			</div>
			<ConfigEditor
				ref={editorRef}
				id="config_editor"
				value={value}
				onChange={onChange}
				fontSize={prefs.fontSize}
				theme={prefs.theme}
				maxLines={fill ? Infinity : maxLines}
				minLines={fill ? FILL_MIN_LINES : undefined}
				className={fill ? "flex flex-1 flex-col px-px" : undefined}
				error={error}
				errorTitle={t("viewer.config_invalid")}
				errorPanelAt={fill ? "top" : "bottom"}
				scrollMargins={fill ? scrollMargins : undefined}
			/>
		</div>
	);
};
