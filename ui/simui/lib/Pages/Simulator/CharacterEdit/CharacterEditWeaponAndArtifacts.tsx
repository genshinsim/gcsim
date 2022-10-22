import { Button } from "@blueprintjs/core";
import { ascLvlMax, ascLvlMin, ascToMaxLvl, maxLvlToAsc } from "~/Util";
import { NumberInput } from "~/Components/NumberInput";
import { CharacterEditArtifactSets } from "./CharacterEditArtifactSets";
import React from "react";
import { WeaponSelect } from "~/Components/Select";
import { Trans, useTranslation } from "react-i18next";
import { Character, IWeapon } from "~/Types";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};
export function CharacterEditWeaponAndArtifacts({ char, onChange }: Props) {
  const { t } = useTranslation();

  const [open, setOpen] = React.useState<boolean>(false);

  const handleChangeWeapon = (weapon: IWeapon) => {
    setOpen(false);
    let w = { ...char.weapon };
    w.name = weapon;

    let asc = maxLvlToAsc(w.max_level);
    if (w.level > w.max_level) {
      w.level = w.max_level;
    } else if (w.level < ascLvlMin(asc)) {
      w.level = ascLvlMin(asc);
    }
    char.weapon = w;
    onChange(char);
  };

  const handleChangeWeaponAttr = (key: "refine" | "max_level" | "level") => {
    return (val: number) => {
      let w = { ...char.weapon };
      w[key] = val;

      let asc = maxLvlToAsc(w.max_level);
      if (w.level > w.max_level) {
        w.level = w.max_level;
      } else if (w.level < ascLvlMin(asc)) {
        w.level = ascLvlMin(asc);
      }

      char.weapon = w;
      onChange(char);
    };
  };

  const handleChangeAsc = (val: number) => {
    let w = { ...char.weapon };
    w.max_level = ascToMaxLvl(val);
    if (w.level > w.max_level) {
      w.level = w.max_level;
    }

    let asc = maxLvlToAsc(w.max_level);
    if (w.level > w.max_level) {
      w.level = w.max_level;
    } else if (w.level < ascLvlMin(asc)) {
      w.level = ascLvlMin(asc);
    }

    char.weapon = w;
    onChange(char);
  };

  const asc = maxLvlToAsc(char.weapon.max_level);

  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full hd:basis-36">
        <img
          src={`https://gcsim.app/api/assets/weapons/${char.weapon.name}.png`}
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
      <WeaponSelect isOpen={open} onClose={() => setOpen(false)} onSelect={handleChangeWeapon} />
    </div>
  );
}
