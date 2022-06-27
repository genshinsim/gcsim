import React from "react";
import axios from "axios";
import { Viewer } from "~src/Components/Viewer";
import {
  extractJSONStringFromBinary,
  parseAndValidate,
  Uint8ArrayFromBase64,
} from "./parse";
import { useAppDispatch } from "~src/store";
import { viewerActions } from "./viewerSlice";
import { ResultsSummary } from "~src/types";
import { useTranslation } from "react-i18next";

axios.defaults.headers.get["Access-Control-Allow-Origin"] = "*";

type Props = {
  path: string;
  version?: string;
  handleClose: () => void;
};

export default function Shared({ path, version = "v2", handleClose }: Props) {
  let { t } = useTranslation();

  const dispatch = useAppDispatch();
  const [msg, setMsg] = React.useState<string>("");
  const [data, setData] = React.useState<ResultsSummary | null>(null);

  React.useEffect(() => {
    //load path
    console.log("loading version: " + version);
    let url = "https://viewer.gcsim.workers.dev/" + path;
    if (path == "local") {
      url = "http://127.0.0.1:8381/data";
    }
    if (version === "v2") {
      //do something with url
      console.log("v2: need to change url");
    }
    axios
      .get(url)
      .then((resp) => {
        console.log(resp.data);

        let data = resp.data;

        // if (data.data === undefined || data.results.length === 0) {
        //   setMsg("Invalid URL");
        //   return;
        // }

        //decode base64
        const binaryStr = Uint8ArrayFromBase64(data.data);

        let jsonData = extractJSONStringFromBinary(binaryStr);

        if (jsonData.err !== "") {
          console.log(
            "error encountered extracting json string: ",
            jsonData.err
          );
          setMsg(t("viewerdashboard.url_does_not"));
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
        setMsg(t("error_retrieving_specified"));
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

  if (data !== null) {
    return (
      <div className="flex-grow">
        <Viewer
          data={data}
          className="h-full flex-grow"
          handleClose={handleClose}
        />
      </div>
    );
  }

  return <div>Unexpected error. Please contact administrator</div>;
}
