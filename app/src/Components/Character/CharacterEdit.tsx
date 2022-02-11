import { Tab, Tabs } from "@blueprintjs/core";
import { CharacterEditStats, CharacterEditWeapon, CharDetail } from ".";
import { SectionDivider } from "../SectionDivider";

type Props = {
  char: CharDetail;
  onChange: (char: CharDetail) => void;
};

export function CharacterEdit({ char, onChange }: Props) {
  const handleOnStatChange = (index: number, value: number) => {
    char.stats[index] = value;
    onChange(char);
  };
  return (
    <div className="flex flex-col gap-2">
      <SectionDivider fontClass="font-bold text-md">Character</SectionDivider>
      <SectionDivider fontClass="font-bold text-md">Weapons</SectionDivider>
      <CharacterEditWeapon />
      <SectionDivider fontClass="font-bold text-md">
        Artifact Sets
      </SectionDivider>
      <div className="p-2"></div>
      <SectionDivider fontClass="font-bold text-md">Stats</SectionDivider>
      <div className="p-2">
        <CharacterEditStats char={char} onChange={handleOnStatChange} />
      </div>
    </div>
  );
}
