import { model } from "@gcsim/types";

//displays one database entry
export default function DBEntryView({ dbEntry }: { dbEntry: model.IDBEntry }) {
  return (
    <div className="flex flex-row bg-slate-800 max-w-fit p-4 gap-4">
      <div className="flex gap-2">
        {dbEntry.charNames &&
          dbEntry.charNames.map((charName, index) => {
            return <DBEntryCharacterPortrait charName={charName} key={index} />;
          })}
      </div>

      <DBEntryDetails
        targetCount={dbEntry.targetCount}
        meanDpsPerTarget={dbEntry.meanDpsPerTarget}
        runDate={dbEntry.runDate}
      />
      <DBEntryActions />
    </div>
  );
}

function DBEntryCharacterPortrait({ charName }: { charName: string }) {
  return (
    <div>
      {
        <img
          src={"https://gcsim.app/api/assets/avatar/" + charName + ".png"}
          alt={charName}
          className="ml-auto h-32"
        />
      }
    </div>
  );
}

function DBEntryDetails({
  targetCount,
  meanDpsPerTarget,
  runDate,
}: model.IDBEntry) {
  return (
    <div className="flex flex-col justify-center">
      {targetCount && <div>Target Count: {targetCount}</div>}
      {meanDpsPerTarget && <div>Mean DPS Per Target: {meanDpsPerTarget}</div>}
      {runDate && <div>Run Date: {JSON.stringify(runDate)}</div>}
    </div>
  );
}

function DBEntryActions() {
  const simulation_key = "test"; // TODO: get simulation key from dbEntry
  return (
    <div className="flex flex-col justify-center">
      <a
        href={`https://gcsim.app/v3/viewer/share/${simulation_key}`}
        target="_blank"
        className="text-white bg-slate-600 rounded-md p-2"
      >
        Open in Viewer
      </a>
    </div>
  );
}
