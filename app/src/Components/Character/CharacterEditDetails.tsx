import { Button } from "@blueprintjs/core";
import React from "react";
import { Character } from "~/src/types";
import { NumberInput } from "../NumberInput";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

export function CharacterEditDetails({ char, onChange }: Props) {
  const [lvl, setlvl] = React.useState<number>(10);
  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full wide:basis-36">
        <img
          src={"/images/avatar/" + char.name + ".png"}
          alt={char.name}
          className="w-28"
        />
        <Button icon="swap-horizontal" fill>
          Change
        </Button>
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 wide:basis-0 flex flex-col gap-y-2">
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
        <NumberInput
          label="Cons"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 wide:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label="Attack"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
        <NumberInput
          label="Skill"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
        <NumberInput
          label="Burst"
          onChange={(v) => setlvl(v)}
          value={lvl}
          integerOnly
        />
      </div>
    </div>
  );
}
