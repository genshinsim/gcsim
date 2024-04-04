import { db, model } from "@gcsim/types";
import { Long } from "protobufjs";
import { ReactI18NextChild, useTranslation } from "react-i18next";
import { DBEntryPortrait } from "./DBEntryViewComponents/DBEntryPortrait";
import DBEntryTags from "./DBEntryViewComponents/DBEntryTags";

//displays one database entry
export default function DBEntryView({ dbEntry }: { dbEntry: db.Entry }) {
  const { t: translate } = useTranslation();
  const t = (key: string) => translate(key) as ReactI18NextChild; // idk why this is needed

  const team = dbEntry.summary?.team ?? [];
  if (team.length < 4) {
    const diff = 4 - team.length;
    for (let i = 0; i < diff; i++) {
      team.push({} as model.Character);
    }
  }
  let link = `https://gcsim.app/sh/${dbEntry.share_key}`;
  if ("_id" in dbEntry) {
    link = `https://gcsim.app/db/${dbEntry["_id"]}`;
  }

  return (
    <>
      <div className="flex flex-row flex-wrap place-content-center bg-slate-700 lg:bg-slate-800 max-w-xs sm:min-w-wsm md:min-w-wmd lg:min-w-wlg xl:min-w-wxl sm:max-w-sm md:max-w-2xl lg:max-w-4xl p-5 border sm:border-0 gap-4 sm:gap-1 ">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
          {team &&
            team.map((char, index) => {
              return <DBEntryPortrait {...char} key={index.toString()} />;
            })}
        </div>
        <div className="flex flex-col grow h-full">
          {visibleTagCount(dbEntry.accepted_tags ?? []) > 0 ? (
            <DBEntryTags tags={dbEntry.accepted_tags} />
          ) : (
            <></>
          )}
          <DBEntryDetails
            {...dbEntry.summary}
            create_date={dbEntry.create_date}
            description={dbEntry.description}
          />
          <div className="hidden p-1 max-h-7 opacity-50 lg:block lg:w-0 lg:min-w-full">
            {dbEntry.description}
          </div>
        </div>
        <div className="flex flex-col justify-center w-full lg:w-fit">
          <a
            href={link}
            target="_blank"
            className="bp4-button bp4-intent-primary w-full"
            rel="noreferrer"
          >
            <div className="m-0">{t("db.openInViewer")}</div>
          </a>
        </div>
        <div className="basis-full text-xs font-bold w-full flex place-content-start">
          {dbEntry.submitter === "migrated"
            ? "Unknown author"
            : `${translate("db.author")} ${dbEntry.submitter}`}
        </div>
      </div>
    </>
  );
}

function DBEntryDetails({
  target_count,
  mean_dps_per_target,
  mode,
  sim_duration,
  create_date,
}: NonNullable<db.Entry["summary"]> & {
  create_date?: number | Long | null;
  description?: string | null;
}) {
  const { t: translate } = useTranslation();

  const t = (key: string) => translate(key) as ReactI18NextChild; // idk why this is needed
  let date = t("db.unknown");
  if (create_date) {
    date = new Date((create_date as number) * 1000).toLocaleDateString();
  }
  return (
    <table className="bp4-html-table w-full">
      <thead>
        <tr className="text-xs">
          <th className="priority-5">
            <div>{t("db.simMode") as ReactI18NextChild}</div>
          </th>
          <th className="priority-5">{t("db.targetCount")}</th>
          <th className="priority-1">{t("db.dpsPerTarget")}</th>
          <th className="priority-1">{t("db.avgSimTime")}</th>
          <th className="priority-3">{t("db.createDate")}</th>
        </tr>
      </thead>
      <tbody>
        <tr className=" text-xs ">
          <td className="priority-5">
            {mode ? t("db.ttk") : t("db.duration")}
          </td>
          <td className="priority-5">{target_count}</td>
          <td className="priority-1">
            {prettyPrintNumberStr(mean_dps_per_target?.toFixed(2) ?? "")}
          </td>
          <td className="priority-1">
            {sim_duration?.mean
              ? `${sim_duration.mean.toPrecision(3)}s`
              : t("db.unknown")}
          </td>
          <td className="priority-3">{date}</td>
        </tr>
      </tbody>
    </table>
  );
}

function prettyPrintNumberStr(num: string): string {
  return num.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function visibleTagCount(tags: model.DBTag[]): number {
  // 1 is the gcsim tag
  return tags.filter((tag) => tag !== 1).length;
}
