import { model } from "@gcsim/types";
import React, { useState } from "react";
import DBEntryView from "./Components/DBEntryView";

export interface ListViewProps {
  query?: any;
  sort?: any;
  skip?: any;
  limit?: any;
}

export function ListView(props: ListViewProps) {
  const [data, setData] = useState<model.IDBEntries["data"]>([]);

  React.useEffect(() => {
    const url = `/api/db?q=${encodeURIComponent(
      JSON.stringify(props.query)
    )}`;
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setData(data.data);
      })
      .catch((e) => {
        console.log(e);
      });
  }, []);
  if (!data) {
    //TODO: add loading spinner or emoji
    return <div>Loading...</div>;
  }
  return (
    <div className="flex flex-col w-full">
      {data.map((entry, index) => {
        return <DBEntryView dbEntry={entry} key={index} />;
      })}
    </div>
  );
}
