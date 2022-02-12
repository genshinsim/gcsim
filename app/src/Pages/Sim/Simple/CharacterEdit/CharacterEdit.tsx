import { SectionDivider } from "~src/Components/SectionDivider";
import { CharacterEditDetails } from "./CharacterEditDetails";
import { CharacterEditWeaponAndArtifacts } from "./CharacterEditWeaponAndArtifacts";
import { CharacterEditStats } from "./CharacterEditStats";

export function CharacterEdit() {
  return (
    <div className="flex flex-col gap-2">
      <SectionDivider fontClass="font-bold text-md">Character</SectionDivider>
      <CharacterEditDetails />
      <SectionDivider fontClass="font-bold text-md">
        Weapons and Artifacts
      </SectionDivider>
      <CharacterEditWeaponAndArtifacts />
      <SectionDivider fontClass="font-bold text-md">Stats</SectionDivider>
      <div className="p-2">
        <CharacterEditStats />
      </div>
    </div>
  );
}
