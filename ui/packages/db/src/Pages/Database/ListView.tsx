import { model } from "@gcsim/types";
import DBEntryView from "./Components/DBEntryView";

export function ListView({ data }: { data: model.IDBEntries["data"] }) {
  if (!data) {
    return <div>no data</div>;
  }

  return (
    <div className="flex flex-col">
      {data.map((entry, index) => {
        return <DBEntryView dbEntry={entry} key={index} />;
      })}
    </div>
  );
}
