import { dynamicKey, resources } from "@gcsim/localization";
import type { IArtifact } from "@gcsim/types";
import i18n from "i18next";
import { createKeyOrLabelPredicate } from "./utils";

export const artifacts: IArtifact[] = Object.keys(
	resources.en.game.artifact_names,
);

export function artifactLabel(artifact: IArtifact): string {
	return i18n.t(dynamicKey(`game:artifact_names.${artifact}`));
}

export const artifactPredicate = createKeyOrLabelPredicate(artifactLabel);
