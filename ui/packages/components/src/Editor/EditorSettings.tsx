import { EDITOR_THEMES } from "@gcsim/editor";
import {
	Button,
	cn,
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@gcsim/primitives";
import { Check, Minus, Plus, SlidersHorizontal } from "lucide-react";
import { useTranslation } from "react-i18next";
import { defaultEditorPrefs, type EditorAppearance } from "./types";

export const MIN_FONT_SIZE = 10;
export const MAX_FONT_SIZE = 28;

export function clampFontSize(size: number) {
	if (!Number.isFinite(size)) return defaultEditorPrefs.fontSize;
	return Math.min(MAX_FONT_SIZE, Math.max(MIN_FONT_SIZE, Math.round(size)));
}

export interface EditorSettingsProps {
	appearance: EditorAppearance;
	onChange: (next: EditorAppearance) => void;
}

export function EditorSettings({ appearance, onChange }: EditorSettingsProps) {
	const { t } = useTranslation();
	const { fontSize, theme } = appearance;
	const setFontSize = (next: number) =>
		onChange({ ...appearance, fontSize: clampFontSize(next) });

	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					variant="ghost"
					size="sm"
					aria-label={t("simple.editor_settings")}
				>
					<SlidersHorizontal />
					<span className="hidden sm:inline">
						{t("simple.editor_settings")}
					</span>
				</Button>
			</PopoverTrigger>
			<PopoverContent
				align="end"
				collisionPadding={16}
				className="scrollbar-surface flex max-h-(--radix-popover-content-available-height) w-[min(18rem,calc(100vw-2rem))] flex-col gap-4 overflow-y-auto"
			>
				<section className="flex items-center justify-between gap-2">
					<span className="text-g-sm font-semibold">
						{t("simple.font_size")}
					</span>
					<div className="flex items-center gap-1">
						<Button
							variant="outline"
							size="icon"
							aria-label={t("simple.decrease_font_size")}
							disabled={fontSize <= MIN_FONT_SIZE}
							onClick={() => setFontSize(fontSize - 1)}
						>
							<Minus />
						</Button>
						<output className="w-8 text-center font-g-mono tabular-nums">
							{fontSize}
						</output>
						<Button
							variant="outline"
							size="icon"
							aria-label={t("simple.increase_font_size")}
							disabled={fontSize >= MAX_FONT_SIZE}
							onClick={() => setFontSize(fontSize + 1)}
						>
							<Plus />
						</Button>
					</div>
				</section>

				<section className="flex flex-col gap-2">
					<span className="text-g-sm font-semibold">
						{t("simple.editor_theme")}
					</span>
					<div role="radiogroup" className="flex flex-col gap-1">
						{EDITOR_THEMES.map(({ id, label, palette }) => {
							const selected = id === theme;
							return (
								<label
									key={id}
									className={cn(
										"flex min-h-11 cursor-pointer items-center gap-3 rounded-g-md border px-2 text-g-sm transition-colors has-focus-visible:ring-2 has-focus-visible:ring-g-accent",
										selected
											? "border-g-accent bg-g-accent/10"
											: "border-transparent hover:bg-g-surface-2",
									)}
								>
									<input
										type="radio"
										name="editor-theme"
										value={id}
										checked={selected}
										onChange={() => onChange({ ...appearance, theme: id })}
										className="sr-only"
									/>
									<span
										aria-hidden
										className="flex h-7 w-12 shrink-0 items-center justify-center gap-1 rounded-g-sm border border-g-line"
										style={{ backgroundColor: palette.background }}
									>
										{[palette.keyword, palette.string, palette.character].map(
											(color) => (
												<span
													key={color}
													className="size-2 rounded-full"
													style={{ backgroundColor: color }}
												/>
											),
										)}
									</span>
									<span className="flex-1">{label}</span>
									{selected ? <Check className="size-4 text-g-accent" /> : null}
								</label>
							);
						})}
					</div>
				</section>
			</PopoverContent>
		</Popover>
	);
}
