import { model } from "@gcsim/types";
import React from "react";
import { ErrorBoundary } from "react-error-boundary";
import { AvatarPortrait } from "../AvatarPortait/AvatarPortrait";
import { Graphs } from "./Graphs";
import { Metadata } from "./Metadata";

function fallbackRender({ error }) {
  return (
    <div id="card" role="alert">
      <p>Something went wrong:</p>
      <pre style={{ color: "red" }}>{error.message}</pre>
    </div>
  );
}

type PreviewCardProps = {
  data?: model.ISimulationResult;
};

export const PreviewCard = ({ data }: PreviewCardProps) => {
  const [ready, setReady] = React.useState<boolean>(false);
  const [error, setError] = React.useState<string>("");
  const [loaded, setLoaded] = React.useState(0);

  const handleError = (error: Error) => {
    setError(error.message);
  };
  const handleImageLoaded = () => {
    if (data === undefined) return;

    if (loaded + 1 == data.character_details?.length) {
      console.log("all loaded");
      setReady(true);
    }
    setLoaded(loaded + 1);
  };

  // do nothing if no data...
  if (data === undefined) {
    return (
      <div>
        <span id="status">error: no data</span>
      </div>
    );
  }

  if (data.character_details === null) {
    return (
      <div>
        <span id="status">error: no character data</span>
      </div>
    );
  }

  return (
    <div className="!w-[540px] !h-[250px] bg-slate-800">
      <span id="status" hidden>
        {error !== "" ? `error: ${error}` : ready ? "ok" : "loading"}
      </span>
      <ErrorBoundary fallbackRender={fallbackRender} onError={handleError}>
        <div>
          <div className="grid grid-cols-4">
            {data.character_details?.map((c, i) => {
              return (
                <AvatarPortrait
                  key={"char-" + i}
                  char={c}
                  invalid={
                    data.incomplete_characters?.includes(c.name ?? "") ?? false
                  }
                  className="m-1"
                  onImageLoaded={handleImageLoaded}
                />
              );
            })}
          </div>
          <Metadata data={data} />
          <Graphs data={data} />
        </div>
      </ErrorBoundary>
    </div>
  );
};
