import { Button } from "@blueprintjs/core";
import { simActions } from "~src/Pages/Sim";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { ascLvlMax, ascLvlMin, ascToMaxLvl, maxLvlToAsc } from "~src/util";
import { NumberInput } from "~src/Components/NumberInput";
import React from "react";
import { CharacterSelect, ICharacter } from "~src/Components/Character";
import { Trans, useTranslation } from "react-i18next";

export function CharacterEditDetails() {
  let { t } = useTranslation();

  const { char, team } = useAppSelector((state: RootState) => {
    return {
      team: state.sim.team,
      char: state.sim.team[state.sim.edit_index],
    };
  });
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);

  const handleChangeCharacter = (w: ICharacter) => {
    setOpen(false);
    //do nothing if this char already exists
    for (let i = 0; i < team.length; i++) {
      if (team[i].name === w.key) {
        return;
      }
    }
    dispatch(
      simActions.setCharacterNameAndEle({ name: w.key, ele: w.element })
    );
  };

  const handleChangeTalent = (key: "attack" | "skill" | "burst") => {
    return (val: number) => {
      let next = { ...char.talents };
      next[key] = val;
      dispatch(simActions.setCharacterTalent({ val: next }));
    };
  };

  const handleChangeAsc = (val: number) => {
    if (val < 0 || val > 6) {
      return;
    }
    const lvl = ascToMaxLvl(val);
    dispatch(simActions.setCharacterMaxLvl({ val: lvl }));
  };

  const handleChangeLvl = (val: number) => {
    if (val <= 0 || val > 90) {
      return;
    }
    dispatch(simActions.setCharacterLvl({ val: val }));
  };

  const handleChangeCons = (val: number) => {
    if (val < 0 || val > 6) {
      return;
    }
    dispatch(simActions.setCharacterCon({ val: val }));
  };

  const asc = maxLvlToAsc(char.max_level);

  const disabled = team.map((c) => c.name);

  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full hd:basis-36">
        <img
          src={"/images/avatar/" + char.name + ".png"}
          alt={char.name}
          className="w-28"
        />
        <Button icon="swap-horizontal" fill onClick={() => setOpen(true)}>
          <Trans>characteredit.change</Trans>
        </Button>
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label={t("characteredit.ascension")}
          onChange={handleChangeAsc}
          value={asc}
          min={0}
          max={6}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.level")}
          onChange={handleChangeLvl}
          value={char.level}
          min={ascLvlMin(asc)}
          max={ascLvlMax(asc)}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.cons")}
          onChange={handleChangeCons}
          value={char.cons}
          integerOnly
        />
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label={t("characteredit.attack")}
          onChange={handleChangeTalent("attack")}
          min={1}
          max={10}
          value={char.talents.attack}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.skill")}
          onChange={handleChangeTalent("skill")}
          min={1}
          max={10}
          value={char.talents.skill}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.burst")}
          onChange={handleChangeTalent("burst")}
          min={1}
          max={10}
          value={char.talents.burst}
          integerOnly
        />
      </div>
      <CharacterSelect
        disabled={disabled}
        isOpen={open}
        onClose={() => setOpen(false)}
        onSelect={handleChangeCharacter}
      />
    </div>
  );
}
