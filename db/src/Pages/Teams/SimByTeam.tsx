import React from "react";
import { useLocation } from "wouter";
import { SimCard } from "Components";
import { useAppSelector, useAppDispatch } from "Store";
import { loadCharacter } from "Store/dbSlice";

type TeamListProps = {
  char: string;
  team: string;
};

export function SimByTeam({ char, team }: TeamListProps) {
  const charSims = useAppSelector((state) => state.db.charSims);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    if (!(char in charSims)) {
      dispatch(loadCharacter(char));
    }
  }, [charSims, dispatch]);

  if (!(char in charSims) || charSims[char].length === 0) {
    return (
      <div className="flex flex-row place-content-center mt-2">
        Sorry, <span className=" capitalize">{team.replaceAll("-", ", ")}</span>
        does not have any simulations yet
      </div>
    );
  }

  const sims = charSims[char].filter((s) => {
    const key = s.metadata.char_names
      .map((x: string) => x)
      .sort()
      .join("-");
    return team === key;
  });

  const rows = sims.map((s) => {
    const details = (
      <div className="text-sm text-white flex flex-col">
        <span>
          Total DPS: {parseInt(s.metadata.dps.mean.toFixed(0)).toLocaleString()}
        </span>
        <span>Number of Targets: {s.metadata.num_targets}</span>
        <span>
          Average DPS per Target:{" "}
          {parseInt(
            (s.metadata.dps.mean / s.metadata.num_targets).toFixed(0)
          ).toLocaleString()}
        </span>
      </div>
    );

    const action = (
      <div className="text-sm text-white">
        <a
          href={`https://gcsim.app/v3/viewer/share/${s.simulation_key}`}
          target="_blank"
        >
          Open in Viewer
        </a>
      </div>
    );

    return (
      <SimCard
        meta={s.metadata}
        key={s.simulation_key}
        summary={details}
        actions={action}
        onCharacterClick={(char) => setLocation(`/db/${char}`)}
      />
    );
  });

  return (
    <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
      <div className="text-white font-bold mb-2 text-xl">
        Showing Teams for{" "}
        <span className=" capitalize">{team.replaceAll("-", ", ")}</span>
      </div>
      <div className="flex flex-col">{rows}</div>
    </main>
  );
}
