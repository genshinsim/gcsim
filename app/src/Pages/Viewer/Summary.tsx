import Graphs from "./Graphs/Graphs";
import { SimResults } from "./DataType";
import TeamView from "./Team/TeamView";
import { HTMLTable } from "@blueprintjs/core";

const DEFP = 1;
const DEF = 2;
const HP = 3;
const HPP = 4;
const ATK = 5;
const ATKP = 6;
const ER = 7;
const EM = 8;
const CR = 9;
const CD = 10;
const Heal = 11;
const PyroP = 12;
const HydroP = 13;
const CryoP = 14;
const ElectroP = 15;
const AnemoP = 16;
const GeoP = 17;
const PhysP = 18;
const DendroP = 19;

export default function Summary({ data }: { data: SimResults }) {
  const chars = data.char_names.map((e) => {
    return (
      <div key={e}>
        <img src={"/images/avatar/" + e + ".png"} className="w-full h-auto" />
      </div>
    );
  });

  //calculate per target damage
  let trgs: JSX.Element[] = [];

  for (const key in data.dps_by_target) {
    trgs.push(
      <div className="w-full flex flex-row" key={key}>
        <span className="w-24">
          <span className="pl-2" />
          {`Target ${key}:`}
        </span>
        <div className="grid grid-cols-4 grow">
          <span className="text-right">-</span>
          <span className="text-right">
            {data.dps_by_target[key].mean.toLocaleString(undefined, {
              maximumFractionDigits: 0,
              minimumFractionDigits: 0,
            })}
          </span>
          <span className="text-right">
            {(
              (100 * data.dps_by_target[key].mean) /
              data.dps.mean
            ).toLocaleString(undefined, {
              maximumFractionDigits: 0,
              minimumFractionDigits: 0,
            })}
            {"%"}
          </span>
          <span className="text-right">
            {data.dps_by_target[key].sd
              ? data.dps_by_target[key].sd!.toLocaleString(undefined, {
                  maximumFractionDigits: 0,
                  minimumFractionDigits: 0,
                })
              : "-"}
          </span>
        </div>
      </div>
    );
  }

  console.log(data);

  let statRows: {
    [key: string]: {
      name: string;
      val: JSX.Element[];
      flat: number;
      per: number;
      count: number;
      t: string;
    };
  } = {
    hp: { name: "hp / hp%", flat: HP, per: HPP, val: [], count: 0, t: "both" },
    atk: {
      name: "atk / atk%",
      flat: ATK,
      per: ATKP,
      val: [],
      count: 0,
      t: "both",
    },
    def: {
      name: "def / def%",
      flat: DEF,
      per: DEFP,
      val: [],
      count: 0,
      t: "both",
    },
    em: { name: "em", flat: EM, per: -1, val: [], count: 0, t: "f" },
    er: { name: "er", flat: -1, per: ER, val: [], count: 0, t: "%" },
    cr: { name: "cr", flat: -1, per: CR, val: [], count: 0, t: "%" },
    cd: { name: "cd", flat: -1, per: CD, val: [], count: 0, t: "%" },
    electro: {
      name: "electro%",
      flat: -1,
      per: ElectroP,
      val: [],
      count: 0,
      t: "%",
    },
    pyro: { name: "pyro%", flat: -1, per: PyroP, val: [], count: 0, t: "%" },
    cryo: { name: "cryo%", flat: -1, per: CryoP, val: [], count: 0, t: "%" },
    hydro: { name: "hydro%", flat: -1, per: HydroP, val: [], count: 0, t: "%" },
    geo: { name: "geo%", flat: -1, per: GeoP, val: [], count: 0, t: "%" },
    anemo: { name: "anemo%", flat: -1, per: AnemoP, val: [], count: 0, t: "%" },
    phys: { name: "phys%", flat: -1, per: PhysP, val: [], count: 0, t: "%" },
    heal: { name: "heal", flat: -1, per: Heal, val: [], count: 0, t: "%" },
  };

  let header: JSX.Element[] = [];

  const charStats = data.char_details.map((char, i) => {
    for (const key in statRows) {
      let s = statRows[key];
      switch (s.t) {
        case "both":
          statRows[key].val.push(
            <td key={char.name} className="text-right">
              {char.stats[s.flat].toFixed(0)}
            </td>
          );
          statRows[key].val.push(
            <td key={char.name + "%"} className="text-right">
              {(char.stats[s.per] * 100).toFixed(2) + "%"}
            </td>
          );

          if (char.stats[s.per] > 0 || char.stats[s.flat] > 0) {
            statRows[key].count++;
          }

          break;
        case "f":
          statRows[key].val.push(
            <td key={char.name} className="text-right">
              {char.stats[s.flat].toFixed(0)}
            </td>
          );
          statRows[key].val.push(<td key={char.name + "%"}></td>);
          if (char.stats[s.flat] > 0) {
            statRows[key].count++;
          }
          break;
        case "%":
          statRows[key].val.push(
            <td key={char.name} className="text-right"></td>
          );
          statRows[key].val.push(
            <td key={char.name + "%"} className="text-right">
              {(char.stats[s.per] * 100).toFixed(2) + "%"}
            </td>
          );
          if (char.stats[s.per] > 0) {
            statRows[key].count++;
          }
          break;
      }
    }
    header.push(
      <th key={i} className="capitalize" colSpan={2}>
        {char.name}
      </th>
    );
  });

  let rows: JSX.Element[] = [];

  for (const key in statRows) {
    const r = statRows[key];
    if (r.count > 0) {
      rows.push(
        <tr key={key}>
          <td>{r.name}</td>
          {r.val}
        </tr>
      );
    }
  }

  return (
    <div>
      <TeamView team={data.char_details} />
      <div className="m-2 grid grid-cols-2 gap-2">
        <div className="rounded-md p-2 bg-gray-600">
          <span className="font-bold">
            Character Stats From Artifacts (Main + Subs Only)
          </span>
          <div className="m-4">
            <table className="w-full p-4">
              <thead className="border-b-2">
                <tr>
                  <th className="w-24"></th>
                  {header}
                </tr>
              </thead>
              <tbody>{rows}</tbody>
            </table>
          </div>
        </div>
        <div className="rounded-md bg-gray-600 flex flex-col gap-1 p-4">
          <div className="w-full">
            <span>
              Simulated{" "}
              {data.sim_duration.mean.toLocaleString(undefined, {
                maximumFractionDigits: 2,
              })}{" "}
              sec of combat ({data.iter} iterations took{" "}
              {(data.runtime / 1000000000).toFixed(3)} seconds to run)
            </span>
          </div>
          <div className="max-w-4xl w-full pl-4 flex flex-col gap-1">
            <div className="flex flex-row border-solid border-b-2 font-bold">
              <span className="w-24">Target</span>
              <div className="grid grid-cols-4 grow">
                <span className="text-right">Level</span>
                <span className="text-right">Avg DPS</span>
                <span className="text-right">%</span>
                <span className="text-right">Std. Dev.</span>
              </div>
            </div>
            {trgs}

            <div className="w-full flex flex-row border-solid border-t-2 font-bold">
              <span className="w-24">Combined</span>
              <div className="grid grid-cols-4 grow">
                <span className="text-right"></span>
                <span className="text-right">
                  {" "}
                  {data.dps.mean.toLocaleString(undefined, {
                    maximumFractionDigits: 0,
                    minimumFractionDigits: 0,
                  })}
                </span>
                <span className="text-right"></span>
                <span className="text-right">
                  {data.dps.sd?.toLocaleString(undefined, {
                    maximumFractionDigits: 0,
                    minimumFractionDigits: 0,
                  })}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <Graphs data={data} />
    </div>
  );
}
