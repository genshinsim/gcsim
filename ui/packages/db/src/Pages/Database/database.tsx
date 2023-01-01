import { useCallback, useEffect, useState } from "react";
import { Filter } from "./Filter";
import { ListView } from "./ListView";
import { model } from "../../../protos_gen/protos";
import { fetchDataFromDB } from "../../api/FetchDataFromDB";
import DBEntryView from "./Components/DBEntryView";

const mockData: model.IDBEntries["data"] = [
  {
    team: [
      {
        name: "nahida",
        cons: 0,
        element: "dendro",
        level: 90,
        weapon: {
          level: 90,
          name: "favoniusgreatsword",
          maxLevel: 90,
          refine: 5,
        },
        sets: {
          gladiatorsfinale: 2,
          noblesseoblige: 2,
        },
        talents: {
          attack: 9,
          burst: 9,
          skill: 9,
        },
        maxLevel: 90,
      },
      {
        name: "dehya",
        cons: 0,
        element: "pyro",
        level: 90,
        weapon: {
          level: 90,
          name: "favoniuscodex",
          maxLevel: 90,
          refine: 5,
        },
        sets: {
          wandererstroupe: 2,
          gildeddreams: 2,
        },
        talents: {
          attack: 9,
          burst: 9,
          skill: 9,
        },
        maxLevel: 90,
      },
    ],
    dpsByTarget: {
      target1: {
        max: 100,
        mean: 100,
        min: 100,
        SD: 100,
      },
    },
    simDuration: {
      mean: 100,
      min: 100,
      max: 100,
      SD: 100,
    },
    config: "config",
    hash: "hash",
    //indexing fields
    charNames: ["dehya", "nahida"],
    targetCount: 1,
    meanDpsPerTarget: 100,
    createDate: 20210101,
    runDate: 20210101,
  },
];

export function Database() {
  const urlParams = window.location.search;
  const [data, setData] = useState<model.IDBEntries["data"]>([]);

  useEffect(() => {
    fetchDataFromDB(urlParams, setData);
    //mock data
    setData(mockData);
  }, [urlParams]);

  if (!data) {
    //TODO: add loading spinner or emoji
    return <div>Loading...</div>;
  }

  return (
    <>
      <Filter />
      {data[0] && <DBEntryView dbEntry={data[0]} />}
      <ListView
        query={{ char_names: "ayaka" }}
        sort="??"
        skip="??"
        limit="??"
      />
    </>
  );
}
