import nahida from "images/nahida.png";
import { FaPlus } from "react-icons/fa";
import { ArtifactSetFilter } from "./Filter.utils";

interface FilterPortraitProps {
  charName: string;
  sets?: ArtifactSetFilter;
  weapon?: string;
}

export function FilterPortrait(props: FilterPortraitProps) {
  return (
    <>
      <div className="hidden lg:flex">
        <FilterDesktopPortrait {...props} />
      </div>
      {/* <div className="lg:hidden block">
        <DBEntryMobilePortrait {...char} />
      </div> */}
    </>
  );
}

function FilterDesktopPortrait({
  charName: name,
  //   sets,
  weapon,
}: FilterPortraitProps) {
  if (!name) {
    return (
      <div className="bg-slate-700 p-2 w-20 flex flex-row  h-fit justify-center">
        <img src={nahida} className=" object-contain opacity-50 " />
      </div>
    );
  }
  return (
    <div className="bg-slate-700 p-2 flex flex-row max-h-fit    w-20">
      <div className="grid grid-cols-2 grid-rows-3 ">
        <div className="col-span-2 row-span-2 border-b border-white/25">
          <div className=" relative ">
            {name && (
              <img
                src={"https://gcsim.app/api/assets/avatar/" + name + ".png"}
                alt={name}
              />
            )}
          </div>
        </div>

        {/* <PortraitArtifactsComponent artifactSet={sets} /> */}
        <PortraitWeaponComponent weapon={weapon} />
      </div>
    </div>
  );
}

// function DBEntryMobilePortrait({ name, sets, weapon, cons }: model.ICharacter) {
//   if (!name) {
//     return (
//       <div className="bg-slate-700 p-2  max-h-20 flex flex-row justify-center">
//         <img src={nahida} className=" object-contain opacity-50" />
//       </div>
//     );
//   }
//   return (
//     <div className="bg-slate-700 p-2 flex flex-row max-h-fit    w-32">
//       <div className="grid grid-cols-3 grid-rows-2 bg-slate-400/50 gap-[1px]  ">
//         <div className="col-span-2 row-span-2 bg-slate-700">
//           <div className=" relative ">
//             {name && (
//               <img
//                 src={"https://gcsim.app/api/assets/avatar/" + name + ".png"}
//                 alt={name}
//               />
//             )}
//             <div className="absolute  right-0 bottom-0 text-xs font-bold opacity-80">
//               {(cons as number) ?? 0}
//             </div>
//           </div>
//         </div>

//         <PortraitArtifactsComponent artifactSet={sets} />
//         <PortraitWeaponComponent weapon={weapon} />
//       </div>
//     </div>
//   );
// }

function PortraitWeaponComponent({ weapon }: { weapon?: string }) {
  return (
    <button className={"bp4-button"}>
      {weapon ? (
        <img
          src={"https://gcsim.app/api/assets/weapons/" + weapon + ".png"}
          alt={weapon}
        />
      ) : (
        <FaPlus />
      )}
    </button>
  );
}

// function PortraitArtifactsComponent({
//   artifactSet,
// }: {
//   artifactSet:
//     | {
//         [k: string]: number | Long;
//       }
//     | undefined
//     | null;
// }) {
//   if (!artifactSet) {
//     return (
//       <div>
//         <img src={kuki} alt="kuki" className="relative max-h-full" />
//       </div>
//     );
//   }

//   const artifacts = Object.entries(artifactSet).filter(
//     ([, setCount]) => (setCount as number) >= 2
//   );

//   switch (artifacts.length) {
//     case 1:
//       return (
//         <div className="relative bg-slate-700">
//           <img
//             src={
//               "https://gcsim.app/api/assets/artifacts/" +
//               artifacts[0][0] +
//               "_flower.png"
//             }
//             alt={artifacts[0][0]}
//           />
//         </div>
//       );
//     case 2:
//       return (
//         <div className="relative">
//           <img
//             src={
//               "https://gcsim.app/api/assets/artifacts/" +
//               artifacts[0][0] +
//               "_flower.png"
//             }
//             alt={artifacts[0][0]}
//             className="artifact-top-right-gap"
//           />
//           <img
//             src={
//               "https://gcsim.app/api/assets/artifacts/" +
//               artifacts[1][0] +
//               "_flower.png"
//             }
//             alt={artifacts[1][0]}
//             className="artifact-bottom-left-gap"
//           />
//         </div>
//       );

//     default:
//     case 0:
//       return (
//         <div>
//           <img src={kuki} alt="kuki" className="relative opacity-30" />
//         </div>
//       );
//   }
// }
