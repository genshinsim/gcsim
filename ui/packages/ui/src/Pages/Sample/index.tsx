import { Sample } from "@gcsim/types";
import axios from "axios";
import Pako from "pako";
import { useCallback, useEffect, useState } from "react";
import SamplePage from "./SamplePage";

export enum SampleTypes {
  Landing,
  Upload,
  Local,
}

type LoaderProps = {
  type: SampleTypes;
}

export const SampleLoader = ({ type }: LoaderProps) => {
  switch (type) {
    case SampleTypes.Landing:
      // TODO:
      return <div></div>;
    case SampleTypes.Upload:
      // TODO:
      return <div></div>;
    case SampleTypes.Local:
      return <FromLocal />;
  }
};

// TODO: move to utils
function Base64ToJson(base64: string) {
  const bytes = Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
  return JSON.parse(Pako.inflate(bytes, { to: "string" }));
}

const FromLocal = ({}) => {
  const [sample, setSample] = useState<Sample | null>(null);
  const [error, setError] = useState<string | null>(null);

  const request = useCallback(() => {
    setError(null);
    axios.get("http://127.0.0.1:8381/sample", { timeout: 30000 })
      .then((resp) => {
        const out = Base64ToJson(resp.data);
        setSample(out);
        console.log(out);
      }).catch((e) => {
        setError(e.message);
      });
  }, []);
  useEffect(() => request(), [request]);

  return (
    <SamplePage sample={sample} error={error} retry={request} />
  );
};