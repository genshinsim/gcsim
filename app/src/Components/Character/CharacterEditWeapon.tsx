import { Button } from "@blueprintjs/core";
import React from "react";
import { Character } from "~src/types";
import { NumberInput } from "../NumberInput";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

export function CharacterEditWeapon({ char, onChange }: Props) {
  const [lvl, setlvl] = React.useState<number>(10);
  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full wide:basis-36">
        <img
          src={`/images/weapons/${char.weapon.name}.png`}
          alt={char.weapon.name}
          className="w-28 "
        />
        <Button icon="swap-horizontal" fill>
          Change
        </Button>
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 wide:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label="Refine"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
        <NumberInput
          label="Ascension"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
        <NumberInput
          label="Level"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
      </div>
    </div>
  );
}
