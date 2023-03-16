import React from "react";
import { SimCard } from "Components";
import { useAppSelector, useAppDispatch } from "Store";
import { loadAllDB } from "Store/dbSlice";
import { useLocation } from "wouter";

export function AllSims() {
  const data = useAppSelector((state) => state.db.all);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    dispatch(loadAllDB());
  }, [dispatch]);

  if (data.length === 0) {
    return (
      <div className="flex flex-row place-content-center mt-2">
        Loading ....
      </div>
    );
  }

  const rows = data.map((s) => {
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
      <div className="text-white font-bold mb-2 text-xl">Showing full DB (HIDDEN DEV VIEW DON'T TELL ANYONE)</div>
      <div className="flex flex-col">{rows}</div>
    </main>
  );
}
