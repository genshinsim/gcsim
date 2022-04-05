import { SectionDivider } from "~src/Components/SectionDivider";
import { CharacterEditDetails } from "./CharacterEditDetails";
import { CharacterEditWeaponAndArtifacts } from "./CharacterEditWeaponAndArtifacts";
import { CharacterEditStats } from "./CharacterEditStats";
import { Trans, useTranslation } from "react-i18next";
import { useAppSelector, RootState } from "~src/store";

export function CharacterEdit() {
  useTranslation();

  const { edit_index } = useAppSelector((state: RootState) => {
    return {
      edit_index: state.sim.edit_index,
    };
  });

  if (edit_index === -1) {
    return null;
  }

  return (
    <div className="flex flex-col gap-2">
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.character</Trans>
      </SectionDivider>
      <CharacterEditDetails />
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.weapons_and_artifacts</Trans>
      </SectionDivider>
      <CharacterEditWeaponAndArtifacts />
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.artifact_stats_main</Trans>
      </SectionDivider>
      <div className="p-2">
        <CharacterEditStats />
      </div>
    </div>
  );
}
