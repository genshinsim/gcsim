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
} from "../../SharedComponents/FilterComponents/Filter.utils";
import { DBView } from "./DBVIew";

export function Database() {
  const [filter, dispatch] = useReducer(filterReducer, initialFilter);
  const [data, setData] = useState<db.IEntry[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [hasMore, setHasMore] = useState<boolean>(true)
  const [page, setPage] = useState<number>(1)

  const appendData = (next: db.IEntry[]) => {
    // let d = [ ...data,...next.filter(e => {
    //   return false
    // })]
    setData([...data, ...next])
  }

  const querydb = (query: DbQuery, nextPage: number, append: boolean) => {
    axios(`/api/db?q=${encodeURIComponent(JSON.stringify(query))}`)
      .then((resp: { data: db.IEntries }) => {
        if (resp.data && resp.data.data) {
          setPage(nextPage)
          setHasMore(true)
          if (append) {
            appendData(resp.data.data)
          } else {
            setData(resp.data.data)
          }
          //check count; if we got less than limit then there's no more data... 
          //TODO: this is bugged if there are exactly limit number of entries...
          //TODO: really server should tell us if there's more data
          if (resp.data.data.length < query.limit) {
            setHasMore(false)
          }
          console.log("data: ", resp.data.data);
        } else {
          setHasMore(false)
          if (!append) {
            setData([]);
          }
        }
        setIsLoading(false);
      })
      .catch((err) => {
        console.log("error: ", err);
      });
  };

  useEffect(() => {
    const query = craftQuery(filter, 1, 25);
    querydb(query, 1, false);
  }, [filter]);

  const fetchData = () => {
    const nextPage = page + 1
    const query = craftQuery(filter, nextPage, 25);
    querydb(query, nextPage, true)
  }

  if (isLoading || !data)
    return (
      <div className="h-screen flex flex-col justify-center items-center">
        <Spinner />
      </div>
    );

  return (
    <FilterContext.Provider value={filter}>
      <FilterDispatchContext.Provider value={dispatch}>
        <DBView
          data={data}
          fetchData={fetchData}
          hasMore={hasMore}
        />
      </FilterDispatchContext.Provider>
    </FilterContext.Provider>
  );
}
