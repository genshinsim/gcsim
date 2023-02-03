import { model } from "@gcsim/types";
import DBEntryView from "./DBEntryView";

export function ListView({ data }: { data: model.IDBEntries["data"] }) {
  if (!data) {
    return <div>no data</div>;
  }

  return (
    <div className="flex flex-col gap-2">
      {data.map((entry, index) => {
        return <DBEntryView dbEntry={entry} key={index} />;
      })}
    </div>
  );
}
