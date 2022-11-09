import axios from "axios";
import { throttle } from "lodash-es";
import Pako from "pako";
import { useCallback, useEffect, useRef, useState } from "react";
import { RootState, useAppDispatch, useAppSelector } from "../../Stores/store";
import UpgradeDialog from "./UpgradeDialog";
import Viewer from "./Viewer";
import { viewerActions } from "../../Stores/viewerSlice";
import { validate as uuidValidate } from "uuid";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { SimResults } from "@gcsim/types";

// TODO: make this flush rate configurable?
export const VIEWER_THROTTLE = 100;

export enum ResultSource {
  Loaded,
  Generated,
}

export enum ViewTypes {
  Landing,
  Upload,
  Web,
  Local,
  Share,
}

type LoaderProps = {
  exec: ExecutorSupplier<Executor>;
  type: ViewTypes;
  id?: string; // only used in share
};

export const ViewerLoader = ({ exec, type, id }: LoaderProps) => {
  switch (type) {
    case ViewTypes.Landing:
      // TODO: figure out what this should be
      document.title = "gcsim - viewer";
      return <div></div>;
    case ViewTypes.Upload:
      // TODO: show upload tsx (dropzone)
      document.title = "gcsim - file upload";
      return <div></div>;
    case ViewTypes.Web:
      document.title = "gcsim - web viewer";
      return <FromState exec={exec} redirect="/simulator" />;
    case ViewTypes.Local:
      document.title = "gcsim - local";
      return <FromUrl exec={exec} url="http://127.0.0.1:8381/data" redirect="/viewer" />;
    case ViewTypes.Share:
      document.title = "gcsim - " + id;
      return <FromUrl exec={exec} url={processUrl(id)} redirect="/viewer" />;
  }
};

function processUrl(id?: string): string {
  if (id == null) {
    throw "share is missing id (should never happen)";
  }

  if (uuidValidate(id)) {
    return "/api/view/" + id;
  }
  const type = id.substring(0, id.indexOf("-"));
  id = id.substring(id.indexOf("-") + 1);
  if (type == "hb") {
    return "/hastebin/get/" + id;
  }
  return "";
}

function Base64ToJson(base64: string) {
  const bytes = Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
  return JSON.parse(Pako.inflate(bytes, { to: "string" }));
}

function useRunningState(exec: ExecutorSupplier<Executor>): boolean {
  const [isRunning, setRunning] = useState(true);

  useEffect(() => {
    const check = setInterval(() => {
      setRunning(exec().running());
    }, VIEWER_THROTTLE - 50);
    return () => clearInterval(check);
  }, [exec]);

  return isRunning;
}

const FromUrl = ({ exec, url, redirect }: {
    exec: ExecutorSupplier<Executor>, url: string, redirect: string }) => {
  const [data, setData] = useState<SimResults | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [src, setSrc] = useState<ResultSource>(ResultSource.Loaded);
  const isRunning = useRunningState(exec);

  const request = useCallback(() => {
    setError(null);
    axios
      .get(url, { timeout: 30000 })
      .then((resp) => {
        const out = Base64ToJson(resp.data.data);
        setData(out);
        console.log(out);
      })
      .catch((e) => {
        setError(e.message);
      });
  }, [url]);
  useEffect(() => request(), [request]);

  const updateResult = useRef(
    throttle(
      (res: SimResults | null) => {
        setData(res);
        setSrc(ResultSource.Generated);
      },
      VIEWER_THROTTLE,
      { leading: true, trailing: true }
    )
  );

  return (
    <>
      <Viewer
          running={isRunning}
          data={data}
          error={error}
          src={src}
          redirect={redirect}
          exec={exec}
          retry={request} />
      <UpgradeDialog
          exec={exec}
          data={data}
          redirect={redirect}
          setResult={updateResult.current}
          setError={setError} />
    </>
  );
};

// TODO: rather than using viewer state, have FromState call RunSim using sim state?
//  - determine if this is the right behavior we want. If I load the /viewer/web, should it:
//    * alert saying "no sim loaded" and confirm button redirects to /simulator (current)
//    * start running sim stored in local store, alert if not valid (proposed)
//  - This would also consolidate run logic into one place (here)
const FromState = ({ exec, redirect }: { exec: ExecutorSupplier<Executor>, redirect: string }) => {
  // TODO: conditionally create this via upgrade?
  const running = useRunningState(exec);
  const dispatch = useAppDispatch();
  const { data, error } = useAppSelector((state: RootState) => {
    return {
      data: state.viewer.data,
      error: state.viewer.error,
    };
  });

  const setResult = (result: SimResults | null) => {
    if (result == null) {
      return;
    }
    dispatch(viewerActions.setResult({ data: result }));
  };
  const updateResult = useRef(
    throttle(setResult, VIEWER_THROTTLE, { leading: true, trailing: true })
  );

  const setError = (error: string | null) => {
    if (error == null) {
      return;
    }
    dispatch(viewerActions.setError({ error: error }));
  };

  return (
    <>
      <Viewer
          running={running}
          data={data}
          src={ResultSource.Generated}
          error={error}
          redirect={redirect}
          exec={exec} />
      <UpgradeDialog
          exec={exec}
          data={data}
          redirect={redirect}
          setResult={updateResult.current}
          setError={setError} />
    </>
  );
};
