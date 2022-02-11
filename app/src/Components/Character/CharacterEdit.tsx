import { Tab, Tabs } from "@blueprintjs/core";
import { CharacterEditStats, CharacterEditWeapon, CharDetail } from ".";

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
    <div>
      <Tabs>
        <Tab
          id="stats"
          title="Stats"
          panel={
            <CharacterEditStats char={char} onChange={handleOnStatChange} />
          }
        />
        <Tab
          id="character"
          title="Character"
          panel={<CharacterEditWeapon />}
        ></Tab>
        <Tab id="weapon" title="Weapon" panel={<CharacterEditWeapon />}></Tab>
        <Tab id="sets" title="Sets" panel={<CharacterEditWeapon />}></Tab>
      </Tabs>
    </div>
  );
}
