import type { IArtifact } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { artifactLabel, artifactPredicate, artifacts } from "./artifacts";
import { OmniSelect } from "./OmniSelect";
import { renderLabeledItem } from "./utils";

export type ArtifactSelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (artifact: IArtifact) => void;
	value?: IArtifact;
};

export function ArtifactSelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: ArtifactSelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<IArtifact>
			isOpen={isOpen}
			onClose={onClose}
			items={artifacts}
			itemKey={(artifact) => artifact}
			itemPredicate={artifactPredicate}
			itemRenderer={(artifact, state) =>
				renderLabeledItem(artifact, artifactLabel(artifact), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("simple.artifacts")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
