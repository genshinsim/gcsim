import React from "react";
import { fuse } from "../App";
import { AppContext } from "../Store";
import SearchBar from "./SearchBar";
import charData from "../data/character_images.json";
import CopyIcon from "../content_copy_white_24dp.svg";

export default function Search() {
  const { state } = React.useContext(AppContext);

  let data = fuse.search(state);
  console.log(data);

  const elements = data.map((r: any, i) => {
    return (
      <div
        className="m-2 p-2 rounded-md bg-gray-600 gap-1 items-center grid lg:grid-cols-4 md:grid-cols-2 sm:grid-cols-1"
        key={i}
      >
        <div className="flex flex-row">
          {r.item.characters.map((c: string) => {
            console.log(c);
            // @ts-ignore: Unreachable code error
            let image = charData[c];
            return (
              <div key={c} className="h-24">
                <img src={image} alt={c} className="object-contain h-full" />
              </div>
            );
          })}
        </div>
        <div className="flex-grow flex flex-row items-center lg:col-span-3 md:col-span-2">
          <div className="flex-grow">
            <div className="font-bold text-lg">{r.item.title}</div>
            <div>
              <strong>Author: </strong>
              {r.item.author}
            </div>
            <div>
              <strong>Version: </strong>
              {r.item.version}
            </div>
            <div>
              <strong>Description: </strong>
              {r.item.description}
            </div>
          </div>
          <div>
            <img
              src={CopyIcon}
              alt="copy"
              className="p-1 rounded-md hover:bg-gray-500"
              onClick={() => {
                navigator.clipboard.writeText(r.item.config);
              }}
            />
          </div>
        </div>
      </div>
    );
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
        {elements}
      </div>
    </div>
  );
}
