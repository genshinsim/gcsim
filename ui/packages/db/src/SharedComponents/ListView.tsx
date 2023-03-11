import { Spinner } from "@blueprintjs/core";
import { model } from "@gcsim/types";
import DBEntryView from "./DBEntryView";

export function ListView({ data }: { data: model.IDBEntries["data"] }) {
  if (!data) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  return (
    <>
      <div className="flex flex-col gap-2 justify-center align-middle items-center ">
        {data.map((entry, index) => {
          return <DBEntryView dbEntry={entry} key={index} />;
        })}
      </div>
      {/* <div className="flex flex-col gap-2">
        {data.map((entry, index) => {
          return <DBEntryView dbEntry={entry} key={index} />;
        })}
      </div> */}
    </>
  );
}
