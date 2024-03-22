import tagData from "@gcsim/data/src/tags.json";
import { db, model } from "@gcsim/types";
import { Card, CardContent, CardFooter } from "../../common/ui/card";
import { AvatarCard } from "../AvatarCard/AvatarCard";
import { CardBadge } from "../CardBadge/CardBadge";

type DBCardProps = {
  entry: db.IEntry;

  //optional send to simulator
  footer?: JSX.Element;
};

export const DBCard = ({ entry, footer }: DBCardProps) => {
  const team: (model.ICharacter | null)[] = entry.summary?.team ?? [];
  if (team.length < 4) {
    const diff = 4 - team.length;
    for (let i = 0; i < diff; i++) {
      team.push(null);
    }
  }
  let date = "unknown";
  if (entry.create_date) {
    date = new Date((entry.create_date as number) * 1000).toLocaleDateString();
  }
  const tags = entry.accepted_tags
    ?.filter((tag) => tag !== 1)
    .map((tag) => (
      <CardBadge
        key={tag}
        value={tagData[tag]?.display_name ?? null}
        valueCase=""
        valueSize="text-xs"
        className="bg-teal-800"
      />
    ));

  return (
    <Card className="m-2 bg-slate-800 min-[1300px]:w-[1225px] ">
      <CardContent className="p-3 flex flex-col gap-y-2">
        <div className="flex flex-row flex-wrap gap-2 place-content-center">
          <Card className="flex flex-col bg-slate-800 border-0 pt-1 min-[420px]:basis-0">
            <AvatarCard chars={team} className="min-[420px]:w-[420px] flex-1" />
            <div className="flex flex-row flex-wrap gap-1 p-2 max-w-[95%] place-content-evenly">
              <CardBadge
                title="mode"
                value={entry.summary?.mode ? "ttk" : "duration"}
              />
              <CardBadge
                title="target count"
                value={(entry.summary?.target_count ?? 0).toString()}
              />
              <CardBadge
                title="dps/target"
                value={(entry.summary?.mean_dps_per_target ?? 0).toLocaleString(
                  navigator.language,
                  {
                    notation: "compact",
                    minimumSignificantDigits: 3,
                    maximumSignificantDigits: 3,
                  }
                )}
              />
              <CardBadge
                valueCase="lowercase"
                title="avg sim time"
                value={
                  entry.summary?.sim_duration?.mean
                    ? `${entry.summary?.sim_duration.mean.toPrecision(3)}s`
                    : "unknown"
                }
              />
              <CardBadge valueCase="lowercase" title="created" value={date} />
              {tags}
            </div>
          </Card>
          <div className="flex flex-col grow min-w-[40%] text-gray-200 p-2 self-stretch">
            <div className="block w-0 min-w-full">
              <span className="font-semibold text-orange-300">
                {entry.submitter === "migrated"
                  ? "Unknown author: "
                  : `Submitted by ${entry.submitter}: `}
              </span>
              {entry.description}
            </div>
            {footer ? (
              <div className="mt-auto flex flex-row flex-wrap justify-end w-full pt-2">
                {footer}
              </div>
            ) : null}
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex flex-row flex-wrap gap-y-2 p-3 pt-0"></CardFooter>
    </Card>
  );
};