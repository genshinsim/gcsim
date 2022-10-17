import axios from "axios";
import Pako from "pako";
import React, { useCallback } from "react";
import { RootState, useAppSelector } from "~src/store";
import { pool } from "../Sim";
import { SimResults } from "./SimResults";
import Viewer from "./Viewer";

axios.defaults.headers.get['Access-Control-Allow-Origin'] = '*';

export enum ViewTypes {
  Landing,
  Upload,
  Web,
  Local,
  Share,
  Empty, // for testing
}

type LoaderProps = {
  type: ViewTypes;
  id?: string; // only used in share
};

export const ViewerLoader = ({ type, id }: LoaderProps) => {
  switch (type) {
    case ViewTypes.Landing:
      // TODO: figure out what this should be
      return <div></div>;
    case ViewTypes.Upload:
      // TODO: show upload tsx (dropzone)
      return <div></div>;
    case ViewTypes.Web:
      return <FromState type={type} redirect="/simulator" />;
    case ViewTypes.Local:
      return <FromUrl url='http://127.0.0.1:8381/data' type={type} redirect="/viewer" />;
    case ViewTypes.Share:
      // TODO: process url function + more request props for supporting more endpoints (hastebin)
      return <FromUrl url={'/api/view/' + id} type={type} redirect="/viewer" />;
    case ViewTypes.Empty:
      return <Viewer data={{}} error={null} type={type} redirect="/viewer/empty" />;
  }
};

function Base64ToJson(base64: string) {
  const bytes = Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
  return JSON.parse(Pako.inflate(bytes, { to: 'string' }));
}

const FromUrl = ({ url, type, redirect }: { url: string, type: ViewTypes, redirect: string }) => {
  const [data, setData] = React.useState<SimResults | null>(null);
  const [error, setError] = React.useState<string | null>(null);

  const request = useCallback(() => {
    setError(null);
    axios.get(url, { timeout: 5000 }).then((resp) => {
      const out = Base64ToJson(resp.data.data);
      setData(out);
      console.log(out);
    }).catch((e) => {
      setError(e.message);
    });
  }, [url]);

  React.useEffect(() => {
    request();
  }, [request]);

  return (
    <Viewer
        data={data}
        error={error}
        type={type}
        redirect={redirect}
        retry={request} />
  );
};

// TODO: rather than using viewer state, have FromState call RunSim using sim state?
//  - determine if this is the right behavior we want. If I load the /viewer/web, should it:
//    * alert saying "no sim loaded" and confirm button redirects to /simulator (current)
//    * start running sim stored in local store, alert if not valid (proposed)
const FromState = ({ type, redirect }: { type: ViewTypes, redirect: string }) => {
  const { data, error } = useAppSelector((state: RootState) => {
    return {
      data: state.viewer_new.data,
      error: state.viewer_new.error,
    };
  });
  const cancel = useCallback(() => pool.cancel(), []);

  return (
    <Viewer
        data={data}
        error={error}
        type={type}
        redirect={redirect}
        cancel={cancel} />
  );
};