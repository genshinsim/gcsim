import { Sample } from "@gcsim/types";
import { base64ToBytes } from "@gcsim/utils";
import axios from "axios";
import classNames from "classnames";
import Pako from "pako";
import { useCallback, useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import SamplePage from "./SamplePage";

export const LocalSample = ({}) => {
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

  return <SamplePage sample={sample} error={error} retry={request} />;
};

export const UploadSample = ({}) => {
  const [sample, setSample] = useState<Sample | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone({
    maxFiles: 1,
    noClick: sample != null || error != null
  });

  useEffect(() => {
    const file = acceptedFiles[0];
    if (file != null) {
      file.text().then((b64) => {
        try {
          setSample(JSON.parse(Pako.inflate(base64ToBytes(b64), { to: "string" })));
        } catch (e) {
          let message = 'Unknown error when parsing sample...';
          if (e instanceof Error) message = e.message;
          setError(message);
        }
      });
    }
  }, [acceptedFiles]);

  const dzClass = classNames(
      "border-dashed border-2 w-full p-8 flex place-content-center items-center cursor-pointer");

  if (sample == null && error == null) {
    return (
      <div className="p-8">
        <div {...getRootProps({ className: dzClass })}>
          <input {...getInputProps()} />
          <span className="text-lg">Drop sample file here, or click to select file</span>
        </div>
      </div>
    );
  }

  return (
    <div {...getRootProps({ className: 'dropzone' })}>
      <input {...getInputProps()} />
      <SamplePage sample={sample} error={error} />
    </div>
  );
};