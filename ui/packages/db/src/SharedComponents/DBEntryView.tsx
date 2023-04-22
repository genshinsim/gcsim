import { model } from "@gcsim/types";
import DBEntryActions from "./DBEntryViewComponents/DBEntryActions";
import { DBEntryPortrait } from "./DBEntryViewComponents/DBEntryPortrait";
import DBEntryTags from "./DBEntryViewComponents/DBEntryTags";

function useTranslation() {
  return (text: string) => text;
}

//displays one database entry
export default function DBEntryView({ dbEntry }: { dbEntry: model.IDBEntry }) {
  // const t = useTranslation();
  const team = dbEntry.team ?? [];
  if (team.length < 4) {
    const diff = 4 - team.length;
    for (let i = 0; i < diff; i++) {
      team.push({} as model.ICharacter);
    }
  }
  return (
    <>
      <div className="hidden lg:flex  flex-row bg-slate-800  p-4 gap-4 w-full max-w-7xl">
        <div className="flex gap-2 flex-row min-w-fit ">
          {team &&
            team.map((char, index) => {
              return <DBEntryPortrait {...char} key={index.toString()} />;
            })}
        </div>
        <div className="flex flex-col grow ">
          <div className="max-w-2xl">
            <div className="flex flex-col ">
              <DBEntryTags tags={dbEntry.accepted_tags} />
              <span className="  overflow-hidden">{dbEntry?.description}</span>
            </div>

            <DBEntryDetails {...dbEntry} />
          </div>
        </div>
        <DBEntryActions simulation_key={dbEntry.id} id={dbEntry.id} />
      </div>
      <div className="lg:hidden flex flex-row bg-slate-700 max-w-xs pr-2 gap-2  ">
        <div className="grid grid-cols-2 grid-row-2  ">
          {team &&
            team.map((char, index) => {
              return <DBEntryPortrait {...char} key={index.toString()} />;
            })}
        </div>
        <DBEntryActions simulation_key={dbEntry.id} id={dbEntry.id} />
      </div>
    </>
  );
}

function DBEntryDetails({
  target_count,
  mean_dps_per_target,
  run_date,
  mode,
  sim_duration,
  total_damage,
}: // total_damage,
// description,
model.IDBEntry) {
  const t = useTranslation();
  return (
    <table className="bp4-html-table  ">
      <thead>
        <tr className="">
          <th className="">{t("Sim Mode")}</th>
          <th className="">{t("Target Count")}</th>
          <th className="">{t("DPS per Target")}</th>
          <th className="">{t("AVG Total Damage")}</th>
          <th className="">{t("AVG Sim Time")}</th>
          <th className="">{t("Run Date")}</th>
        </tr>
      </thead>
      <tbody>
        <tr className=" text-xs ">
          <td className="">{mode ? t("Time to kill") : t(" Duration")}</td>
          <td className="">{target_count}</td>
          <td className="">{mean_dps_per_target?.toPrecision(8)}</td>
          <td className="">{total_damage?.mean?.toPrecision(8)}</td>
          <td className="">
            {sim_duration?.mean?.toPrecision(3) ?? t("Unknown")}s
          </td>
          <td className="">{run_date?.toString() ?? t("Unknown")}</td>
        </tr>
      </tbody>
    </table>
  );
}
