import React from "react";
import { Spinner, Viewport } from "Components";
import { useAppDispatch, useAppSelector } from "Store";
import { loadDB } from "Store/dbSlice";
import IngameNamesJson from "./IngameNamesJson.json";
import { CharsGrid } from "./CharsGrid";

const charNames: { [key in string]: string } =
  IngameNamesJson.English.character_names;

export function Database() {
  const status = useAppSelector((state) => state.db.status);
  const errorMsg = useAppSelector((state) => state.db.errorMsg);
  const chars = useAppSelector((state) => state.db.characters);
  const [searchString, setSearchString] = React.useState<string>("");
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    if (status === "idle") {
      dispatch(loadDB());
    }
  }, [status, dispatch]);

  switch (status) {
    case "loading":
    case "idle":
      return (
        <Viewport>
          <div className="flex flex-row place-content-center mt-2">
            <Spinner />
          </div>
        </Viewport>
      );
    case "error":
      return (
        <Viewport>
          <div className="flex flex-row place-content-center mt-2">
            {errorMsg}
          </div>
        </Viewport>
      );
  }

  var avail = chars.map((c) => c.avatar_name);
  avail.sort();
  var charsEntries: string[][] = [];
  avail.forEach((c) => {
    if (c in charNames) {
      charsEntries.push([c, charNames[c]]);
    }
  });

  const filteredChars = searchString
    ? charsEntries.filter(([, longName]) =>
        longName.toLocaleLowerCase().includes(searchString)
      )
    : charsEntries;

  return (
    <div className="flex flex-row justify-center">
      <Viewport>
        <div className=" flex flex-row gap-x-1 justify-end">
          <div>
            <div className="relative mt-1 rounded-md">
              <input
                type="text"
                className="block w-full rounded-md pl-4 pr-4 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm text-white bg-slate-700 h-8"
                placeholder="filter by character"
                value={searchString}
                onChange={(e) => setSearchString(e.currentTarget.value)}
              />
            </div>
          </div>
        </div>
        <div className="border-b-2 mt-2 border-gray-300" />
        <CharsGrid characters={filteredChars} />
      </Viewport>
    </div>
  );
}
