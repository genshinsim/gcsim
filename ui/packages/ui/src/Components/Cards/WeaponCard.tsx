import { Weapon } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import placeholder from "../Images/default.png";

export function WeaponCard({ weapon, isSkeleton }: { weapon: Weapon, isSkeleton?: boolean }) {
  const { t } = useTranslation();

  const content = (
    <div className="flex flex-row">
      <div className="w-12 h-12">
        <img
          src={`/api/assets/weapons/${weapon.name}.png`}
            alt={weapon.name}
            className="object-contain w-full"
            onError={(e) => (e.target as HTMLImageElement).src = placeholder} />
      </div>
      <div className="flex-grow text-sm pl-2 flex flex-col justify-center">
        <div className="font-medium text-left">
          {t("game:weapon_names." + weapon.name).replace(
            /(.{20})..+/,
            "$1…"
          ) +
            " R" +
            weapon.refine}
        </div>
        <div className="justify-center items-center rounded-md">
          Lvl {weapon.level}/{weapon.max_level}
        </div>
      </div>
    </div>
  );

  return (
    <div className="weapon-parent ml-2 mr-2 p-2 bg-[#252A31] border-gray-600 border">
      {isSkeleton ? <div className="h-12"></div> : content}
    </div>
  );
}
