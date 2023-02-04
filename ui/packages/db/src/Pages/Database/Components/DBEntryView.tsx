import { model } from "@gcsim/types";
import { Long } from "protobufjs";
import { useState } from "react";

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
    <div className="flex flex-row bg-slate-800  p-4 gap-4 w-full">
      <div className="flex gap-2 flex-row min-w-fit ">
        {team &&
          team.map((char, index) => {
            return (
              <DBEntryCharacterPortrait {...char} key={index.toString()} />
            );
          })}
      </div>
      <div className="flex flex-col grow ">
        {/* <div className="capitalize text-lg font-semibold ">
          {dbEntry?.char_names?.toString().replaceAll(",", ", ")}
        </div> */}
        <div>
          <div className="flex flex-col max-w-4xl">
            <DBEntryTags tags={dbEntry.tags} />
            <span className="  overflow-hidden">{dbEntry?.description}</span>
          </div>

          <DBEntryDetails {...dbEntry} />
        </div>
      </div>
      <DBEntryActions simulation_key={dbEntry.key} />
    </div>
  );
}

function DBEntryCharacterPortrait({
  name,
  sets,
  weapon,
  cons,
}: model.ICharacter) {
  if (!name) {
    return <div className="bg-slate-700 p-2 flex flex-row w-20 h-20"></div>;
  }
  return (
    <div className="bg-slate-700 p-2 flex flex-row w-20 h-20">
      <div className="flex flex-col  min-h-min">
        <div className=" relative">
          {name && (
            <img
              src={"https://gcsim.app/api/assets/avatar/" + name + ".png"}
              alt={name}
            />
          )}
          <div className="absolute  right-0 bottom-0 text-xs font-bold">
            {(cons as number) ?? 0}
          </div>
        </div>
        <PortraitWeaponComponent weapon={weapon} />
      </div>

      <PortraitArtifactsComponent artifactSet={sets} />
    </div>
  );
}
function PortraitWeaponComponent({
  weapon,
}: {
  weapon: model.IWeapon | undefined | null;
}) {
  if (!weapon || !weapon.name) {
    return <div className="h-16 w-16">No weapon</div>;
  }
  return (
    <div className=" relative   ">
      <img
        src={"https://gcsim.app/api/assets/weapons/" + weapon.name + ".png"}
        alt={weapon.name}
      />
      <div className=" absolute bottom-0 right-0  text-xs  font-semibold">
        R{weapon?.refine?.toString()}
      </div>
    </div>
  );
}
function PortraitArtifactsComponent({
  artifactSet,
}: {
  artifactSet:
    | {
        [k: string]: number | Long;
      }
    | undefined
    | null;
}) {
  if (!artifactSet) {
    return <div className="h-4 w-4">No artifact</div>;
  }
  return (
    <div className="flex flex-col   ">
      {Object.entries(artifactSet).map(
        ([setName, setCount]) =>
          (setCount as number) > 0 && (
            <div className=" relative  ">
              <img
                src={
                  "https://gcsim.app/api/assets/artifacts/" +
                  setName +
                  "_flower.png"
                }
                alt={setName}
              />
              <div className=" absolute bottom-0 right-0  text-xs  font-semibold">
                {setCount.toString()}
              </div>
            </div>
          )
      )}
    </div>
  );
}

function DBEntryTags({ tags }: { tags: string[] | undefined | null }) {
  const t = useTranslation();
  const [showAll, setShowAll] = useState(false);
  return (
    <div
      className={
        "flex flex-row h-full flex-wrap  relative  max-w-md " +
        (showAll ? " " : " truncate")
      }
    >
      {tags?.map((tag) => (
        <div
          className="bg-slate-700 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
          key={tag}
        >
          {t(tag)}
        </div>
      ))}
      <button
        className=" absolute right-0 top-1  text-xs font-semibold     bg-blue-400/30 p-1 mr-1   whitespace-nowrap rounded-sm h-fit   "
        onClick={() => setShowAll(!showAll)}
      >
        {showAll ? "▲" : "▼"}
      </button>
    </div>
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
        <tr className=" ">
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

function DBEntryActions({
  simulation_key,
}: {
  simulation_key: string | undefined | null;
}) {
  return (
    <div className="flex flex-col justify-center">
      <a
        href={`https://gcsim.app/v3/viewer/share/${simulation_key}`}
        target="_blank"
        className="bp4-button    bp4-intent-primary"
        rel="noreferrer"
      >
        Open in Viewer
      </a>
    </div>
  );
}
