import Graphs from "./Graphs/Graphs";
import { SimResults } from "./DataType";
import TeamView from "./Team/TeamView";
import { Button, HTMLTable } from "@blueprintjs/core";
import DPSOverTime from "./Graphs/DPSOverTime";

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

  return (
    <div className="wide:w-[70rem] ml-auto mr-auto">
      <TeamView team={data.char_details} />
      <div className="bg-gray-600 relative rounded-md p-2 m-2 pt-10">
        <DPSOverTime data={data} />
        <div className="w-full text-center">
          Simulated{" "}
          {data.sim_duration.mean.toLocaleString(undefined, {
            maximumFractionDigits: 2,
          })}{" "}
          sec of combat ({data.iter} iterations took{" "}
          {(data.runtime / 1000000000).toFixed(3)} seconds to run).
          <br />
          git hash:{" "}
          {data.version ? (
            <a
              href={
                "https://github.com/genshinsim/gcsim/commits/" + data.version
              }
            >
              {data.version.substring(0, 8)}
            </a>
          ) : (
            "unknown"
          )}
          , built on: {data.build_date ? data.build_date : "unknown"}
        </div>
        <div className=" pl-4 pt-2 flex flex-row place-content-center">
          <div className="max-w-4xl w-full flex flex-col gap-1">
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
