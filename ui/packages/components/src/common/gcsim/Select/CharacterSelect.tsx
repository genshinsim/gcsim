import type { ICharacter } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { characterLabel, characterPredicate, characters } from "./characters";
import { OmniSelect } from "./OmniSelect";
import { renderLabeledItem } from "./utils";

export type CharacterSelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (character: ICharacter) => void;
	value?: ICharacter;
};

export function CharacterSelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: CharacterSelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<ICharacter>
			isOpen={isOpen}
			onClose={onClose}
			items={characters}
			itemKey={(character) => character}
			itemPredicate={characterPredicate}
			itemRenderer={(character, state) =>
				renderLabeledItem(character, characterLabel(character), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("db.characters")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
