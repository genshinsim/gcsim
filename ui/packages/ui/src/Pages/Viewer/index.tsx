import axios from "axios";
import { throttle } from "lodash-es";
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

type ViewerProps = {
  exec: ExecutorSupplier<Executor>;
  id?: string; // only used in share
};

export const ShareViewer = ({ exec, id }: ViewerProps) => (
  <FromUrl exec={exec} url={processUrl(id)} redirect="/" />
);

export const LocalViewer = ({ exec }: ViewerProps) => (
  <FromUrl exec={exec} url="http://127.0.0.1:8381/data" redirect="/" />
);

export const WebViewer = ({ exec }: ViewerProps) => (
  <FromState exec={exec} redirect="/simulator" />
);

function processUrl(id?: string): string {
  if (id == null) {
    throw "share is missing id (should never happen)";
  }

  if (uuidValidate(id)) {
    return "/api/legacy-share/" + id;
  }
  return "/api/share/" + id;
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

type FromUrlProps = {
  exec: ExecutorSupplier<Executor>;
  redirect: string;
  url: string;
};

const FromUrl = ({ exec, url, redirect }: FromUrlProps) => {
  const [data, setData] = useState<SimResults | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [src, setSrc] = useState<ResultSource>(ResultSource.Loaded);
  const isRunning = useRunningState(exec);

  const request = useCallback(() => {
    setError(null);
    axios
      .get(url, { timeout: 30000 })
      .then((resp) => {
        setData(resp.data);
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
    <UpgradableViewer
        data={data}
        error={error}
        src={src}
        exec={exec}
        retry={request}
        running={isRunning}
        redirect={redirect}
        setResult={updateResult.current}
        setError={setError} />
  );
};

type FromStateProps = {
  exec: ExecutorSupplier<Executor>;
  redirect: string;
}

const FromState = ({ exec, redirect }: FromStateProps) => {
  const isRunning = useRunningState(exec);
  const { data, error } = useAppSelector((state: RootState) => {
    return {
      data: state.viewer.data,
      error: state.viewer.error,
    };
  });
  const dispatch = useAppDispatch();

  const setResult = useRef(
    throttle((result: SimResults | null) => {
      if (result == null) {
        return;
      }
      dispatch(viewerActions.setResult({ data: result }));
    }, VIEWER_THROTTLE, { leading: true, trailing: true })
  );

  const setError = (error: string | null) => {
    if (error == null) {
      return;
    }
    dispatch(viewerActions.setError({ error: error }));
  };

  return (
    <UpgradableViewer
        data={data}
        error={error}
        src={ResultSource.Generated}
        exec={exec}
        running={isRunning}
        redirect={redirect}
        setResult={setResult.current}
        setError={setError} />
  );
};

type UpgradableViewerProps = {
  data: SimResults | null;
  error: string | null;
  src: ResultSource;
  running: boolean;
  redirect: string;
  exec: ExecutorSupplier<Executor>;
  retry?: () => void;
  setResult: (r: SimResults | null) => void;
  setError: (err: string | null) => void;
}

const UpgradableViewer = (props: UpgradableViewerProps) => {
  return (
    <>
      <Viewer
          running={props.running}
          data={props.data}
          src={props.src}
          error={props.error}
          redirect={props.redirect}
          exec={props.exec}
          retry={props.retry} />
      <UpgradeDialog
          exec={props.exec}
          data={props.data}
          redirect={props.redirect}
          setResult={props.setResult}
          setError={props.setError} />
    </>
  );
};