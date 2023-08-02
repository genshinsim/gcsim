import { Button } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
import {
  IconAnemo,
  IconAtk,
  IconCD,
  IconCR,
  IconCryo,
  IconDef,
  IconDendro,
  IconElectro,
  IconEM,
  IconER,
  IconGeo,
  IconHeal,
  IconHP,
  IconHydro,
  IconPhysical,
  IconPyro,
} from "../Icons";
import { WeaponCard } from "./WeaponCard";
import { Trans, useTranslation } from "react-i18next";
import { Character } from "@gcsim/types";
import { CharStatBlock } from "../../Pages/Simulator/Components/character";
import classNames from "classnames";
import placeholder from "../Images/default.png";
import React from "react";

type Props = {
  char: Character;
  stats: CharStatBlock[];
  statsRows: number;
  className?: string;
  showDetails?: boolean;
  viewerMode? : boolean; //hide the delete and edit button; add toggle for showing stats
  isSkeleton?: boolean;
  handleDelete?: () => void;
  toggleEdit?: () => void;
  handleToggleDetail?: () => void;
};

function statKeyToIcon(key: string): JSX.Element {
  switch (key) {
    case "hp":
      return <IconHP />;
    case "atk":
      return <IconAtk />;
    case "def":
      return <IconDef />;
    case "er":
      return <IconER />;
    case "em":
      return <IconEM />;
    case "cr":
      return <IconCR />;
    case "cd":
      return <IconCD />;
    case "electro":
      return <IconElectro />;
    case "pyro":
      return <IconPyro />;
    case "cryo":
      return <IconCryo />;
    case "hydro":
      return <IconHydro />;
    case "geo":
      return <IconGeo />;
    case "anemo":
      return <IconAnemo />;
    case "phys":
      return <IconPhysical />;
    case "heal":
      return <IconHeal />;
    case "dendro":
      return <IconDendro />;
    default:
      return <span />;
  }
}

function charBG(element: string) {
  switch (element) {
    case "cryo":
      return "bg-gradient-to-r from-gray-700 to-blue-300";
    case "hydro":
      return "bg-gradient-to-r from-gray-700 to-blue-500";
    case "pyro":
      return "bg-gradient-to-r from-gray-700 to-red-400";
    case "electro":
      return "bg-gradient-to-r from-gray-700 to-purple-300";
    case "anemo":
      return "bg-gradient-to-r from-gray-700 to-teal-500";
    case "dendro":
      return "bg-gradient-to-r from-gray-700 to-lime-700";
    case "geo":
      return "bg-gradient-to-r from-gray-700 to-yellow-400";
  }
  return "bg-gray-700";
}

export function CharacterCard({
  char,
  stats,
  statsRows,
  showDetails = true,
  viewerMode = false,
  isSkeleton,
  toggleEdit,
  handleDelete,
  handleToggleDetail,
  className = "",
}: Props) {
  const { t } = useTranslation();

  const arts: JSX.Element[] = [];

  for (const key in char.sets) {
    arts.push(
      <div className="w-8 flex flex-col rounded-md" key={key}>
        <Tooltip2 content={key}>
          <img
            key="key"
            src={`/api/assets/artifacts/${key}_flower.png`}
            className="w-full h-8"
            onError={(e) => (e.target as HTMLImageElement).src = placeholder}
          />
        </Tooltip2>

        <span className="text-center text-xs">{char.sets[key]}</span>
      </div>
    );
  }

  let count = 0;
  const rows: JSX.Element[] = [];

  stats.forEach((s, i) => {
    const val: JSX.Element[] = [];
    if (s.flat === 0 && s.percent === 0) {
      return;
    }

    count++;

    switch (s.t) {
      case "both":
        val.push(
          <td key={"flat-" + i} className="text-right">
            {s.flat.toFixed(0)}
          </td>
        );
        val.push(
          <td key={"per-" + i} className="text-right">
            {(s.percent * 100).toFixed(2) + "%"}
          </td>
        );
        break;
      case "f":
        val.push(
          <td key={"flat-" + i} className="text-right">
            {s.flat.toFixed(0)}
          </td>
        );
        val.push(<td key={"per-" + i}></td>);
        break;
      case "%":
        val.push(<td key={"flat-" + i}></td>);
        val.push(
          <td key={"per-" + i} className="text-right">
            {(s.percent * 100).toFixed(2) + "%"}
          </td>
        );
    }

    rows.push(
      <tr key={count}>
        <td className="flex flex-row place-items-center">
          <div className="w-4 mr-1 fill-gray-100">{statKeyToIcon(s.key)}</div>{" "}
          {s.name}
        </td>
        {val}
      </tr>
    );
  });

  for (; count < statsRows; count++) {
    rows.push(
      <tr key={count + 1}>
        <td>
          <br />
        </td>
        <td></td>
        <td></td>
      </tr>
    );
  }

  const skeleton = classNames({ ["bp4-skeleton"]: isSkeleton });

  return (
    <div className={className}>
      <div className="min-h-24 bg-bp4-dark-gray-400 shadow text-sm flex flex-col justify-center gap-2 border border-gray-600">
        <div
          className={
            "character-parent flex flex-row pt-4 pl-4 pr-2 " +
            charBG(char.element)
          }
        >
          <div className="absolute top-1 right-1">
            {
              viewerMode ?
              <Button icon={showDetails ? "caret-up" : "caret-down"} small onClick={handleToggleDetail} />
              :
              <Button icon="cross" intent="danger" small onClick={handleDelete} />
            }
          </div>
          <div className="character-header"></div>
          <div className={"character-name font-medium m-4 capitalize " + skeleton}>
            <>
              <Trans>character.c_pre</Trans>
              {char.cons ? char.cons : 0}
              <Trans>character.c_post</Trans>{" "}
              {t(
                `game:character_names.${char.name}`
              )}{" "}
            </>
          </div>
          <div className="w-1/2 text-sm">
            <div className={"pl-1 pr-1 mt-6 " + skeleton}>
              <div>
                <Trans>character.lvl</Trans> {char.level}/{char.max_level}
              </div>
              <div>
                <Trans>character.talents</Trans> {char.talents.attack}/
                {char.talents.skill}/{char.talents.burst}
              </div>
              <div className="mt-1 mr-2 grid grid-cols-5">{arts}</div>
            </div>
          </div>
          <div className="w-1/2 h-32">
            {isSkeleton ? null : <img
              src={
                "https://gcsim.app/api/assets/avatar/" +
                char.name +
                ".png"
              }
              alt={char.name}
              className="ml-auto h-32"
              onError={(e) => (e.target as HTMLImageElement).src = placeholder}
            />}
          </div>
        </div>

        <WeaponCard weapon={char.weapon} isSkeleton={isSkeleton} />

        {showDetails ? (
            <div className="ml-2 mr-2 p-2 bg-[#252A31] border-gray-600 border">
            <span className="font-bold">
              <Trans>character.artifact_stats</Trans>
            </span>
            <div className="px-2">
              <table className="w-full">
                <tbody>{rows}</tbody>
              </table>
            </div>
          </div>
        ) :  null
        }

        {
          viewerMode ?
          null
          :
        <div
          className="ml-auto pl-2 pt-2 pr-2 flex flex-row gap-4" 
        >
          <Button icon="edit" onClick={toggleEdit} />
        </div>
        }

        <div className="" />
      </div>
    </div>
  );
}
