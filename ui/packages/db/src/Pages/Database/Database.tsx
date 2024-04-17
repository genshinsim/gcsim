import { Spinner } from "@blueprintjs/core";
import { db } from "@gcsim/types";
import axios from "axios";
import { useEffect, useReducer, useState } from "react";
import { craftQuery, DbQuery } from "SharedHooks/databaseQuery";
import {
  FilterContext,
  FilterDispatchContext,
  filterReducer,
  FilterState,
  initialFilter as defaultFilter,
} from "../../SharedComponents/FilterComponents/Filter.utils";
import { DBView } from "./DBVIew";

type Props = {
  initialFilter?: FilterState;
};

export const Database = ({ initialFilter = defaultFilter }: Props) => {
  const [filter, dispatch] = useReducer(filterReducer, initialFilter);
  const [data, setData] = useState<db.Entry[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [page, setPage] = useState<number>(1);

  const appendData = (next: db.Entry[]) => {
    // let d = [ ...data,...next.filter(e => {
    //   return false
    // })]
    setData([...data, ...next]);
  };

  const querydb = (query: DbQuery, nextPage: number, append: boolean) => {
    axios(`/api/db?q=${encodeURIComponent(JSON.stringify(query))}`)
      .then((resp: { data: db.Entries }) => {
        if (resp.data && resp.data.data) {
          setPage(nextPage);
          setHasMore(true);
          if (append) {
            appendData(resp.data.data);
          } else {
            setData(resp.data.data);
          }
          //check count; if we got less than limit then there's no more data...
          //TODO: this is bugged if there are exactly limit number of entries...
          //TODO: really server should tell us if there's more data
          if (resp.data.data.length < query.limit) {
            setHasMore(false);
          }
        } else {
          setHasMore(false);
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
    const nextPage = page + 1;
    const query = craftQuery(filter, nextPage, 25);
    querydb(query, nextPage, true);
  };

  if (isLoading || !data)
    return (
      <div className="h-screen flex flex-col justify-center items-center">
        <Spinner />
      </div>
    );

  return (
    <FilterContext.Provider value={filter}>
      <FilterDispatchContext.Provider value={dispatch}>
        <DBView data={data} fetchData={fetchData} hasMore={hasMore} />
      </FilterDispatchContext.Provider>
    </FilterContext.Provider>
  );
};
