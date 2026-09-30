import { useTranslation } from "react-i18next";
import { AceEditorWrapper, themes } from "./AceEditorWrapper";
import type { EditorAppearance, EditorProps } from "./types";

export const Editor = ({
	value,
	onChange,
	theme = "tomorrow_night",
	fontSize = 14,
	onAppearanceChange,
	maxLines,
}: EditorProps) => {
	const { t } = useTranslation();
	const updateAppearance = (patch: Partial<EditorAppearance>) =>
		onAppearanceChange?.({ theme, fontSize, ...patch });

	return (
		<div className="flex flex-col">
			{onAppearanceChange ? (
				<div className="flex flex-wrap items-center justify-end gap-4">
					<label className="flex items-center gap-2">
						{t("simple.font_size")}
						<input
							type="number"
							value={fontSize}
							onChange={(e) =>
								updateAppearance({ fontSize: Number(e.currentTarget.value) })
							}
						/>
					</label>
					<label className="flex items-center gap-2">
						{t("simple.editor_theme")}
						<select
							value={theme}
							onChange={(e) =>
								updateAppearance({ theme: e.currentTarget.value })
							}
						>
							{themes.map((th) => (
								<option key={th} value={th}>
									{th}
								</option>
							))}
						</select>
					</label>
				</div>
			) : null}
			<AceEditorWrapper
				value={value}
				onChange={onChange}
				theme={theme}
				fontSize={fontSize}
				maxLines={maxLines}
			/>
		</div>
	);
};
