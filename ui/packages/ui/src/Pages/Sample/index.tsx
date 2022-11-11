import { Sample } from "@gcsim/types";
import axios from "axios";
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

const FromLocal = ({}) => {
  const [sample, setSample] = useState<Sample | null>(null);
  const [error, setError] = useState<string | null>(null);

  const request = useCallback(() => {
    setError(null);
    axios.get("http://127.0.0.1:8381/sample", { timeout: 30000 })
      .then((resp) => {
        setSample(resp.data);
      }).catch((e) => {
        setError(e.message);
      });
  }, []);
  useEffect(() => request(), [request]);

  return (
    <SamplePage sample={sample} error={error} retry={request} />
  );
};