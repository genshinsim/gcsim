import { PreviewCard } from "@gcsim/components";
import "@gcsim/components/src/index.css";
import { SimulationResult } from "@gcsim/types/src/generated/index.model";
import axios from "axios";
import React from "react";
import { ErrorBoundary } from "react-error-boundary";
import { Route, Switch } from "wouter";

function fallbackRender({ error }) {
  return (
    <div id="card" role="alert">
      <p className="text-white">Something went wrong:</p>
      <pre style={{ color: "red" }}>{error.message}</pre>
    </div>
  );
}

const App = ({ id, src }: { id: string; src: string }) => {
  const [ready, setReady] = React.useState<boolean>(false);
  const [error, setError] = React.useState<string>("");
  const [loaded, setLoaded] = React.useState(0);
  const [data, setData] = React.useState<SimulationResult | undefined>(
    undefined
  );
  React.useEffect(() => {
    //https://gcsim.app/api/share/db/nFLhjtD9dfFN
    axios.get("/api/share/" + src + "/" + id).then((res) => {
      console.log(res);
      if (res.data) {
        setData(res.data);
      }
    });
  }, []);
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
    return <div>no data</div>;
  }

  let status = "loading";
  if (error !== "") {
    status = "error: " + error;
  } else if (ready) {
    status = "done";
  }

  return (
    <>
      <span id="status" hidden>
        {status}
      </span>
      <ErrorBoundary fallbackRender={fallbackRender} onError={handleError}>
        <PreviewCard data={data} onImageLoaded={handleImageLoaded} />
      </ErrorBoundary>
    </>
  );
};

const Routes = () => {
  return (
    <>
      <Switch>
        <Route path="/db/:id">
          {(params) => <App id={params.id} src="db" />}
        </Route>
        <Route path="/sh/:id">
          {(params) => <App id={params.id} src="sh" />}
        </Route>
        <Route>404 Not Found</Route>
      </Switch>
    </>
  );
};

export default Routes;
