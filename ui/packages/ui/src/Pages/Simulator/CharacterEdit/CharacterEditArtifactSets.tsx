import { Button, Checkbox } from "@blueprintjs/core";
import React from "react";
import { ArtifactSelect } from "../../../Components/Select";
import { Trans, useTranslation } from "react-i18next";
import { Character, IArtifact } from "@gcsim/types";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

export function CharacterEditArtifactSets({ char, onChange }: Props) {
  useTranslation();

  const [open, setOpen] = React.useState<boolean>(false);

  const handleAddSet = (set: IArtifact) => {
    //close the display
    setOpen(false);
    //make sure it doesn't exist already
    if (set in char.sets) {
      return;
    }
    const next = JSON.parse(JSON.stringify(char));
    next.sets[set] = 0;
    onChange(next);
  };

  const handleDeleteSetBonus = (set: string) => {
    return () => {
      const next = JSON.parse(JSON.stringify(char));
      delete next.sets[set];
      onChange(next);
    };
  };

  const handleChangeSetBonus = (set: string, bonus: 2 | 4) => {
    return () => {
      const next = JSON.parse(JSON.stringify(char));
      const current = next.sets[set];
      switch (bonus) {
        case 2:
          //if i click on 2 and current >=2 then set it to 0
          if (current >= 2) {
            next.sets[set] = 0;
          } else {
            //otherwise set to 2
            next.sets[set] = 2;
          }
          onChange(next);
          break;
        case 4:
          //if current is already 4 then set to 2
          if (current >= 4) {
            next.sets[set] = 2;
          } else {
            next.sets[set] = 4;
            //otherwise set to 4
          }
          onChange(next);
          break;
      }
    };
  };

  //this is the total number of artifact set bonuses
  let total = 0;
  for (const key in char.sets) {
    total += char.sets[key];
  }

  const checkDisabled = (key: string, bonus: 2 | 4): boolean => {
    console.log(
      "set: " +
        key +
        "bonus for " +
        bonus +
        " total ticked: " +
        total +
        " in set: " +
        char.sets[key] +
        " check: " +
        (total + bonus - char.sets[key])
    );
    return total + bonus - char.sets[key] > 4 && char.sets[key] < bonus;
  };

  const arts: JSX.Element[] = [];

  for (const key in char.sets) {
    arts.push(
      <div
        key={key}
        className="basis-full sm:basis-320 rounded-md bg-gray-700 flex flex-row place-items-center pl-1 pr-2"
      >
        <img
          key="key"
          src={`/api/assets/artifacts/${key}_flower.png`}
          className="w-12"
        />
        <span className="font-bold">
          <Trans>characteredit.set_bonus</Trans>
        </span>
        <div className="flex flex-row gap-2 flex-grow justify-center">
          <Checkbox
            large
            style={{ marginBottom: 0 }}
            checked={char.sets[key] >= 2}
            onClick={handleChangeSetBonus(key, 2)}
            disabled={checkDisabled(key, 2)}
          >
            2
          </Checkbox>
          <Checkbox
            large
            style={{ marginBottom: 0 }}
            checked={char.sets[key] >= 4}
            onClick={handleChangeSetBonus(key, 4)}
            disabled={checkDisabled(key, 4)}
          >
            4
          </Checkbox>
        </div>
        <Button
          className="ml-auto"
          icon="trash"
          intent="danger"
          onClick={handleDeleteSetBonus(key)}
        />
      </div>
    );
  }

  return (
    <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col place-items-center">
      <div className="flex flex-row flex-wrap gap-2 justify-center w-full">
        {arts}
      </div>
      <div className="mt-2 w-full xs:w-[25rem]">
        <Button icon="add" fill intent="success" onClick={() => setOpen(true)}>
          <Trans>characteredit.add_set_bonus</Trans>
        </Button>
      </div>
      <ArtifactSelect
        isOpen={open}
        onClose={() => setOpen(false)}
        onSelect={handleAddSet}
      />
    </div>
  );
}
