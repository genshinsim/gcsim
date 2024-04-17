import {PreviewCard} from '@gcsim/components';
import '@gcsim/components/src/index.css';
import {SimulationResult} from '@gcsim/types/src/generated/index.model';
import axios from 'axios';
import React from 'react';
import {ErrorBoundary} from 'react-error-boundary';
import {Route, Switch} from 'wouter';

function fallbackRender({error}) {
  return (
    <div id="card" role="alert">
      <input id="has-error" disabled hidden value={JSON.stringify(error)} />
      <p className="text-white">Something went wrong:</p>
      <pre style={{color: 'red'}}>{error.message}</pre>
    </div>
  );
}

const App = ({id, src}: {id: string; src: string}) => {
  const [err, setError] = React.useState<string>('');
  const [data, setData] = React.useState<SimulationResult | undefined>(
    undefined,
  );
  const [loaded, setLoaded] = React.useState(0);
  const [completed, setCompleted] = React.useState(false);
  React.useEffect(() => {
    //https://gcsim.app/api/share/db/nFLhjtD9dfFN
    const url = `/api/share/` + (src === 'db' ? 'db/' : '') + id;
    axios
      .get(url)
      .then((res) => {
        console.log(res);
        if (res.data) {
          setData(res.data);
        } else {
          setError('unexpected no data');
        }
      })
      .catch((e) => {
        setError(JSON.stringify(e));
      });
  }, []);
  React.useEffect(() => {
    if (loaded >= (data?.character_details?.length ?? 0)) {
      setCompleted(true);
    }
  }, [loaded]);

  const handleOnImageLoaded = () => {
    setLoaded(loaded + 1);
  };

  if (err !== '') {
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
      {completed ? (
        <span
          className="hidden absolute top-0 left-0"
          id="images_loaded"></span>
      ) : null}
      <PreviewCard data={data} onImageLoaded={handleOnImageLoaded} />
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
