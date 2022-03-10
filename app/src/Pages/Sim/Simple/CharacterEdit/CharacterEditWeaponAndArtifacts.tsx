import { Button } from "@blueprintjs/core";
import { simActions } from "~src/Pages/Sim";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { ascLvlMax, ascLvlMin, ascToMaxLvl, maxLvlToAsc } from "~src/util";
import { NumberInput } from "~src/Components/NumberInput";
import { CharacterEditArtifactSets } from "./CharacterEditArtifactSets";
import React from "react";
import { IWeapon, WeaponSelect } from "~src/Components/Weapon";
import { Trans, useTranslation } from "react-i18next";

export function CharacterEditWeaponAndArtifacts() {
  let { t } = useTranslation()

  const { char } = useAppSelector((state: RootState) => {
    return {
      char: state.sim.team[state.sim.edit_index],
    };
  });
  const dispatch = useAppDispatch();

  const [open, setOpen] = React.useState<boolean>(false);

  const handleChangeWeapon = (w: IWeapon) => {
    setOpen(false);
    let next = { ...char.weapon };
    next.name = w.key;
    dispatch(simActions.setCharacterWeapon({ val: next }));
  };

  const handleChangeWeaponAttr = (key: "refine" | "max_level" | "level") => {
    return (val: number) => {
      let next = { ...char.weapon };
      next[key] = val;
      dispatch(simActions.setCharacterWeapon({ val: next }));
    };
  };

  const handleChangeAsc = (val: number) => {
    let next = { ...char.weapon };
    next.max_level = ascToMaxLvl(val);
    if (next.level > next.max_level) {
      next.level = next.max_level;
    }
    dispatch(simActions.setCharacterWeapon({ val: next }));
  };

  let asc = maxLvlToAsc(char.weapon.max_level);

  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full hd:basis-36">
        <img
          src={`/images/weapons/${char.weapon.name}.png`}
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
      <CharacterEditArtifactSets />
      <WeaponSelect
        isOpen={open}
        onClose={() => setOpen(false)}
        onSelect={handleChangeWeapon}
      />
    </div>
  );
}
