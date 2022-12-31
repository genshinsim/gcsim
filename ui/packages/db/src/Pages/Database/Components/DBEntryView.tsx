import SimDurRollupCard from "@gcsim/ui/src/Pages/Viewer/Components/Overview/RollupCards/SimDurRollupCard";
import { model } from "../../../../protos_gen/protos";

//displays one database entry
export default function DBEntryView({ dbEntry }: { dbEntry: model.IDBEntry }) {
  return (
    <div className="flex flex-row bg-slate-800 max-w-fit p-4 gap-2">
      <div className="flex gap-2">
        {dbEntry.team &&
          dbEntry.team.map((character, index) => {
            return (
              <DBEntryCharacterPortrait character={character} key={index} />
            );
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

function DBEntryCharacterPortrait({
  character,
}: {
  character: model.ICharacter;
}) {
  return (
    <div>
      {character.name && <img src={character.name} alt={character.name} />}
    </div>
  );
}

function DBEntryDetails({
  targetCount,
  meanDpsPerTarget,
  runDate,
}: model.IDBEntry) {
  return (
    <div>
      {targetCount && <div>Target Count: {targetCount}</div>}
      {meanDpsPerTarget && <div>Mean DPS Per Target: {meanDpsPerTarget}</div>}
      {runDate && <div>Run Date: {JSON.stringify(runDate)}</div>}
    </div>
  );
}

function DBEntryActions() {
  const simulation_key = "test"; // TODO: get simulation key from dbEntry
  return (
    <div>
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
