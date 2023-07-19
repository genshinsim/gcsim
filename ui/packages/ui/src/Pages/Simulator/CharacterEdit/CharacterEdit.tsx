import { SectionDivider } from "../../../Components/SectionDivider";
import { CharacterEditDetails } from "./CharacterEditDetails";
import { CharacterEditWeaponAndArtifacts } from "./CharacterEditWeaponAndArtifacts";
import { CharacterEditStats } from "./CharacterEditStats";
import { Trans, useTranslation } from "react-i18next";
import { Character } from "@gcsim/types";

export function CharacterEdit({
  char,
  index,
  handleChange,
}: {
  char: Character | null;
  index: number;
  handleChange: (char: Character, index: number) => void;
}) {
  useTranslation();

  if (char === null) {
    return null;
  }

  const onChange = (char: Character) => {
    handleChange(char, index);
  };

  return (
    <div className="flex flex-col gap-2">
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.character</Trans>
      </SectionDivider>
      <CharacterEditDetails char={char} onChange={onChange} />
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.weapons_and_artifacts</Trans>
      </SectionDivider>
      <CharacterEditWeaponAndArtifacts char={char} onChange={onChange} />
      <SectionDivider fontClass="font-bold text-md">
        <Trans>characteredit.artifact_stats_main</Trans>
      </SectionDivider>
      <div className="p-2">
        <CharacterEditStats char={char} onChange={onChange} />
      </div>
    </div>
  );
}
