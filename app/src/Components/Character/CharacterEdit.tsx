import { Character } from "~src/types";
import { CharacterEditStats, CharacterEditWeapon } from ".";
import { SectionDivider } from "../SectionDivider";
import { CharacterEditArtifactSets } from "./CharacterEditArtifactSets";
import { CharacterEditDetails } from "./CharacterEditDetails";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

function maxLvlToAsc(lvl: number): number {
  switch (lvl) {
    case 90:
      return 6;
    case 80:
      return 5;
    case 70:
      return 4;
    case 60:
      return 3;
    case 50:
      return 2;
    case 40:
      return 1;
    default:
      return 0;
  }
}

export function CharacterEdit({ char, onChange }: Props) {
  const handleOnStatChange = (index: number, value: number) => {
    console.log("before: ", char.stats);
    char.stats[index] = value;
    console.log("changing: ", index, value);
    console.log("after: ", char.stats);

    onChange(char);
  };
  return (
    <div className="flex flex-col gap-2">
      <SectionDivider fontClass="font-bold text-md">Character</SectionDivider>
      <CharacterEditDetails char={char} onChange={onChange} />
      <SectionDivider fontClass="font-bold text-md">Weapons</SectionDivider>
      <CharacterEditWeapon char={char} onChange={onChange} />
      <SectionDivider fontClass="font-bold text-md">
        Artifact Sets
      </SectionDivider>
      <div className="p-2">
        <CharacterEditArtifactSets char={char} onChange={onChange} />
      </div>
      <SectionDivider fontClass="font-bold text-md">Stats</SectionDivider>
      <div className="p-2">
        <CharacterEditStats char={char} onChange={handleOnStatChange} />
      </div>
    </div>
  );
}
