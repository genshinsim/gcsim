import { Spinner } from "@blueprintjs/core";
import { db } from "@gcsim/types";
import axios from "axios";
import eula from "images/eula.png";
import { useEffect, useReducer, useState } from "react";
import { ActionBar } from "SharedComponents/ActionBar";
import { PaginationButtons } from "SharedComponents/Pagination";
import { craftQuery, DbQuery } from "SharedHooks/databaseQuery";
import {
  FilterContext,
  FilterDispatchContext,
  filterReducer,
  initialFilter,
} from "../SharedComponents/FilterComponents/Filter.utils";
import { ListView } from "../SharedComponents/ListView";

export function Database() {
  const [filter, dispatch] = useReducer(filterReducer, initialFilter);
  const [data, setData] = useState<db.IEntry[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  const querydb = (query: DbQuery) => {
    axios(`/api/db?q=${encodeURIComponent(JSON.stringify(query))}`)
      .then((resp: { data: db.IEntries }) => {
        if (resp.data && resp.data.data) {
          setData(resp.data.data);
          console.log("data: ", resp.data.data);
        } else {
          setData([]);
        }
        setIsLoading(false);
      })
      .catch((err) => {
        console.log("error: ", err);
      });
  };

  useEffect(() => {
    const query = craftQuery(filter);
    querydb(query);
  }, [filter]);

  if (isLoading || !data)
    return (
      <div className="h-screen flex flex-col justify-center items-center">
        <Spinner />
      </div>
    );

  return (
    <FilterContext.Provider value={filter}>
      <FilterDispatchContext.Provider value={dispatch}>
        <div className="flex flex-col  gap-4 m-8 my-4 items-center">
          <ActionBar simCount={data.length} />
          {data.length === 0 ? (
            <div className="6 flex flex-col justify-center items-center h-screen">
              <img
                src={eula}
                className=" object-contain opacity-50 w-32 h-32"
              />
            </div>
          ) : (
            <ListView data={data} />
          )}
          <PaginationButtons />
        </div>
      </FilterDispatchContext.Provider>
    </FilterContext.Provider>
  );
}
