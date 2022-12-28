export function Database() {
  const data = await fetch("https://simimpact.app/api/db");
  //   data
  //     .then((res) => {
  //       idk = res.body;
  //     })
  //     .catch((e) => {
  //       console.log(e);
  //     });

  return (
    <>
      <div className="flex justify-center my-40">
        <div className="">
          {/* {dbData.map((sim) => {
            return (
              <div key={sim.hash}>
                <div className="flex flex-row">
                  <div className="flex flex-col"></div>
                </div>
              </div>
            );
          })} */}
        </div>
      </div>
    </>
  );
}

async function fetchDatabase() {
  const data = await fetch("https://simimpact.app/api/db");
  return data;
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
