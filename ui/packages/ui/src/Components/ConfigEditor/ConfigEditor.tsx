import AceEditor from "react-ace";
//THESE IMPORTS NEEDS TO BE AFTER IMPORTING AceEditor
import "ace-builds/src-noconflict/ext-language_tools";
import "../../util/mode-gcsim.dm.js";
//manually import supported themes cause we can't get for loop to work here
import {
	Label,
	NumberInput,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/theme-kuroir";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/theme-solarized_dark";
import "ace-builds/src-noconflict/theme-solarized_light";
import "ace-builds/src-noconflict/theme-terminal";
import "ace-builds/src-noconflict/theme-textmate";
import "ace-builds/src-noconflict/theme-tomorrow";
import "ace-builds/src-noconflict/theme-tomorrow_night";
import "ace-builds/src-noconflict/theme-twilight";
import "ace-builds/src-noconflict/theme-xcode";
import React from "react";
import { useTranslation } from "react-i18next";

export const themes = [
	"monokai",
	"github",
	"tomorrow",
	"tomorrow_night",
	"kuroir",
	"twilight",
	"xcode",
	"textmate",
	"solarized_dark",
	"solarized_light",
	"terminal",
];
export type Theme = (typeof themes)[number];

type Props = {
	cfg: string;
	onChange: (v: string) => void;
	hideThemeSelector?: boolean;
};

const LOCALSTORAGE_THEME_KEY = "gcsim-config-editor-theme";
const LOCALSTORAGE_FONT_SIZE_KEY = "gcsim-config-editor-font-size";

export function ConfigEditor(props: Props) {
	const { t } = useTranslation();
	const [theme, setTheme] = React.useState<Theme>(() => {
		return localStorage.getItem(LOCALSTORAGE_THEME_KEY) ?? "tomorrow_night";
	});
	const [fontSize, setFontSize] = React.useState(() => {
		return localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY)
			? Number(localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY))
			: 14;
	});
	React.useEffect(() => {
		localStorage.setItem(LOCALSTORAGE_THEME_KEY, theme);
		localStorage.setItem(LOCALSTORAGE_FONT_SIZE_KEY, fontSize.toString());
	}, [theme, fontSize]);
	const hideThemeSelector = props.hideThemeSelector ?? true;
	return (
		<div className="p-1 md:p-2">
			{hideThemeSelector ? null : (
				<div className="my-1 w-full flex flex-col gap-0.5 items-center md:flex-row-reverse md:gap-4 md:items-start">
					<div className="flex items-center gap-2">
						<Label>{t("simple.editor_theme")}</Label>
						<Select value={theme} onValueChange={(v) => setTheme(v)}>
							<SelectTrigger>
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								{themes.map((th) => (
									<SelectItem key={th} value={th}>
										{th}
									</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>
					<div className="flex items-center gap-2">
						<Label>{t("simple.font_size")}</Label>
						<NumberInput value={fontSize} onValueChange={setFontSize} />
					</div>
				</div>
			)}
			<AceEditor
				mode="gcsim"
				theme={theme}
				width="100%"
				onChange={props.onChange}
				value={props.cfg}
				name="config_editor"
				editorProps={{
					$blockScrolling: true,
				}}
				setOptions={{
					maxLines: Infinity,
					fontSize: fontSize,
					tabSize: 2,
					highlightActiveLine: false,
				}}
			/>
		</div>
	);
}
