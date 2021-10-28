import React from "react";
import { fuse } from "../App";
import { AppContext } from "../Store";
import SearchBar from "./SearchBar";
import TeamRow from "./TeamRow";
import charData from "../data/character_images.json";
import CopyIcon from "../content_copy_white_24dp.svg";

export default function SearchResults() {
  const { state } = React.useContext(AppContext);

  let data = fuse.search(state);
  console.log(data);

  const elements = data.map((r: any, i) => {
    return <TeamRow {...r.item} key={i} />;
  });
  //load fuzzy and search for str
  return (
    <div className="flex-grow flex p-10 flex-col">
      <div className="flex justify-end w-full ">
        <SearchBar />
      </div>
      <div className="flex flex-col">
        <div className="font-title p-4 text-2xl">
          {state === ""
            ? "Search for something first"
            : elements.length > 0
            ? "I found the following teams..."
            : "Sorry I didn't find anything :("}
        </div>
        <div>{elements}</div>
      </div>
    </div>
  );
}
