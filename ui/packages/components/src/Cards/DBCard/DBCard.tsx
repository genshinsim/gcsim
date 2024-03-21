import { db, model } from "@gcsim/types"
import { AvatarCard } from "../AvatarCard/AvatarCard";
import { Card, CardContent, CardFooter } from "../../common/ui/card";
import { Table, TableBody, TableHead, TableHeader, TableRow } from "../../common/ui/table";
import { prettyPrintNumberStr } from "../../lib/helper";
import { Long } from "protobufjs";


const DBDetails = ({ summary, create_date }: { summary: db.IEntrySummary, create_date: number | Long | null }) => {
    let date = "unknown";
    if (create_date) {
        date = new Date((create_date as number) * 1000).toLocaleDateString();
    }
    return (
        <Table className="w-full grow">
            <TableHeader>
                <TableRow>
                    <TableHead className="priority-5 font-semibold">sim mode</TableHead>
                    <TableHead className="priority-5 font-semibold">target count</TableHead>
                    <TableHead className="priority-1 font-semibold">dps per target</TableHead>
                    <TableHead className="priority-1 font-semibold">avg sim time</TableHead>
                    <TableHead className="priority-3 font-semibold">create date</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                <TableRow>
                    <TableHead className="priority-5">{summary.mode ? "ttk" : "duration"}</TableHead>
                    <TableHead className="priority-5">{summary.target_count}</TableHead>
                    <TableHead className="priority-1">{prettyPrintNumberStr(summary.mean_dps_per_target?.toFixed(2) ?? "")}</TableHead>
                    <TableHead className="priority-1">
                        {summary.sim_duration?.mean
                            ? `${summary.sim_duration.mean.toPrecision(3)}s`
                            : "unknown"}
                    </TableHead>
                    <TableHead className="priority-3">{date}</TableHead>
                </TableRow>
            </TableBody>
        </Table>
    )

}

type DBCardProps = {
    entry: db.IEntry

    //optional send to simulator
    footer?: JSX.Element
}

export const DBCard = ({ entry, footer }: DBCardProps) => {
    const team: (model.ICharacter | null)[] = entry.summary?.team ?? [];
    if (team.length < 4) {
        const diff = 4 - team.length;
        for (let i = 0; i < diff; i++) {
            team.push(null);
        }
    }
    return (
        <Card className="m-2 bg-slate-800 min-[1300px]:w-[1225px] ">
            <CardContent className="p-3 flex flex-col  gap-y-2">
                <div className="flex flex-row flex-wrap gap-y-2 place-content-center">
                    <AvatarCard chars={team} className="min-[420px]:w-[420px] -ml-2" />
                    {
                        entry.summary ? <DBDetails summary={entry.summary} create_date={entry.create_date ? entry.create_date : null} /> : null
                    }

                </div>
                <div className="flex-grow flex flex-col  text-gray-400 font-san font-medium text-sm">
                    <div className="block w-0 min-w-full">
                        <span className="font-semibold text-orange-300">
                            {entry.submitter === "migrated"
                                ? "Unknown author: "
                                : `Submitted by ${entry.submitter}: `}
                        </span>
                        {entry.description}
                    </div>
                </div>
            </CardContent>
            <CardFooter className="flex flex-row flex-wrap gap-y-2 p-3 pt-0">
                {
                    footer ?
                        <div className="flex flex-row flex-wrap justify-end w-full">
                            {footer}
                        </div>
                        : null
                }
            </CardFooter>
        </Card>
    )


}