import { model } from "@gcsim/types";
import kuki from "images/kuki.png";
import nahida from "images/nahida.png";
import { Long } from "protobufjs";

export function DBEntryPortrait(char: model.Character) {
  return (
    <>
      <div className="hidden lg:flex">
        <DBEntryDesktopPortrait {...char} />
      </div>
      <div className="lg:hidden block">
        <DBEntryMobilePortrait {...char} />
      </div>
    </>
  );
}

function DBEntryDesktopPortrait({ name, sets, weapon, cons }: model.Character) {
  if (!name) {
    return (
      <div className="bg-slate-700 p-2 w-20 flex flex-row  h-fit justify-center">
        <img src={nahida} className=" object-contain opacity-50 " />
      </div>
    );
  }
  return (
    <div className="bg-slate-700 p-2 flex flex-row h-32  w-20">
      <div className="grid grid-cols-2 grid-rows-3 ">
        <div className="col-span-2 row-span-2 border-b border-white/25">
          <div className=" relative ">
            {name && (
              <img
                src={"https://gcsim.app/api/assets/avatar/" + name + ".png"}
                alt={name}
              />
            )}
            <div className="absolute  right-0 bottom-0 text-xs font-bold opacity-80">
              {(cons as number) ?? 0}
            </div>
          </div>
        </div>

        <PortraitArtifactsComponent artifactSet={sets} />
        <PortraitWeaponComponent weapon={weapon} />
      </div>
    </div>
  );
}

function DBEntryMobilePortrait({ name, sets, weapon, cons }: model.Character) {
  if (!name) {
    return (
      <div className="bg-slate-700 p-2  max-h-20 flex flex-row justify-center">
        <img src={nahida} className=" object-contain opacity-50" />
      </div>
    );
  }
  return (
    <div className="bg-slate-700 flex flex-row max-h-fit max-w-[128px]">
      <div className="grid grid-cols-3 grid-rows-2 bg-slate-500/10 gap-[2px]  ">
        <div className="col-span-2 row-span-2 bg-slate-700">
          <div className=" relative ">
            {name && (
              <img
                src={"https://gcsim.app/api/assets/avatar/" + name + ".png"}
                alt={name}
              />
            )}
            <div className="absolute  right-0 bottom-0 text-xs font-bold opacity-80">
              {(cons as number) ?? 0}
            </div>
          </div>
        </div>

        <PortraitArtifactsComponent artifactSet={sets} />
        <PortraitWeaponComponent weapon={weapon} />
      </div>
    </div>
  );
}

function PortraitWeaponComponent({
  weapon,
}: {
  weapon: model.Weapon | undefined | null;
}) {
  if (!weapon || !weapon.name) {
    return <div className="h-16 w-16">?</div>;
  }
  return (
    <div className="bg-slate-700 relative">
      <img
        src={"https://gcsim.app/api/assets/weapons/" + weapon.name + ".png"}
        alt={weapon.name}
      />
      <div className=" absolute bottom-0 right-0  text-xs  font-semibold opacity-80">
        R{weapon?.refine?.toString()}
      </div>
    </div>
  );
}

function PortraitArtifactsComponent({
  artifactSet,
}: {
  artifactSet:
    | {
        [k: string]: number | Long;
      }
    | undefined
    | null;
}) {
  if (!artifactSet) {
    return (
      <div>
        <img src={kuki} alt="kuki" className="relative max-h-full" />
      </div>
    );
  }

  const artifacts = Object.entries(artifactSet).filter(
    ([, setCount]) => (setCount as number) >= 2
  );

  switch (artifacts.length) {
    case 1:
      return (
        <div className="relative bg-slate-700">
          <img
            src={
              "https://gcsim.app/api/assets/artifacts/" +
              artifacts[0][0] +
              "_flower.png"
            }
            alt={artifacts[0][0]}
          />
        </div>
      );
    case 2:
      return (
        <div className="relative">
          <img
            src={
              "https://gcsim.app/api/assets/artifacts/" +
              artifacts[0][0] +
              "_flower.png"
            }
            alt={artifacts[0][0]}
            className="artifact-top-right-gap"
          />
          <img
            src={
              "https://gcsim.app/api/assets/artifacts/" +
              artifacts[1][0] +
              "_flower.png"
            }
            alt={artifacts[1][0]}
            className="artifact-bottom-left-gap"
          />
        </div>
      );

    default:
    case 0:
      return (
        <div>
          <img src={kuki} alt="kuki" className="relative opacity-30" />
        </div>
      );
  }
}
