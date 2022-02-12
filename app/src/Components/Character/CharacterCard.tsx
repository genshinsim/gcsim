import { Button } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
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
} from "~src/Components/Character/Icons";
import { WeaponCard } from "~src/Components/Weapon";
import { CharDetail, CharStatBlock } from "/src/Components/Character";

type Props = {
  char: CharDetail;
  stats: CharStatBlock[];
  statsRows: number;
  className?: string;
  showDelete?: boolean;
  showEdit?: boolean;
  handleDelete?: () => void;
  toggleEdit?: () => void;
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
      return "bg-gradient-to-r from-gray-700 to-green-300";
    case "geo":
      return "bg-gradient-to-r from-gray-700 to-yellow-400";
  }
  return "bg-gray-700";
}

export function CharacterCard({
  char,
  stats,
  statsRows,
  showDelete = false,
  showEdit = false,
  toggleEdit,
  handleDelete,
  className = "",
}: Props) {
  const arts: JSX.Element[] = [];

  for (const key in char.sets) {
    arts.push(
      <div className="w-8 flex flex-col rounded-md" key={key}>
        <Tooltip2 content={key}>
          <img
            key="key"
            src={`/images/artifacts/${key}_flower.png`}
            className="w-full h-8"
          />
        </Tooltip2>

        <span className="text-center text-xs">{char.sets[key]}</span>
      </div>
    );
  }

  let count = 0;
  let rows: JSX.Element[] = [];

  stats.forEach((s, i) => {
    let val: JSX.Element[] = [];
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
        val.push(<td key={"flat-" + i}> </td>);
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
        <td> </td>
        <td> </td>
      </tr>
    );
  }

  return (
    <div className={className}>
      <div className="min-h-24 bg-gray-600 shadow rounded-md text-sm flex flex-col justify-center gap-2">
        <div
          className={
            "character-parent flex flex-row pt-4 pl-4 pr-2 -mt-2 rounded-t-md " +
            charBG(char.element)
          }
        >
          <div className={showDelete ? "absolute top-1 right-1" : "hidden"}>
            <Button icon="cross" intent="danger" small onClick={handleDelete} />
          </div>
          <div className="character-header rounded-t-md" />
          <div className="character-name font-medium m-4 capitalize">
            {char.name} C{char.cons}
          </div>
          <div className="w-1/2 text-sm">
            <div className="rounded-md pl-1 pr-1 mt-6">
              <div>
                Lvl {char.level}/{char.max_level}
              </div>
              <div>
                Talents {char.talents.attack}/{char.talents.skill}/
                {char.talents.burst}
              </div>
              <div className="mt-1 mr-2 grid grid-cols-5">{arts}</div>
            </div>
          </div>
          <div className="w-1/2">
            <img
              src={"/images/avatar/" + char.name + ".png"}
              alt={char.name}
              className="ml-auto h-32 wide:h-auto "
            />
          </div>
        </div>

        <WeaponCard weapon={char.weapon} />

        <div className="ml-2 mr-2 p-2 bg-gray-800 rounded-md">
          <span className="font-bold">Artifact Stats</span>
          <div className="px-2">
            <table className="w-full">
              <tbody>{rows}</tbody>
            </table>
          </div>
        </div>

        <div
          className={
            showEdit ? "ml-auto pl-2 pt-2 pr-2 flex flex-row gap-4" : "hidden"
          }
        >
          <Button icon="edit" onClick={toggleEdit} />
        </div>
        <div className="mb-2" />
      </div>
    </div>
  );
}
