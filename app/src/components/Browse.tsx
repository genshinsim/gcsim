import React from "react";
import data from "../data/configs.json";
import TeamRow from "./TeamRow";

export default function Browse() {
  const elements = data.map((r, i) => {
    return <TeamRow {...r} key={i} />;
  });
  //load fuzzy and search for str
  return (
    <div className="flex-grow flex p-10 flex-col overflow-y-scroll overflow-x-hidden">
      <div className="flex flex-col">{elements}</div>
    </div>
  );
}
