import { Button, FormGroup, Switch } from "@blueprintjs/core";
import { Character } from "~/src/types";
import { StatRow } from "./CharacterEditStatRow";

import {
  IconAnemo,
  IconAtk,
  IconCD,
  IconCR,
  IconCryo,
  IconDef,
  IconElectro,
  IconEM,
  IconER,
  IconGeo,
  IconHeal,
  IconHP,
  IconHydro,
  IconPhysical,
  IconPyro,
} from "~src/Components/Icons";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "~src/Pages/Sim";
import { StatToIndexMap } from "~src/util";
import { useTranslation } from "react-i18next";

export type subDisplayLine = {
  stat?: string;
  stat_?: string;
  label: string;
  val: number;
  val_: number;
  icon: React.ReactElement;
};

type StatRowsProp = {
  stats: number[];
  onChange: (index: number, value: number) => void;
};

function StatRows(props: StatRowsProp) {
  let { t } = useTranslation()

  let subs: subDisplayLine[] = [
    {
      stat: "HP",
      stat_: "HPP",
      label: t("characteredit.hp_hp"),
      val_: 0,
      val: 0,
      icon: <IconHP className="fill-gray-100" />,
    },
    {
      stat: "ATK",
      stat_: "ATKP",
      label: t("characteredit.atk_atk"),
      val: 0,
      val_: 0,
      icon: <IconAtk className="fill-gray-100" />,
    },
    {
      stat: "DEF",
      stat_: "DEFP",
      label: t("characteredit.def_def"),
      val_: 0,
      val: 0,
      icon: <IconDef className="fill-gray-100" />,
    },
    {
      stat: "EM",
      label: t("characteredit.em"),
      val_: 0,
      val: 0,
      icon: <IconEM className="fill-gray-100" />,
    },
    {
      stat_: "ER",
      label: t("characteredit.er"),
      val_: 0,
      val: 0,
      icon: <IconER className="fill-gray-100" />,
    },
    {
      stat_: "CR",
      label: t("characteredit.cr"),
      val_: 0,
      val: 0,
      icon: <IconCR className="fill-gray-100" />,
    },
    {
      stat_: "CD",
      label: t("characteredit.cd"),
      val_: 0,
      val: 0,
      icon: <IconCD className="fill-gray-100" />,
    },
  ];

  let eleSubs: subDisplayLine[] = [
    {
      stat_: "Heal",
      label: t("characteredit.heal"),
      val_: 0,
      val: 0,
      icon: <IconHeal className="fill-gray-100" />,
    },
    {
      stat_: "PyroP",
      label: t("characteredit.pyro"),
      val_: 0,
      val: 0,
      icon: <IconPyro className="fill-gray-100" />,
    },
    {
      stat_: "HydroP",
      label: t("characteredit.hydro"),
      val_: 0,
      val: 0,
      icon: <IconHydro className="fill-gray-100" />,
    },
    {
      stat_: "CryoP",
      label: t("characteredit.cryo"),
      val_: 0,
      val: 0,
      icon: <IconCryo className="fill-gray-100" />,
    },
    {
      stat_: "ElectroP",
      label: t("characteredit.electro"),
      val_: 0,
      val: 0,
      icon: <IconElectro className="fill-gray-100" />,
    },
    {
      stat_: "AnemoP",
      label: t("characteredit.anemo"),
      val_: 0,
      val: 0,
      icon: <IconAnemo className="fill-gray-100" />,
    },
    {
      stat_: "GeoP",
      label: t("characteredit.geo"),
      val_: 0,
      val: 0,
      icon: <IconGeo className="fill-gray-100" />,
    },
    {
      stat_: "PhyP",
      label: t("characteredit.phy"),
      val_: 0,
      val: 0,
      icon: <IconPhysical className="fill-gray-100" />,
    },
    // {
    //   stat_: "DendroP",
    //   label: "Dendro%",
    //   val_: 0,
    //   val: 0,
    //   icon: <IconHP className="fill-gray-100" />,
    // },
  ];

  const rows = subs.map((sub, index) => {
    //find the stats
    if (sub.stat) {
      sub.val = props.stats[StatToIndexMap[sub.stat]];
    }
    if (sub.stat_) {
      sub.val_ =
        Math.round(props.stats[StatToIndexMap[sub.stat_]] * 10000) / 100;
    }
    return <StatRow key={index} sub={sub} onChange={props.onChange} />;
  });

  const eleRows = eleSubs.map((sub, index) => {
    if (sub.stat_) {
      sub.val_ =
        Math.round(props.stats[StatToIndexMap[sub.stat_]] * 10000) / 100;
    }
    return <StatRow key={index} sub={sub} onChange={props.onChange} />;
  });

  return (
    <div className="flex flex-row flex-wrap">
      <div className="basis-full hd:basis-1/2 pl-2 pr-2 ">{rows}</div>
      <div className="basis-full hd:basis-1/2 pl-2 pr-2">{eleRows}</div>
    </div>
  );
}

export function CharacterEditStats() {
  const { char } = useAppSelector((state: RootState) => {
    return {
      char: state.sim.team[state.sim.edit_index],
    };
  });
  const dispatch = useAppDispatch();

  const handleChangeStat = (index: number, value: number) => {
    dispatch(simActions.setCharacterStats({ index: index, val: value }));
  };

  return (
    <div className="flex flex-col">
      <StatRows stats={char.stats} onChange={handleChangeStat} />
    </div>
  );
}
