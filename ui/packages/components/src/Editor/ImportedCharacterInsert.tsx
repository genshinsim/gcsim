import { Button, CommandItem } from "@gcsim/primitives";
import { UserPlus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { characterLabel, OmniSelect } from "../common/gcsim";
import { charToCfg } from "./teamConfig";
import type { ImportedCharacterOption } from "./types";

export function insertCharacterBlock(
	config: string,
	option: ImportedCharacterOption,
): string {
	return `${charToCfg(option.character)}\n${config}`;
}

export interface ImportedCharacterInsertProps {
	characters: ImportedCharacterOption[];
	onInsert: (option: ImportedCharacterOption) => void;
}

export function ImportedCharacterInsert({
	characters,
	onInsert,
}: ImportedCharacterInsertProps) {
	const { t } = useTranslation();
	const [open, setOpen] = useState(false);
	const name = (option: ImportedCharacterOption) =>
		characterLabel(option.character.name ?? "") || option.label || option.key;

	return (
		<>
			<Button
				variant="ghost"
				size="sm"
				aria-label={t("simple.imported_characters")}
				onClick={() => setOpen(true)}
			>
				<UserPlus />
				<span className="hidden sm:inline">
					{t("simple.imported_characters")}
				</span>
			</Button>
			<OmniSelect<ImportedCharacterOption>
				isOpen={open}
				onClose={() => setOpen(false)}
				items={characters}
				itemKey={(option) => option.key}
				itemPredicate={(option, query) =>
					`${option.key} ${option.label ?? ""} ${name(option)}`
						.toLowerCase()
						.includes(query.trim().toLowerCase())
				}
				itemRenderer={(option, state) => (
					<CommandItem value={option.key} onSelect={state.onSelect}>
						<span className="flex-1">{name(option)}</span>
						<span className="font-g-mono text-g-xs text-g-ink-mute">
							C{option.character.cons} · Lv{option.character.level}
						</span>
					</CommandItem>
				)}
				onSelect={(option) => {
					setOpen(false);
					onInsert(option);
				}}
				title={t("simple.imported_characters")}
				placeholder={t("db.type_to_search")}
				emptyMessage={t("simple.no_imported_characters")}
			/>
		</>
	);
}
