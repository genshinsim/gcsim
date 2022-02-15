import React from "react";
import { useLocation } from "wouter";
import axios from "axios";
import { Viewer } from "~src/Components/Viewer";
import {
  extractJSONStringFromBinary,
  parseAndValidate,
  Uint8ArrayFromBase64,
} from "./parse";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { viewerActions } from "./viewerSlice";
import { ResultsSummary } from "~src/types";
import { Viewport } from "~src/Components/Viewport";

axios.defaults.headers.get["Access-Control-Allow-Origin"] = "*";

type Props = {
  path: string;
  version?: string;
};

export default function Shared({ path, version = "v2" }: Props) {
  const dispatch = useAppDispatch();
  const [location, setLocation] = useLocation();
  const [msg, setMsg] = React.useState<string>("");
  const [data, setData] = React.useState<ResultsSummary | null>(null);

  React.useEffect(() => {
    //load path
    console.log("loading version: " + version);
    let url = "https://api.gcsim.app/viewer/" + path;
    if (version === "v2") {
      //do something with url
      console.log("v2: need to change url");
    }
    axios
      .get("https://api.gcsim.app/viewer/" + path)
      .then((resp) => {
        console.log(resp.data);

        let data = resp.data;

        if (data.results === undefined || data.results.length === 0) {
          setMsg("Invalid URL");
          return;
        }

        //decode base64
        const binaryStr = Uint8ArrayFromBase64(data.results[0].data);

        let jsonData = extractJSONStringFromBinary(binaryStr);

        if (jsonData.err !== "") {
          console.log(
            "error encountered extracting json string: ",
            jsonData.err
          );
          setMsg("URL does not contain valid gzipped JSON file");
          return;
        }

        //try parsing
        const parsed = parseAndValidate(jsonData.data);

        if (typeof parsed === "string") {
          setMsg(parsed);
          return;
        }

        dispatch(
          viewerActions.addViewerData({
            key: path,
            data: parsed,
          })
        );

        setData(parsed);
      })
      .catch(function (error) {
        // handle error
        setMsg("error retrieving specified url");
        console.log(error);
      });
  }, [path]);

  if (data === null && msg == "") {
    return <div>loading {path}... please wait</div>;
  }

  if (msg != "") {
    return (
      <div className="h-full p-8 flex place-content-center items-center">
        <div className="p-8 h-full w-full flex place-content-center items-center">
          <div>
            <p className="text-lg text-red-700">{msg}</p>
          </div>
        </div>
      </div>
    );
  }

  const handleClose = () => {
    setLocation("/viewer");
  };

  if (data !== null) {
    return (
      <Viewer
        data={data}
        className="h-full flex-grow"
        handleClose={handleClose}
      />
    );
  }

  return <div>Unexpected error. Please contact administrator</div>;
}
