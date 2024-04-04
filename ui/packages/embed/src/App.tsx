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
      <input id="has-error" disabled hidden value={JSON.stringify(error)} />
      <p className="text-white">Something went wrong:</p>
      <pre style={{ color: "red" }}>{error.message}</pre>
    </div>
  );
}

const App = ({ id, src }: { id: string; src: string }) => {
  const [err, setError] = React.useState<string>("");
  const [data, setData] = React.useState<SimulationResult | undefined>(
    undefined
  );
  React.useEffect(() => {
    //https://gcsim.app/api/share/db/nFLhjtD9dfFN
    axios
      .get("/api/share/" + src + "/" + id)
      .then((res) => {
        console.log(res);
        if (res.data) {
          setData(res.data);
        } else {
          setError("unexpected no data");
        }
      })
      .catch((e) => {
        setError(JSON.stringify(e));
      });
  }, []);

  if (err !== "") {
    return (
      <>
        <input id="has-error" disabled hidden value={err} />
        <div>{err}</div>
      </>
    );
  }

  if (data === undefined) {
    return (
      <div id="status" className="disabled">
        no data
      </div>
    );
  }

  return (
    <ErrorBoundary fallbackRender={fallbackRender}>
      <PreviewCard data={data} />
    </ErrorBoundary>
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
