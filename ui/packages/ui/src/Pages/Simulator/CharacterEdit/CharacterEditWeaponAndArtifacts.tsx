import { Button } from "@blueprintjs/core";
import { ascLvlMax, ascLvlMin, ascToMaxLvl, maxLvlToAsc } from "../Components/util";
import { NumberInput } from "../../../Components/NumberInput";
import { CharacterEditArtifactSets } from "./CharacterEditArtifactSets";
import React from "react";
import { WeaponSelect } from "../../../Components/Select";
import { Trans, useTranslation } from "react-i18next";
import { Character, IWeapon } from "@gcsim/types";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};
export function CharacterEditWeaponAndArtifacts({ char, onChange }: Props) {
  const { t } = useTranslation();

  const [open, setOpen] = React.useState<boolean>(false);

  const weapLvlCheck = (level: number, max_level: number): number => {
    const asc = maxLvlToAsc(max_level);
    if (level > max_level) {
      return max_level;
    } else if (level < ascLvlMin(asc)) {
      return ascLvlMin(asc);
    }
    return 1;
  };

  const handleChangeWeapon = (weapon: IWeapon) => {
    setOpen(false);
    const next = JSON.parse(JSON.stringify(char));
    next.weapon.name = weapon;

    next.level = weapLvlCheck(next.level, next.max_level);
    onChange(next);
  };

  const handleChangeWeaponAttr = (key: "refine" | "max_level" | "level") => {
    return (val: number) => {
      const next = JSON.parse(JSON.stringify(char));
      next.weapon[key] = val;
      next.level = weapLvlCheck(next.level, next.max_level);
      onChange(next);
    };
  };

  const handleChangeAsc = (val: number) => {
    const next = JSON.parse(JSON.stringify(char));
    next.weapon.max_level = ascToMaxLvl(val);
    next.level = weapLvlCheck(next.level, next.max_level);
    onChange(next);
  };

  const asc = maxLvlToAsc(char.weapon.max_level);

  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full hd:basis-36">
        <img
          src={`/api/assets/weapons/${char.weapon.name}.png`}
          alt={char.weapon.name}
          className="w-28 "
        />
        <Button
          icon="swap-horizontal"
          fill
          onClick={() => {
            setOpen(true);
          }}
        >
          <Trans>characteredit.change</Trans>
        </Button>
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col gap-y-2 ">
        <NumberInput
          label={t("characteredit.refine")}
          onChange={handleChangeWeaponAttr("refine")}
          value={char.weapon.refine}
          min={1}
          max={5}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.ascension")}
          onChange={handleChangeAsc}
          value={asc}
          integerOnly
          min={0}
          max={6}
        />
        <NumberInput
          label={t("characteredit.level")}
          onChange={handleChangeWeaponAttr("level")}
          value={char.weapon.level}
          integerOnly
          min={ascLvlMin(asc)}
          max={ascLvlMax(asc)}
        />
      </div>
      <CharacterEditArtifactSets onChange={onChange} char={char} />
      <WeaponSelect
        isOpen={open}
        onClose={() => setOpen(false)}
        onSelect={handleChangeWeapon}
      />
    </div>
  );
}
