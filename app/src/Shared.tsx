import React from "react";
import { Viewer } from "./Viewer/Viewer";
import { useLocation } from "wouter";
import axios from "axios";
import pako from "pako";

function Uint8ArrayFromBase64(base64: string) {
  return Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
}

axios.defaults.headers.get["Access-Control-Allow-Origin"] = "*";

export default function Shared({ path }: { path: string }) {
  const [data, setData] = React.useState<string>("");
  const [location, setLocation] = useLocation();
  const [msg, setMsg] = React.useState<string>("");

  React.useEffect(() => {
    //load path
    console.log("loading stuff");
    axios
      .get("https://api.gcsim.app/viewer/" + path)
      .then((resp) => {
        console.log(resp.data);

        let data = resp.data;
        if (data.results === undefined || data.results.length === 0) {
          setMsg("Invalid URL");
        } else {
          //decode base64
          let binaryStr = Uint8ArrayFromBase64(data.results[0].data);
          //ungzip
          try {
            const restored = pako.inflate(binaryStr, { to: "string" });
            // ... continue processing
            // console.log(restored);
            setData(restored);
            setMsg("");
          } catch (err) {
            console.log(err);
            setMsg("URL does not contain valid gzipped JSON file");
          }
        }
      })
      .catch(function (error) {
        // handle error
        setMsg("error retrieving specified url");
        console.log(error);
      });
  }, [path]);

  if (data == "" && msg == "") {
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
    setLocation("/");
  };

  return (
    <div className="p-8 h-screen flex flex-col">
      <Viewer data={data} names="grow" handleClose={handleClose} />
    </div>
  );
}
