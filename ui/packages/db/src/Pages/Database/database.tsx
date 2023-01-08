import { useEffect, useState } from "react";
import { model } from "@gcsim/types";

import { Filter } from "./Filter";
import { ListView } from "./ListView";
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
        name: "xingqiu",
        cons: 0,
        element: "hydro",
        level: 90,
        weapon: {
          level: 90,
          name: "favoniuscodex",
          maxLevel: 90,
          refine: 5,
        },
        sets: {
          wandererstroupe: 4,
        },
        talents: {
          attack: 9,
          burst: 9,
          skill: 9,
        },
        maxLevel: 90,
      },
      {
        name: "yelan",
        cons: 0,
        element: "hydro",
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
      {
        name: "raiden",
        cons: 0,
        element: "cryo",
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
    charNames: ["nahida", "xingqiu", "yelan", "raiden"],
    targetCount: 1,
    meanDpsPerTarget: 100,
    createDate: 20210101,
    runDate: 20210101,
    description:
      "  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean ut augue dapibus, interdum ante quis, congue nisi. Nulla ut sagittis lorem. Aliquam lobortis, urna sit amet fringilla porttitor, diam metus laoreet odio, rutrum maximus ligula ipsum ac nunc. Vivamus euismod nec neque at pharetra. Curabitur aliquam lectus diam. Suspendisse quis ultrices odio. Aenean eleifend condimentum nibh auctor ultrices. Curabitur sit amet erat vitae odio imperdiet consectetur. Suspendisse at massa bibendum, rutrum sem eu, laoreet magna. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Proin a maximus quam. Duis tempor viverra dui, at lacinia ligula Aliquam erat volutpat. Etiam pulvinar, elit a blandit sodales, erat libero bibendum turpis, quis sollicitudin felis ipsum ac velit. Nunc lobortis est eget lacus gravida consectetur. Morbi sed venenatis odio. Quisque vel finibus risus. Duis ut cursus magna. Donec egestas ante vitae neque finibus, at blandit ex porta. Proin molestie orci vitae velit ornare facilisis. Morbi auctor sodales maximus. Etiam eu posuere augue. Nunc vehicula ut est ac placerat. Pellentesque dignissim vitae lectus ac faucibus. Aliquam accumsan mi ut magna rutrum, commodo dapibus augue malesuada. Aenean vitae erat nec elit pharetra semper id sit amet orci. Nam ullamcorper euismod elit, nec elementum sem porta eu. Nam et velit dictum, varius odio sit amet, venenatis eros. ",
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
