import { useTranslation } from "react-i18next";
import { CharacterQuickSelect } from "./CharacterQuickSelect";
import { Filter } from "./Filter";

export function ActionBar({ simCount }: { simCount: number | null }) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-row justify-between items-center w-full max-w-7xl gap-4">
      <Filter />
      <CharacterQuickSelect />

      <div className="text-base  hidden md:block md:text-2xl">
        {t<string>("db.showing_simulations", { i: simCount ?? 0 })}
      </div>
      {/* <Sorter /> */}
    </div>
  );
}
