import { useCallback, useEffect, useState } from "react";
import { Filter } from "./Filter";
import { ListView } from "./ListView";
import { model } from "../../../protos_gen/protos";
import axios from "axios";
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
    charNames: ["dehya, nahida"],
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

// export interface sim {
//   team: Character[];
//   dps_by_target: {
//     dps: SummaryStats;
//     //room for other stats in future
//   };
//   iter: number;
//   sim_duration: SummaryStats;
//   config: string;
//   hash: string;
//   //indexing fields
//   char_names: string[];
//   target_count: number;
//   mean_dps_per_target: number;
//   create_date: string;
//   run_date: string;
// }

// type Character = Record<string, unknown>;
// type SummaryStats = Record<string, unknown>;

// const mockData: sim[] = [
//   {
//     team: [
//       {
//         id: "1000056",
//         name: "Nahida",
//       },
//       {
//         id: "1000069",
//         name: "Dehya",
//       },
//     ],
//     dps_by_target: {
//       dps: {},
//     },
//     iter: 1000,
//     sim_duration: {},
//     config: "config",
//     hash: "hash",
//     //indexing fields
//     char_names: ["Dehya, Nahida"],
//     target_count: 1,
//     mean_dps_per_target: 100,
//     create_date: "2021-01-01",
//     run_date: "2021-01-01",
//   },
// ];

// function DisplaySim({ sim }: { sim: sim }) {
//   const legendaryCss = " bg-[#FFB13F] ";
//   const rareCss = " bg-[#D28FD6]";
//   return (
//     <>
//       {sim.char_names.map((charName, index) => {
//         return (
//           <div
//             key={index}
//             className={
//               "rounded-md bg-opacity-70" +
//               (rareCharNames.includes(charName) ? rareCss : legendaryCss)
//             }
//           >
//             <img
//               src={"/api/assets/avatar/" + charName + ".png"}
//               alt={charName}
//               className="margin-auto"
//             />
//             <div className="text-xs flex items-center justify-center text-center h-8 bg-slate-600 text-white">
//               {charName}
//             </div>
//           </div>
//         );
//       })}
//     </>
//   );
// }

// const rareCharNames = [
//   "amber",
//   "barbara",
//   "beidou",
//   "bennett",
//   "candace",
//   "chongyun",
//   "collei",
//   "diona",
//   "dori",
//   "fischl",
//   "gorou",
//   "kaeya",
//   "lisa",
//   "kuki",
//   "ningguang",
//   "noelle",
//   "razor",
//   "heizou",
//   "rosaria",
//   "sara",
//   "sucrose",
//   "sayu",
//   "thoma",
//   "xiangling",
//   "xinyan",
//   "xingqiu",
//   "yanfei",
//   "yunjin",
// ];
