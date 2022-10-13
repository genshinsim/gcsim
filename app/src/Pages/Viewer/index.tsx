import axios from "axios";
import Pako from "pako";
import React from "react";
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
      // TODO: get from a slice state
      return <div></div>;
    case ViewTypes.Local:
      return <FromUrl url = 'http://127.0.0.1:8381/data' />;
    case ViewTypes.Share:
      // TODO: process url function + more request props for supporting more endpoints (hastebin)
      return <FromUrl url = {'/api/view/' + id} />;
    case ViewTypes.Empty:
      return <Viewer data={{}} error={null} />
  }
};

function Base64ToJson(base64: string) {
  const bytes = Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
  return JSON.parse(Pako.inflate(bytes, { to: 'string' }));
}

const FromUrl = ({ url }: { url: string }) => {
  const [data, setData] = React.useState<any | null>(null);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    axios.get(url).then((resp) => {
      const out = Base64ToJson(resp.data.data);
      setData(out);
      console.log(out);
    }).catch((e) => {
      setError(e.message);
    });
  }, [url]);

  return <Viewer data={data} error={error} />
};