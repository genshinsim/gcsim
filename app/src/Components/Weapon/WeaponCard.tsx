import { Weapon } from "~src/types";
import { useTranslation } from "react-i18next";

export function WeaponCard({ weapon }: { weapon: Weapon }) {
  let { t } = useTranslation();

  console.log(weapon.name);
  return (
    <div className="weapon-parent ml-2 mr-2 p-2 bg-gray-800 rounded-md">
      <div className="flex flex-row">
        <div className="w-12 ">
          <img
            src={`/images/weapons/${weapon.name}.png`}
            alt={weapon.name}
            className="object-contain h-auto w-full"
          />
        </div>
        <div className="flex-grow text-sm pl-2 flex flex-col justify-center">
          <div className="font-medium text-left">
            {t("weapon_names." + weapon.name).replace(/(.{20})..+/, "$1…") +
              " R" +
              weapon.refine}
          </div>
          <div className="justify-center items-center rounded-md">
            Lvl {weapon.level}/{weapon.max_level}
          </div>
        </div>
      </div>
    </div>
  );
}
