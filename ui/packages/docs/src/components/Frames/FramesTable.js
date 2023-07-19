import React from "react";
import ReactPlayer from "react-player";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import character_data from "./character_data.json";

function Vid({ vid }) {
  return (
    <div>
      <ReactPlayer controls url={vid.vid} />
      Video credit: {vid.vid_credit === "" ? "Unknown" : vid.vid_credit}
      <br />
      Count :{" "}
      {vid.count !== "" ? (
        <>
          <a href={vid.count} target="_blank" rel="noreferrer">
            Sheet
          </a>{" "}
          (credit: {vid.count_credit === "" ? "Unknown" : vid.count_credit})
        </>
      ) : (
        <span>None available</span>
      )}
    </div>
  );
}

function MultiVids({ vids }) {
  const tabs = vids.map((e, i) => {
    return (
      <TabItem value={i} label={`Video #${i + 1}`} key={i}>
        <Vid vid={e} />
      </TabItem>
    );
  });
  return <Tabs>{tabs}</Tabs>;
}

export default function FramesTable({ item_key }) {
  if (!(item_key in character_data)) {
    return <div>Character does not have any frames video</div>;
  }
  const char = character_data[item_key];
  if (char.length === 0) {
    return <div>Character does not have any frames video</div>;
  }
  if (char.length === 1) {
    return <Vid vid={char[0]} />;
  }
  return <MultiVids vids={char} />;
}
