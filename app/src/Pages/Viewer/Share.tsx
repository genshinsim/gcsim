import {
  Dialog,
  Classes,
  Button,
  Spinner,
  SpinnerSize,
} from "@blueprintjs/core";
import axios from "axios";
import pako from "pako";
import React from "react";
import { bytesToBase64 } from "./base64";
import { SimResults } from "./DataType";

export interface ShareProps {
  isOpen: boolean;
  handleClose: () => void;
  data: SimResults;
}

export default function Share(props: ShareProps) {
  const [loading, setIsLoading] = React.useState<boolean>(false);
  const [errMsg, setErrMsg] = React.useState<string>("");
  const [url, setURL] = React.useState<string>("");

  const handleUpload = () => {
    //encode data
    let compressed = pako.deflate(JSON.stringify(props.data));

    // const restored = JSON.parse(pako.inflate(compressed, { to: "string" }));

    // console.log(restored);

    let s = bytesToBase64(compressed);

    // console.log(s);
    setIsLoading(true);
    axios({
      method: "post",
      url: "https://api.gcsim.app/upload/",
      headers: { "Access-Control-Allow-Origin": "*" },
      data: {
        author: "anon",
        description: "none",
        data: s,
      },
    })
      .then((response) => {
        console.log(response);
        if (response.data.id) {
          setErrMsg("");
          setURL(response.data.id);
          setIsLoading(false);
        } else {
          setErrMsg("upload failed");
          setURL("");
          setIsLoading(false);
        }
      })
      .catch((error) => {
        console.log(error);
        setErrMsg("error encountered: " + error);
        setURL("");
        setIsLoading(false);
      });
  };

  return (
    <Dialog
      canEscapeKeyClose
      canOutsideClickClose
      autoFocus
      enforceFocus
      shouldReturnFocusOnClose
      isOpen={props.isOpen}
      onClose={() => {
        if (loading) {
          return;
        }
        setErrMsg("");
        setURL("");
        props.handleClose();
      }}
      title="Share this file"
      icon="share"
    >
      <div className="p-2">
        <div className={Classes.DIALOG_BODY}>
          {loading ? <Spinner size={SpinnerSize.LARGE} /> : null}
          {errMsg === "" ? (
            url !== "" ? (
              <div>
                Upload ok. View results at:
                <div>
                  <pre>{`https://viewer.gcsim.app/share/${url}`}</pre>
                </div>
              </div>
            ) : null
          ) : (
            <div>{errMsg}</div>
          )}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button
              intent="none"
              onClick={() => {
                navigator.clipboard
                  .writeText(`https://viewer.gcsim.app/share/${url}`)
                  .then(
                    () => {
                      alert("URL copied ok");
                    },
                    () => {
                      alert("Error copying :( Not sure what went wrong");
                    }
                  );
              }}
              disabled={url === ""}
            >
              Copy
            </Button>
            <Button intent="primary" onClick={handleUpload} disabled={loading}>
              Upload
            </Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}
